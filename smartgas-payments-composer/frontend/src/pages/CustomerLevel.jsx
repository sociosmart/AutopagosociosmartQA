import { Alert, Box, Tooltip, Typography } from "@mui/material";
import { useMemo, useState, useEffect } from "react";

import { LoadingButton } from "@mui/lab";
import { CloudSync } from "@mui/icons-material";
import Table from "../components/Table";
import { useDispatch } from "react-redux";
import {
  checkPermissionInUser,
  formatDate,
  generateRandomNumber,
} from "../utils";
import { useGetMeQuery } from "../services/user";
import {
  useCreateCustomerLevelMutation,
  useGetAllLevelsQuery,
  useGetCustomerLevelsQuery,
  useUpdateCustomerLevelMutation,
} from "../services/elegibility";
import { GridRowModes } from "@mui/x-data-grid";
import { showSnackbar } from "../store/snackbar";
import { useGetAllCustomersQuery } from "../services/customers";
import {
  useGetLastSyncQuery,
  useSyncNowMutation,
} from "../services/synchronization";

const range = (start, stop, step) =>
  Array.from({ length: (stop - start) / step + 1 }, (_, i) => start + i * step);

const monthToLabel = [
  "Enero",
  "Febrero",
  "Marzo",
  "Abril",
  "Mayo",
  "Junio",
  "Julio",
  "Agosto",
  "Septiembre",
  "Octubre",
  "Noviembre",
  "Diciembre",
];

const monthsToOptions = monthToLabel.map((val, i) => ({
  value: i + 1,
  label: val,
}));

export default function CustomerLevelPage() {
  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");
  const [newRow, setNewRow] = useState(null);
  const [loading, setLoading] = useState(false);
  const [pollingInterval, setPollingInterval] = useState(5 * 1000);

  const dispatch = useDispatch();

  const { data: user } = useGetMeQuery();

  const canAddCustomerLevel = useMemo(
    () => checkPermissionInUser("add_customer_level", user),
    [user],
  );

  const canEditCustomerLevel = useMemo(
    () => checkPermissionInUser("edit_customer_level", user),
    [user],
  );

  const canViewSync = useMemo(
    () => checkPermissionInUser("view_synchronizations", user),
    [user],
  );
  const canAddSync = useMemo(
    () => checkPermissionInUser("add_synchronization", user),
    [user],
  );

  const {
    data: syncData,
    isLoading: isSyncLoading,
    error: syncError,
    isError: isSyncError,
  } = useGetLastSyncQuery(
    { type: "customer_levels" },
    { pollingInterval, skip: !canViewSync },
  );

  useEffect(() => {
    if (syncData?.status === "done") {
      setPollingInterval(0);
    }
  }, [syncData]);

  const [syncNow, { isLoading: isLoadingSyncNow }] = useSyncNowMutation();

  const {
    isLoading,
    data: { data = [], total_rows } = {},
    error,
    isError,
  } = useGetCustomerLevelsQuery({
    page: page + 1,
    limit: limit,
    search,
  });

  const [createCustomerLevel] = useCreateCustomerLevelMutation();
  const [updateCustomerLevel] = useUpdateCustomerLevelMutation();

  const { data: levels = [] } = useGetAllLevelsQuery();
  const { data: customers = [] } = useGetAllCustomersQuery();

  const levelsOptions = useMemo(
    () =>
      levels.map((level) => {
        return {
          value: level["id"],
          label: `${level["name"]}`,
        };
      }),
    [levels],
  );

  const customerOptions = useMemo(
    () =>
      customers.map((customer) => {
        return {
          value: customer["id"],
          label: `${customer["first_name"]} ${customer["first_last_name"]} ${customer["second_last_name"]}`,
        };
      }),
    [customers],
  );

  const onRowUpdate = async (row) => {
    try {
      setLoading(true);
      if (row.isNew) {
        const { id } = await createCustomerLevel(row).unwrap();
        delete row["isNew"];
        setNewRow(null);
        dispatch(
          showSnackbar({
            message: "Record Agregado",
          }),
        );
        return { ...row, id };
      }
      await updateCustomerLevel(row).unwrap();
      dispatch(
        showSnackbar({
          message: "Record Actualizado",
        }),
      );
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
        field: "customer_id",
        headerName: "Cliente",
        type: "singleSelect",
        flex: 1,
        editable: true,
        valueFormatter: ({ value }) => {
          let customer = customers.find((customer) => customer.id === value);

          return customer
            ? `${customer["first_name"]} ${customer["first_last_name"]} ${customer["second_last_name"]}`
            : "No encontrado";
        },
        valueOptions: [
          { value: "", label: "Selecione opcion..." },
          ...customerOptions,
        ],
      },
      {
        field: "elegibility_level_id",
        headerName: "Nivel",
        type: "singleSelect",
        flex: 1,
        editable: true,
        valueFormatter: ({ value }) => {
          let level = levels.find((level) => level.id === value);

          return level ? level.name : "No encontrado";
        },
        valueOptions: [
          { value: "", label: "Selecione opcion..." },
          ...levelsOptions,
        ],
      },
      {
        field: "validity_year",
        headerName: "Año de validez",
        editable: true,
        flex: 1,
        type: "singleSelect",
        valueOptions: range(2015, 2030, 1),
      },
      {
        field: "validity_month",
        headerName: "Mes de validez",
        type: "singleSelect",
        flex: 1,
        editable: true,
        valueFormatter: ({ value }) => {
          return monthToLabel[value - 1] || "Mes no valido";
        },
        valueOptions: [
          { value: "", label: "Selecione opcion..." },
          ...monthsToOptions,
        ],
      },
    ],
    [data, levelsOptions, customerOptions],
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
      elegibility_level_id: "",
      customer_id: "",
      validity_month: new Date().getMonth() + 1,
      validity_year: new Date().getFullYear(),
      isNew: true,
    });

    setRowModesModel((oldModel) => ({
      ...oldModel,
      [ranId]: { mode: GridRowModes.Edit, fieldToFocus: "customer_id" },
    }));
  };

  if (isError) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box>
      {canViewSync && (
        <Box
          sx={{
            my: 1,
            display: "flex",
            flexDirection: "row-reverse",
            alignItems: "center",
          }}
        >
          {canAddSync && (
            <Tooltip title="Sincronizar Niveles de clientes">
              <LoadingButton
                loading={isLoadingSyncNow || syncData?.status === "running"}
                onClick={() => syncNow({ type: "customer_levels" })}
              >
                <CloudSync />
              </LoadingButton>
            </Tooltip>
          )}
          <Typography>
            {isSyncLoading
              ? "Cargando..."
              : isSyncError
              ? syncError
              : `Ultima actualizacion: ${formatDate(syncData.created_at)} ${
                  syncData.status === "running" || isLoadingSyncNow
                    ? "(Sincronizando...)"
                    : ""
                }`}
          </Typography>
        </Box>
      )}
      <Table
        isLoading={isLoading || loading}
        disableAddButton={!user.is_admin && !canAddCustomerLevel}
        disableEditButton={!user.is_admin && !canEditCustomerLevel}
        onRowUpdate={onRowUpdate}
        disableActions={!user.is_admin && !canEditCustomerLevel}
        addButtonAction={addButtonAction}
        rows={computedRows}
        addButtonText="Nivel a cliente"
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
