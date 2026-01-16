import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  isLoggedIn: !!localStorage.getItem("access_token"), 
  accessToken: localStorage.getItem("access_token") || "",
  refreshToken: localStorage.getItem("refresh_token") || "",
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
      return { isLoggedIn: false, accessToken: "", refreshToken: "" };
    },
  },
});

export const { logIn, logOut } = authSlice.actions;
export default authSlice.reducer;
