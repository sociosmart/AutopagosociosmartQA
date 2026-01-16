import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    open: false,
    message: "",
    severity: "success",
    position: {
    vertical: "bottom",
    horizontal: "right",
  },
    autoHideDuration: 5000,
};

export const snackBarSlice = createSlice({
    name: "snackbar",
    initialState,
    reducers: {
    showSnackbar: (state, action) => {
        const { message, severity, autoHideDuration, position } = action.payload;
        state.open = true;

        state.message = message;
        state.severity = severity;
        if (position) state.position = position;
        state.autoHideDuration = autoHideDuration;
    },
    closeSnackbar: (state) => {
        state.open = false;
    },
    },
});

export const { closeSnackbar, showSnackbar } = snackBarSlice.actions;

export default snackBarSlice.reducer;