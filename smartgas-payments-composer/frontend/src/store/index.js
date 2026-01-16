import { configureStore } from "@reduxjs/toolkit";
import { setupListeners } from "@reduxjs/toolkit/dist/query";
import { authApi } from "../services/auth";
import { gasPumpApi } from "../services/gasPumps";
import {
  gasStationApi,
  socioSmartGasStationApi,
} from "../services/gasStations";
import { userApi } from "../services/user";
import { paymentApi } from "../services/payment";
import appSlice from "./app";
import authSlice from "./authSlice";
import dialogSlice from "./dialogSlice";
import { synchronizationApi } from "../services/synchronization";
import { permissionApi } from "../services/permissions";
import snackBarSlice from "./snackbar";
import { settingsApi } from "../services/setting";
import { campaignsApi } from "../services/campaigns";
import { elegibilityApi } from "../services/elegibility";
import { customersApi } from "../services/customers";

const store = configureStore({
  reducer: {
    auth: authSlice,
    app: appSlice,
    dialog: dialogSlice,
    snackbar: snackBarSlice,
    [authApi.reducerPath]: authApi.reducer,
    [userApi.reducerPath]: userApi.reducer,
    [gasStationApi.reducerPath]: gasStationApi.reducer,
    [gasPumpApi.reducerPath]: gasPumpApi.reducer,
    [socioSmartGasStationApi.reducerPath]: socioSmartGasStationApi.reducer,
    [paymentApi.reducerPath]: paymentApi.reducer,
    [synchronizationApi.reducerPath]: synchronizationApi.reducer,
    [permissionApi.reducerPath]: permissionApi.reducer,
    [settingsApi.reducerPath]: settingsApi.reducer,
    [campaignsApi.reducerPath]: campaignsApi.reducer,
    [elegibilityApi.reducerPath]: elegibilityApi.reducer,
    [customersApi.reducerPath]: customersApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }).concat(
      authApi.middleware,
      userApi.middleware,
      gasStationApi.middleware,
      gasPumpApi.middleware,
      socioSmartGasStationApi.middleware,
      paymentApi.middleware,
      synchronizationApi.middleware,
      permissionApi.middleware,
      settingsApi.middleware,
      campaignsApi.middleware,
      elegibilityApi.middleware,
      customersApi.middleware,
    ),
});
export default store;

setupListeners(store.dispatch);
