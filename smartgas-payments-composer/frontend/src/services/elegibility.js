import { createApi } from "@reduxjs/toolkit/query/react";
import { baseQueryWithReauth } from "./reAuth";

export const elegibilityApi = createApi({
  reducerPath: "elegibilityApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["Levels", "CustomerLevels", "AllLevels"],
  endpoints: (builder) => ({
    getAllLevels: builder.query({
      query: () => `/api/v1/elegibility/levels/all`,
      providesTags: ["AllLevels"],
    }),
    getLevels: builder.query({
      query: ({ page, limit, search }) =>
        `/api/v1/elegibility/levels?page=${page}&limit=${limit}&search=${search}`,
      providesTags: ["Levels"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 401:
            return "No tienes permisos para ver los pagos";
          default:
            return "Internal server error";
        }
      },
    }),
    getCustomerLevels: builder.query({
      query: ({ page, limit, search }) =>
        `/api/v1/elegibility/customers/levels?page=${page}&limit=${limit}&search=${search}`,
      providesTags: ["CustomerLevels"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 401:
            return "No tienes permisos para ver los pagos";
          default:
            return "Internal server error";
        }
      },
    }),
    addElegibility: builder.mutation({
      query: (body) => ({
        url: "/api/v1/elegibility/levels",
        body,
        method: "POST",
      }),
      invalidatesTags: ["Levels", "AllLevels"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 409:
            return "Este nombre de nivel ya existe, por favor intente otro";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
    }),
    createCustomerLevel: builder.mutation({
      query: (body) => ({
        url: `/api/v1/elegibility/customers/levels`,
        body,
        method: "POST",
      }),
      invalidatesTags: ["CustomerLevels"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 404:
            return "No se encontro a quien actualizar";
          case 406:
            return "No se encontro este cliente o nivel";
          case 409:
            return "Este cliente ya tiene un nivel asignado para esta fecha";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
    }),
    updateCustomerLevel: builder.mutation({
      query: ({ id, ...body }) => ({
        url: `/api/v1/elegibility/customers/levels/${id}`,
        body,
        method: "PUT",
      }),
      invalidatesTags: ["CustomerLevels"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 404:
            return "No se encontro a quien actualizar";
          case 406:
            return "No se encontro este cliente o nivel";
          case 409:
            return "Este cliente ya tiene un nivel asignado para esta fecha";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
    }),
    updateElegibilityLevel: builder.mutation({
      query: ({ id, ...body }) => ({
        url: `/api/v1/elegibility/levels/${id}`,
        body,
        method: "PUT",
      }),
      invalidatesTags: ["Levels", "AllLevels"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 404:
            return "No se encontro a quien actualizar";
          case 409:
            return "Este nombre de nivel ya existe, por favor intente otro";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
    }),
  }),
});

export const {
  useGetLevelsQuery,
  useAddElegibilityMutation,
  useUpdateElegibilityLevelMutation,
  useGetCustomerLevelsQuery,
  useGetAllLevelsQuery,
  useUpdateCustomerLevelMutation,
  useCreateCustomerLevelMutation,
} = elegibilityApi;
