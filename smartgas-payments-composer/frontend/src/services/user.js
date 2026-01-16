import { createApi } from "@reduxjs/toolkit/query/react";
import { gatherPermissions, gatherPureGroups } from "../utils";
import { baseQueryWithReauth } from "./reAuth";

export const userApi = createApi({
  reducerPath: "userApi",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["Me", "Users", "UserDetail"],
  endpoints: (builder) => ({
    getMe: builder.query({
      query: () => "/api/v1/users/me",
      providesTags: ["Me"],
      transformResponse: (response) => {
        const { permissions, groups } = response;

        response.permissions = gatherPermissions(permissions, groups);
        response.groups = gatherPureGroups(groups);
        return response;
      },
    }),
    getUserDetail: builder.query({
      query: (id) => `/api/v1/users/${id}`,
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Su parametro ID no es UUID4";
          case 404:
            return "Este usuario no existe";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
      providesTags: ["UserDetail"],
    }),
    addUser: builder.mutation({
      query: (body) => ({ url: "/api/v1/users", body, method: "POST" }),
      invalidatesTags: ["Users"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 409:
            return "Este email ya esta en uso";
          case 406:
            return "Esta intentando enviar un permiso o grupo que no existe";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
    }),
    editUser: builder.mutation({
      query: ({ id, ...body }) => ({
        url: `/api/v1/users/${id}`,
        body,
        method: "PUT",
      }),
      invalidatesTags: ["Users", "UserDetail"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 406:
            return "Esta intentando enviar un permiso o grupo que no existe";
          case 500:
            return "Internal server error";
          default:
            return "Error conectando con el servidor";
        }
      },
    }),
    getUsers: builder.query({
      query: ({ search, page, limit }) =>
        `/api/v1/users?search=${search}&page=${page}&limit=${limit}`,
      providesTags: ["Users"],
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 400:
            return "Error con los parametros enviados, por favor revise";
          case 406:
            return "Esta intentando enviar un permiso o grupo que no existe";
          case 500:
            return "Internal server error";
          case 401:
            return "No tienes permisos para ver los usuarios";
          default:
            return "Error conectando con el servidor";
        }
      },
      transformResponse: (response) => {
        let { data } = response;

        let newData = data.map((v) => {
          const { permissions, groups } = v;

          v.permissions = gatherPermissions(permissions, groups);
          v.groups = gatherPureGroups(groups);
          return v;
        });

        return { ...response, data: newData };
      },
    }),
  }),
});

export const {
  useGetMeQuery,
  useGetUsersQuery,
  useAddUserMutation,
  useGetUserDetailQuery,
  useEditUserMutation,
} = userApi;
