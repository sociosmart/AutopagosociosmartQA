import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

import { baseQueryWithReauth } from "./reAuth";

export const socioSmartGasStationApi = createApi({
  reducerPath: "socioSmartGasStationApi",
  baseQuery: fetchBaseQuery({
    baseUrl: import.meta.env.VITE_SM_HOST,
  }),
  endpoints: (builder) => ({
    getSMStations: builder.query({
      query: () => `/rest/operacion`,
    }),
  }),
});

export const gasStationApi = createApi({
  reducerPath: "gasStationApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["GasStation", "GasStationAll"],
  endpoints: (builder) => ({
    getGasStationsAll: builder.query({
      query: () => `/api/v1/gas-stations/all`,
      providesTags: ["GasStationAll"],
    }),
    getGasStations: builder.query({
      query: ({ page, limit, search }) =>
        `/api/v1/gas-stations?page=${page}&limit=${limit}&search=${search}`,
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Revisa tus parametros";
          case 401:
            return "No tienes permisos para ver las estaciones";
          case 500:
            return "Internal server error";
          default:
            return "Error al intentar conectar con el servidor";
        }
      },
      providesTags: (result) =>
        result?.data
          ? [
              ...result?.data.map(({ id }) => ({ type: "GasStation", id })),
              "GasStation",
            ]
          : ["GasStation"],
    }),
    // TODO: Add manual cache in order to invalidate whe whole cache
    createGasStation: builder.mutation({
      query: (body) => ({
        url: "/api/v1/gas-stations",
        method: "POST",
        body,
      }),
      invalidatesTags: ["GasStation", "GasStationAll"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Hay un error con sus datos, por favor modifica los parametros";
          case 401:
            return "No tienes permisos para crear una estacion";
          case 404:
            return "No se encontro el id de la estacion que quiere modificar";
          case 409:
            return "El nombre campo nombre y IP ya existen, prueba una combinacion de ip y nombre diferente";
          default:
            return "Internal server error, intenta mas tarde.";
        }
      },
    }),
    // TODO: Perform mutation in order to invalidate the correct tag for gas station
    updateGasStation: builder.mutation({
      query: ({ id, ...body }) => ({
        url: `/api/v1/gas-stations/${id}`,
        method: "PUT",
        body,
      }),
      invalidatesTags: ["GasStation", "GasStationAll"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Hay un error con sus datos, por favor modifica los parametros";
          case 404:
            return "No se encontro el id de la estacion que quiere modificar";
          case 409:
            return "El nombre campo nombre y IP ya existen, prueba una combinacion de ip y nombre diferente";
          case 401:
            return "No tienes permisos para actualizar una estacion";
          default:
            return "Internal server error, intenta mas tarde.";
        }
      },
    }),
  }),
});

export const {
  useGetGasStationsQuery,
  useUpdateGasStationMutation,
  useCreateGasStationMutation,
  useGetGasStationsAllQuery,
} = gasStationApi;

export const { useGetSMStationsQuery } = socioSmartGasStationApi;
