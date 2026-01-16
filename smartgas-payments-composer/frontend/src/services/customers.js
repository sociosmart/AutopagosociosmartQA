import { createApi } from "@reduxjs/toolkit/query/react";
import { baseQueryWithReauth } from "./reAuth";

export const customersApi = createApi({
  reducerPath: "customersApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["AllCustomers"],
  endpoints: (builder) => ({
    getAllCustomers: builder.query({
      query: () => `/api/v1/customers/all`,
      providesTags: ["AllCustomers"],
    }),
  }),
});

export const { useGetAllCustomersQuery } = customersApi;
