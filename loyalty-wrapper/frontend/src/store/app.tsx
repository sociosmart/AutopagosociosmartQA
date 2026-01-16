import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  loading: true
}

export const appSlice = createSlice({
  name: "app",
  initialState,
  reducers: {
    setLoadingApp: (state, action) => {
      const { loading } = action.payload

      state.loading = loading
    }
  }
})


export const { setLoadingApp } = appSlice.actions
export default appSlice.reducer