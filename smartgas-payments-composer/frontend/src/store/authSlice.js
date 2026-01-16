import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  isLoggedIn: false,
  accessToken: "",
  refreshToken: "",
};
export const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    logIn: (state, action) => {
      const { access_token, refresh_token } = action.payload;

      state.isLoggedIn = true;
      state.accessToken = access_token;
      state.refreshToken = refresh_token;

      localStorage.setItem("access_token", access_token);
      localStorage.setItem("refresh_token", refresh_token);
    },
    logOut: () => {
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");

      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");

      return initialState;
    },
  },
});

export const { logIn, logOut } = authSlice.actions;

export default authSlice.reducer;
