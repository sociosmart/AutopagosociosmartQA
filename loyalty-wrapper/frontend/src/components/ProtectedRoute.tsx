import { RootState } from "../store"; 
import { useSelector } from "react-redux";
import { Navigate, Outlet } from "react-router-dom";

export default function ProtectedRoute() {
  const isLoggedIn = useSelector((state: RootState) => state.auth.isLoggedIn);

  if (isLoggedIn) {
    return <Outlet />;
  } else {
    return <Navigate to="/login" replace />;
  }
}
