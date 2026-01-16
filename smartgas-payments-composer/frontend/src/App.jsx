import React, { useEffect } from "react";

import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import { RouterProvider } from "react-router-dom";
import router from "./router";
import {
  Box,
  CircularProgress,
  ThemeProvider,
  Typography,
} from "@mui/material";
import { lightTheme } from "./theme";
import { Provider, useDispatch, useSelector } from "react-redux";
import store from "./store";
import { logIn } from "./store/authSlice";
import { setLoadingApp } from "./store/app";
import Dialog from "./components/Dialog";
import Snackbar from "./components/Snackbar";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";

const baseUrl = import.meta.env.VITE_BASE_URL || "";

function App() {
  const { loading: isAppLoading } = useSelector((state) => state.app);
  const dispatch = useDispatch();

  useEffect(() => {
    const accessToken = localStorage.getItem("access_token");
    const refreshToken = localStorage.getItem("refresh_token");

    if (accessToken && refreshToken) {
      let url = window.location.pathname;

      if (!url.startsWith(`${baseUrl}/protected`)) {
        url = `${baseUrl}/protected`;
      }
      dispatch(
        logIn({
          access_token: accessToken,
          refresh_token: refreshToken,
        }),
      );
      dispatch(setLoadingApp(false));
      router.navigate(url);
    } else {
      dispatch(setLoadingApp(false));
      router.navigate(`${baseUrl}/login`);
    }
  }, []);

  if (isAppLoading) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          flexDirection: "column",
        }}
        height="100vh"
      >
        <CircularProgress />
        <Typography variant="body2" align="center" mt={3}>
          Cargando...
        </Typography>
      </Box>
    );
  }

  return (
    <ThemeProvider theme={lightTheme}>
      <Dialog />
      <Snackbar />
      <RouterProvider router={router} />
    </ThemeProvider>
  );
}

function AppWithStore() {
  return (
    <React.StrictMode>
      <LocalizationProvider dateAdapter={AdapterDateFns}>
        <Provider store={store}>
          <App />
        </Provider>
      </LocalizationProvider>
    </React.StrictMode>
  );
}

export default AppWithStore;
