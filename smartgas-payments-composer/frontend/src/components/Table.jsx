import { useState } from "react";
import TableToolBar from "../components/TableToolBar";
import { TableNoRows } from "../components/TableNoRows";
import SaveIcon from "@mui/icons-material/Save";
import EditIcon from "@mui/icons-material/Edit";
import CancelIcon from "@mui/icons-material/Cancel";
import { DataGrid, GridRowModes } from "@mui/x-data-grid";
import {
  Alert,
  IconButton,
  LinearProgress,
  Slide,
  Snackbar,
} from "@mui/material";

export default function Table({
  rows = [],
  isLoading,
  setLimit,
  setPage,
  columns = [],
  limit = 0,
  totalRows = 0,
  addButtonText = "",
  onRowUpdate = () => {},
  addButtonAction = () => {},
  disableActions = false,
  onFilterChange = () => {},
  tableProps = {},
  disableQuickFilter = false,
  disableAddButton = false,
  disableEditButton = false,
}) {
  const [rowModesModel, setRowModesModel] = useState({});

  const [snackbar, setSnackbar] = useState(null);
  const handleCloseSnackbar = () => setSnackbar(null);

  var columnsWithActions = [...columns];
  if (!disableActions) {
    columnsWithActions.push({
      field: "actions",
      type: "actions",
      headerName: "Actions",
      width: 100,
      cellClassName: "actions",
      getActions: ({ id, row }) => {
        const isInEditMode = rowModesModel[id]?.mode === GridRowModes.Edit;

        if (isInEditMode) {
          return [
            <IconButton onClick={handleSaveClick(id)}>
              <SaveIcon />
            </IconButton>,
            <IconButton onClick={handleCancelClick(id)}>
              <CancelIcon />
            </IconButton>,
          ];
        }

        if (disableEditButton && !row.isNew) {
          return [];
        }
        return [
          <IconButton
            label="Edit"
            className="textPrimary"
            onClick={handleEditClick(id)}
            color="inherit"
          >
            <EditIcon />
          </IconButton>,
        ];
      },
    });
  }

  const handleEditClick = (id) => () => {
    setRowModesModel({ ...rowModesModel, [id]: { mode: GridRowModes.Edit } });
  };

  const handleSaveClick = (id) => () => {
    setRowModesModel({ ...rowModesModel, [id]: { mode: GridRowModes.View } });
  };

  const handleRowEditStart = (_, event) => {
    event.defaultMuiPrevented = true;
  };

  const handleRowEditStop = (_, event) => {
    event.defaultMuiPrevented = true;
  };

  const handleCancelClick = (id) => () => {
    setRowModesModel({
      ...rowModesModel,
      [id]: { mode: GridRowModes.View, ignoreModifications: true },
    });
  };

  return (
    <>
      {!!snackbar && (
        <Snackbar
          open
          anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
          onClose={handleCloseSnackbar}
          autoHideDuration={6000}
          TransitionComponent={(props) => <Slide {...props} direction="up" />}
        >
          <Alert
            {...snackbar}
            onClose={handleCloseSnackbar}
            sx={{ width: "100%", minWidth: 400 }}
          />
        </Snackbar>
      )}
      <DataGrid
        components={{
          Toolbar: TableToolBar,
          LoadingOverlay: LinearProgress,
          NoResultsOverlay: TableNoRows,
        }}
        componentsProps={{
          toolbar: {
            addText: addButtonText,
            addButtonAction: () => addButtonAction(setRowModesModel),
            disableActions,
            disableQuickFilter,
            disableAddButton,
          },
        }}
        loading={isLoading}
        columns={columnsWithActions}
        rows={rows}
        rowCount={totalRows}
        pageSize={limit}
        autoHeight
        hideFooterSelectedRowCount
        paginationMode="server"
        onFilterModelChange={onFilterChange}
        rowModesModel={rowModesModel}
        onRowModesModelChange={(newModel) => setRowModesModel(newModel)}
        filterMode="server"
        rowsPerPageOptions={[5, 10, 25, 50]}
        onPageSizeChange={(limit) => setLimit(limit)}
        editMode="row"
        experimentalFeatures={{ newEditingApi: true }}
        processRowUpdate={onRowUpdate}
        onProcessRowUpdateError={(error) => {
          setSnackbar({ children: error.message, severity: "error" });
        }}
        onRowEditStart={handleRowEditStart}
        onRowEditStop={handleRowEditStop}
        onPageChange={(page) => setPage(page)}
        {...tableProps}
      />
    </>
  );
}
