import { createApi } from "@reduxjs/toolkit/query/react";
import { baseQueryWithReauth } from "./reAuth";

export const settingsApi = createApi({
  reducerPath: "settingsApi",
  baseQuery: baseQueryWithReauth,
  endpoints: (builder) => ({
    getSettings: builder.query({
      query: () =>
        `/api/v1/settings`,
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 401:
            return "No tienes permisos para ver los pagos";
          default:
            return "Internal server error";
        }
      },
    }),
    setSetting: builder.mutation({
      query: (body) => ({ url: "/api/v1/settings", body, method: "POST" }),
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise, posiblemente parametro vacio.";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
    }),
  }),
});

export const { useGetSettingsQuery, useSetSettingMutation } = settingsApi;
