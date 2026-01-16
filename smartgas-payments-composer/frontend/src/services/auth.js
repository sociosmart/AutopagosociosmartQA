import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { logIn } from '../store/authSlice'

export const authApi = createApi({
  reducerPath: "authApi",
  baseQuery: fetchBaseQuery({
    baseUrl: import.meta.env.VITE_HOST
  }),
  endpoints: builder => ({
    logIn: builder.mutation(({
      query: ({ ...body }) => ({
        url: "/api/v1/auth/login",
        method: "POST",
        body
      }),
      transformErrorResponse: (response) => {
        switch (response.status) {
          case 404:
            return "Correo o Contraseña incorrecto"
          case 400:
            return "Error en los parametros enviados"

          default:
            return "Internal server error"
        }
      },
      async onQueryStarted(
        _,
        {
          dispatch,
          queryFulfilled
        }
      ) {

        try {
          let { data } = await queryFulfilled
          dispatch(logIn({ ...data }))
        } catch {
        }
      },
    })),
  })
})


export const {
  useLogInMutation,
} = authApi
