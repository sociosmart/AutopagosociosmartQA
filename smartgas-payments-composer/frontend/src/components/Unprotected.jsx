import { useSelector } from "react-redux";
import { Navigate } from "react-router-dom";

export default function ({ children }) {

  const { isLoggedIn } = useSelector(state => state.auth)

  if (isLoggedIn) {
    return <Navigate to="/protected" replace />
  }

  return children

}
