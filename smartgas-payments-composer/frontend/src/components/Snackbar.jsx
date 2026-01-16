import { Alert, Snackbar } from "@mui/material";
import { useDispatch, useSelector } from "react-redux";
import { closeSnackbar } from "../store/snackbar";

export default function SnackbarComponent() {
  const dispatch = useDispatch();

  const {
    open,
    message,
    severity,
    autoHideDuration,
    position: { horizontal, vertical },
  } = useSelector((state) => state.snackbar);

  const handleClose = () => dispatch(closeSnackbar());
  return (
    <Snackbar
      open={open}
      onClose={handleClose}
      autoHideDuration={autoHideDuration}
      anchorOrigin={{ vertical, horizontal }}
    >
      <Alert severity={severity}>{message}</Alert>
    </Snackbar>
  );
}
