import { createApi } from "@reduxjs/toolkit/query/react";

import { baseQueryWithReauth } from "./reAuth";

export const gasPumpApi = createApi({
  reducerPath: "gasPumpApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["GasPump"],
  endpoints: (builder) => ({
    getGasPumps: builder.query({
      query: ({ page, limit, search }) =>
        `/api/v1/gas-pumps?page=${page}&limit=${limit}&search=${search}`,
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Revisa tus parametros";
          case 401:
            return "No tienes permisos para ver las bombas";
          default:
            return "Internal server error";
        }
      },
      providesTags: (result) =>
        result?.data
          ? [
              ...result?.data.map(({ id }) => ({ type: "GasPump", id })),
              "GasPump",
            ]
          : ["GasPump"],
    }),
    createGasPump: builder.mutation({
      query: (body) => ({
        url: "/api/v1/gas-pumps",
        method: "POST",
        body,
      }),
      invalidatesTags: ["GasPump"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Hay un error con sus datos, por favor modifica los parametros";
          case 401:
            return "No tienes permisos para crear una bomba";
          case 404:
            return "No se encontro el id de la estacion que quiere modificar";
          case 406:
            return "Estas mandando una estacion incorrecta, intenta con otra";
          case 409:
            return "El nombre campo numero y estacion ya existen";
          default:
            return "Internal server error, intenta mas tarde.";
        }
      },
    }),
    updateGasPump: builder.mutation({
      query: ({ id, ...body }) => ({
        url: `/api/v1/gas-pumps/${id}`,
        method: "PUT",
        body,
      }),
      //invalidatesTags: ["GasStation"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Hay un error con sus datos, por favor modifica los parametros";
          case 401:
            return "No tienes permisos para actualizar una bomba";
          case 404:
            return "No se encontro el id de la estacion que quiere modificar";
          case 406:
            return "Estas mandando una estacion incorrecta, intenta con otra";
          case 409:
            return "El nombre campo numero y estacion ya existen";
          default:
            return "Internal server error, intenta mas tarde.";
        }
      },
    }),
  }),
});

export const {
  useGetGasPumpsQuery,
  useUpdateGasPumpMutation,
  useCreateGasPumpMutation,
} = gasPumpApi;
