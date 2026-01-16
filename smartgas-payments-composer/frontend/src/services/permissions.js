import { createApi } from "@reduxjs/toolkit/query/react";
import { baseQueryWithReauth } from "./reAuth";

export const permissionApi = createApi({
  reducerPath: "permissionApi",
  baseQuery: baseQueryWithReauth,
  endpoints: (builder) => ({
    getAllPermissions: builder.query({
      query: () => `/api/v1/permissions/all`,
      transformErrorResponse: (response) => {
        switch (response.status) {
          default:
            return "Internal server error";
        }
      },
    }),
    getAllGroups: builder.query({
      query: () => `/api/v1/permissions/all-groups`,
      transformErrorResponse: (response) => {
        switch (response.status) {
          default:
            return "Internal server error";
        }
      },
    }),
  }),
});

export const { useGetAllPermissionsQuery, useGetAllGroupsQuery } =
  permissionApi;
