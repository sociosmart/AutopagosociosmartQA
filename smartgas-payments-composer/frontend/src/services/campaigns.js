import { createApi } from "@reduxjs/toolkit/query/react";
import { baseQueryWithReauth } from "./reAuth";

export const campaignsApi = createApi({
  reducerPath: "campaignsApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["Campaigns", "Campaign"],
  endpoints: (builder) => ({
    getCampaigns: builder.query({
      query: ({ page, limit, search }) =>
        `/api/v1/campaigns?page=${page}&limit=${limit}&search=${search}`,
      providesTags: ["Campaigns"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 401:
            return "No tienes permisos para ver los pagos";
          default:
            return "Internal server error";
        }
      },
    }),
    getCampaign: builder.query({
      query: (id) => `/api/v1/campaigns/${id}`,
      providesTags: ["Campaign"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 401:
            return "No tienes permisos para ver los pagos";
          case 404:
            return "No encontrado";
          default:
            return "Internal server error";
        }
      },
    }),
    addCampaign: builder.mutation({
      query: (body) => ({ url: "/api/v1/campaigns", body, method: "POST" }),
      invalidatesTags: ["Campaigns", "Campaign"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 406:
            return "Esta intentando enviar una campa;a que no existe";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
    }),
    editCampaign: builder.mutation({
      query: ({ id, ...body }) => ({
        url: `/api/v1/campaigns/${id}`,
        body,
        method: "PUT",
      }),
      invalidatesTags: ["Campaigns", "Campaign"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 406:
            return "Esta intentando enviar una campa;a que no existe";
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
  useGetCampaignsQuery,
  useGetCampaignQuery,
  useAddCampaignMutation,
  useEditCampaignMutation,
} = campaignsApi;
