import { createApi } from "@reduxjs/toolkit/query/react";
import { baseQueryWithReauth } from "./reAuth";

export const paymentApi = createApi({
  reducerPath: "paymentApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["Payments"],
  endpoints: (builder) => ({
    getPayments: builder.query({
      providesTags: ["Payments"],
      query: ({ page, limit, search }) =>
        `/api/v1/payments?page=${page}&limit=${limit}&search=${search}`,
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Revisa tus parametros";
          case 401:
            return "No tienes permisos para ver los pagos";
          default:
            return "Internal server error";
        }
      },
    }),
    doAction: builder.mutation({
      invalidatesTags: ["Payments"],
      query: ({ id, ...body }) => ({
        url: `/api/v1/payments/actions/${id}`,
        method: "POST",
        body,
      }),
      transformErrorResponse: (response, { dispatch }) => {
        switch (response.status) {
          case 400:
            return "Hay un error con sus datos, por favor modifica los parametros";
          case 404:
            return "No se encontro el id";
          case 412:
            return "El estatus de esta solicitud de carga es diferente de pagado";
          case 402:
            return "Esta solicitud de carga aun no ha sido pagado";
          case 503:
            return "Hubo un error al intentar prefijar la bomba remotamente"
          case 428:
            return "La bombas no estan habilitadas para ser utilizadas en AutoPago, favor de contactar a un admin."
          default:
            return "Internal server error, intenta mas tarde.";
        }
      },
    })
  }),
});

export const { useGetPaymentsQuery, useDoActionMutation } = paymentApi;
