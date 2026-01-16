import { fetchBaseQuery } from "@reduxjs/toolkit/query";
//import { tokenReceived, loggedOut } from './authSlice'
import { Mutex } from "async-mutex";
import { logIn, logOut } from "../store/authSlice";
import { setAuthorizationHeaders } from "../utils";

// create a new mutex
const mutex = new Mutex();
const baseQuery = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_HOST,
  prepareHeaders: setAuthorizationHeaders,
});

export const baseQueryWithReauth = async (args, api, extraOptions) => {
  // wait until the mutex is available without locking it
  await mutex.waitForUnlock();
  let result = await baseQuery(args, api, extraOptions);
  if (result.error && result.error.status === 401) {
    if (!mutex.isLocked()) {
      const release = await mutex.acquire();
      try {
        const refreshToken = api.getState().auth.refreshToken;
        const refreshResult = await fetchBaseQuery({
          baseUrl: import.meta.env.VITE_HOST + "/api/v1/auth/refresh-token",
          body: {
            refresh_token: refreshToken,
          },
          method: "POST",
        })("", api, extraOptions);
        if (refreshResult.data) {
          api.dispatch(logIn(refreshResult.data));
          // retry the initial query
          result = await baseQuery(args, api, extraOptions);
        } else {
          api.dispatch(logOut());
        }
      } finally {
        // release must be called once the mutex should be released again.
        release();
      }
    } else {
      // wait until the mutex is available without locking it
      await mutex.waitForUnlock();
      result = await baseQuery(args, api, extraOptions);
    }
  }
  return result;
};
