import { Alert, Box } from "@mui/material";
import { useMemo, useState } from "react";

import Table from "../components/Table";
import { useDispatch } from "react-redux";
import { checkPermissionInUser, generateRandomNumber } from "../utils";
import { useGetMeQuery } from "../services/user";
import {
  useAddElegibilityMutation,
  useGetLevelsQuery,
  useUpdateElegibilityLevelMutation,
} from "../services/elegibility";
import { GridRowModes } from "@mui/x-data-grid";
import { showSnackbar } from "../store/snackbar";

export default function ElegibilityLevelPage() {
  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");
  const [newRow, setNewRow] = useState(null);
  const [loading, setLoading] = useState(false);

  const dispatch = useDispatch();

  const { data: user } = useGetMeQuery();

  const canAddLevel = useMemo(
    () => checkPermissionInUser("add_elegibility_level", user),
    [user],
  );

  const canEditLevel = useMemo(
    () => checkPermissionInUser("edit_elegibility_level", user),
    [user],
  );

  const {
    isLoading,
    data: { data = [], total_rows } = {},
    error,
    isError,
  } = useGetLevelsQuery({
    page: page + 1,
    limit: limit,
    search,
  });

  const [addElegibilityLevel] = useAddElegibilityMutation();
  const [updateElegibilityLevel] = useUpdateElegibilityLevelMutation();

  const onRowUpdate = async (row) => {
    try {
      setLoading(true);
      if (row.isNew) {
        const { id } = await addElegibilityLevel(row).unwrap();
        delete row["isNew"];
        setNewRow(null);
        dispatch(
          showSnackbar({
            message: "Nivel Agregado",
          }),
        );
        return { ...row, id };
      }
      await updateElegibilityLevel(row).unwrap();
      dispatch(
        showSnackbar({
          message: "Nivel editado",
        }),
      );
      // await updateGasStation(row).unwrap();
      return row;
    } catch (error) {
      throw new Error(error);
    } finally {
      setLoading(false);
    }
  };

  const columns = useMemo(
    () => [
      {
        field: "id",
        headerName: "ID",
        flex: 1,
        hide: true,
      },
      {
        field: "name",
        headerName: "Nombre del nivel",
        flex: 1,
        editable: true,
      },
      {
        field: "discount",
        headerName: "Descuento por litro",
        valueFormatter: ({ value }) => `\$ ${value.toFixed(2)}`,
        flex: 1,
        editable: true,
        type: "number",
      },
      {
        field: "min_amount",
        headerName: "Monto Minimo",
        flex: 1,
        editable: true,
        type: "number",
        valueFormatter: ({ value }) => `\$ ${value.toFixed(2)}`,
      },
      {
        field: "min_charges",
        headerName: "Cargas Minimas",
        flex: 1,
        editable: true,
        type: "number",
      },
      {
        field: "active",
        headerName: "Activo",
        type: "boolean",
        flex: 1,
        editable: true,
      },
    ],
    [data],
  );

  const onFilterChange = ({ quickFilterValues }) => {
    let search = quickFilterValues.join(" ");

    setSearch(search);
  };

  const computedRows = useMemo(() => {
    if (newRow) {
      return [newRow, ...data];
    }
    return data;
  }, [newRow, data]);

  const addButtonAction = (setRowModesModel) => {
    let ranId = generateRandomNumber();

    setNewRow({
      id: ranId,
      name: "",
      discount: 0,
      min_amount: 0,
      min_charges: 0,
      active: true,
      isNew: true,
    });

    setRowModesModel((oldModel) => ({
      ...oldModel,
      [ranId]: { mode: GridRowModes.Edit, fieldToFocus: "name" },
    }));
  };

  if (isError) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box>
      <Table
        isLoading={isLoading || loading}
        disableAddButton={!user.is_admin && !canAddLevel}
        disableEditButton={!user.is_admin && !canEditLevel}
        onRowUpdate={onRowUpdate}
        disableActions={!user.is_admin && !canEditLevel}
        addButtonAction={addButtonAction}
        rows={computedRows}
        addButtonText="NIVEL"
        columns={columns}
        onFilterChange={onFilterChange}
        totalRows={total_rows}
        limit={limit}
        setLimit={setLimit}
        page={page}
        setPage={setPage}
      />
    </Box>
  );
}
