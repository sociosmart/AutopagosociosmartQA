import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { BrowserRouter } from "react-router-dom";
import "./App.css";
import Router from "./Router";
import { logIn, logOut } from "./store/authSlice";
import { setLoadingApp } from "./store/app";

function AppWithRouter() {
  return (
    <BrowserRouter>
      <Router />
    </BrowserRouter>
  );
}

function App() {
  const dispatch = useDispatch();
  const { loading: isAppLoading } = useSelector((state: any) => state.app);

  useEffect(() => {
    const accessToken = localStorage.getItem("access_token");
    const refreshToken = localStorage.getItem("refresh_token");

    if (accessToken && refreshToken) {
      dispatch(logIn({ access_token: accessToken, refresh_token: refreshToken }));
      dispatch(setLoadingApp(false));
    } else {
      dispatch(logOut());
      dispatch(setLoadingApp(false));
    }
  }, [dispatch]);

  if (isAppLoading) {
    return <div>Cargando...</div>;
  }

  return <AppWithRouter />;
}

export default App;
