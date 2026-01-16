import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Button,
} from "@mui/material";
import { useDispatch, useSelector } from "react-redux";
import { closeDialog } from "../store/dialogSlice";

export default function DialogComponent() {
  const { open, title, content, scroll, customButtonText, customButtonAction } =
    useSelector((state) => state.dialog);

  const dispatch = useDispatch();

  const handleClose = () => dispatch(closeDialog());
  return (
    <Dialog open={open} onClose={handleClose} scroll={scroll}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        {typeof content === "string" ? (
          <DialogContentText>{content}</DialogContentText>
        ) : (
          content
        )}
      </DialogContent>
      <DialogActions>
        {customButtonAction && (
          <Button onClick={customButtonAction}>{customButtonText}</Button>
        )}
        <Button onClick={handleClose}>Cerrar</Button>
      </DialogActions>
    </Dialog>
  );
}
