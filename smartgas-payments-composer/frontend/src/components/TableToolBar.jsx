import { Add } from "@mui/icons-material";
import { Button } from "@mui/material";
import { Box } from "@mui/system";
import {
  GridToolbarColumnsButton,
  GridToolbarContainer,
  GridToolbarDensitySelector,
  GridToolbarExport,
  GridToolbarQuickFilter,
} from "@mui/x-data-grid";

export default function TableToolBar({
  addText,
  addButtonAction,
  disableActions,
  disableQuickFilter,
  disableAddButton,
}) {
  const disableAddButtonI = disableActions || disableAddButton;
  return (
    <GridToolbarContainer>
      <Box sx={{ flexGrow: 1 }}>
        <GridToolbarColumnsButton />
        <GridToolbarDensitySelector />
        <GridToolbarExport />
      </Box>
      <Box
        sx={{
          flexGrow: 1,
        }}
      >
        {!disableQuickFilter && <GridToolbarQuickFilter variant="standard" />}
      </Box>
      <Box sx={{ flexGrow: 1, display: "flex", justifyContent: "flex-end" }}>
        {!disableAddButtonI && (
          <Button startIcon={<Add />} onClick={addButtonAction}>
            {addText || "Agregar"}
          </Button>
        )}
      </Box>
    </GridToolbarContainer>
  );
}
