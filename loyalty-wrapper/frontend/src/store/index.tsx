import { configureStore } from "@reduxjs/toolkit";
import appSlice from "./app";
import authSlice from "./authSlice";
import dialogSlice from "./dialogSlice";
import snackBarSlice from "./snackbar";
import { setupListeners } from "@reduxjs/toolkit/query/react";

export const store = configureStore({
  reducer: {
    auth: authSlice,
    app: appSlice,
    dialog: dialogSlice,
    snackbar: snackBarSlice,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }),
});

setupListeners(store.dispatch);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
