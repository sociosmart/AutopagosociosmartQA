import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  open: false,
  title: "",
  content: "",
  scroll: "body",
  customButtonAction: null,
  customButtonText: "",
};

export const dialogSlice = createSlice({
  name: "dialog",
  initialState,
  reducers: {
    openDialog: (state, action) => {
      const { title, content, scroll, customButtonText, customButtonAction } =
        action.payload;
      state.open = true;

      state.title = title;
      state.content = content;
      state.scroll = scroll;

      state.customButtonText = customButtonText;
      state.customButtonAction = customButtonAction;
    },
    closeDialog: (state) => {
      state.open = false;
      state.scroll = "body";

      state.customButtonText = "";
      state.customButtonAction = null;
    },
  },
});

export const { openDialog, closeDialog } = dialogSlice.actions;

export default dialogSlice.reducer;
