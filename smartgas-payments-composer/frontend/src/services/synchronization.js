import { createApi } from "@reduxjs/toolkit/query/react";
import { baseQueryWithReauth } from "./reAuth";
import { openDialog } from "../store/dialogSlice";
import { gasPumpApi } from "./gasPumps";
import { gasStationApi } from "./gasStations";
import { elegibilityApi } from "./elegibility";

export const synchronizationApi = createApi({
  reducerPath: "synchronizationApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["Synchronization", "SynchronizationDetails"],
  endpoints: (builder) => ({
    getSynchronizations: builder.query({
      query: ({ limit, page, type }) =>
        `/api/v1/synchronizations?type=${type}&page=${page}&limit=${limit}`,
      providesTags: ["Synchronization"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Revisa tus parametros";
          case 401:
            return "No tienes permisos para ver las sincronizaciones";
          default:
            return "Internal server error";
        }
      },
    }),
    getSynchronizationDetails: builder.query({
      query: ({ id, limit, page }) =>
        `/api/v1/synchronizations/${id}/details?limit=${limit}&page=${page}`,
      providesTags: ["SynchronizationDetails"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Revisa tus parametros";
          default:
            return "Internal server error";
        }
      },
    }),
    getLastSync: builder.query({
      query: ({ type }) => `/api/v1/synchronizations/last?type=${type}`,
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Revisa tus parametros";
          case 404:
            return "No se ha hecho ninguna sincronizacion aun";
          default:
            return "Internal server error";
        }
      },
      providesTags: ["Synchronization"],
    }),
    syncNow: builder.mutation({
      query: (body) => ({
        url: `/api/v1/synchronizations/now`,
        method: "POST",
        body,
      }),
      invalidatesTags: ["Synchronization", "SynchronizationDetails"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Revisa tus parametros";
          case 423:
            return "Ya hay una sincronizacion que esta corriendo actualmente";
          default:
            return "Internal server error";
        }
      },
      async onQueryStarted({ type }, { dispatch, queryFulfilled }) {
        try {
          await queryFulfilled;

          //Invalidate tags
          if (type === "gas_pumps") {
            dispatch(gasPumpApi.util.invalidateTags(["GasPump"]));
          } else if (type === "gas_stations") {
            dispatch(gasStationApi.util.invalidateTags(["GasStation"]));
          } else if (type === "customer_levels") {
            dispatch(elegibilityApi.util.invalidateTags(["CustomerLevels"]));
          }
        } catch (error) {
          dispatch(
            openDialog({
              title: "Error",
              content: error.error,
            }),
          );
        }
      },
    }),
  }),
});

export const {
  useGetLastSyncQuery,
  useSyncNowMutation,
  useGetSynchronizationsQuery,
  useGetSynchronizationDetailsQuery,
} = synchronizationApi;
