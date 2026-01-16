import { useState, useMemo, useEffect } from "react";
import { useGetGasStationsAllQuery } from "../services/gasStations";
import { Alert, Tooltip, Typography } from "@mui/material";
import Table from "../components/Table";
import { Box } from "@mui/system";

import {
  useGetGasPumpsQuery,
  useUpdateGasPumpMutation,
  useCreateGasPumpMutation,
} from "../services/gasPumps";
import { checkPermissionInUser, formatDate } from "../utils";
import {
  useGetLastSyncQuery,
  useSyncNowMutation,
} from "../services/synchronization";
import { LoadingButton } from "@mui/lab";
import { CloudSync } from "@mui/icons-material";
import { useGetMeQuery } from "../services/user";

export default function GasPump() {
  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);
  const [loading, setLoading] = useState(false);
  const [newRow, setNewRow] = useState(null);
  const [search, setSearch] = useState("");
  const [pollingInterval, setPollingInterval] = useState(5 * 1000);

  const {
    data: { data = [], total_rows = 0 } = {},
    isLoading,
    isError,
    error,
  } = useGetGasPumpsQuery({ page: page + 1, limit, search });

  const { data: user } = useGetMeQuery();

  const canViewSync = useMemo(
    () => checkPermissionInUser("view_synchronizations", user),
    [user]
  );
  const canAddSync = useMemo(
    () => checkPermissionInUser("add_synchronization", user),
    [user]
  );

  const { data: gasStations = [] } = useGetGasStationsAllQuery();

  const [updateGasPump] = useUpdateGasPumpMutation();

  const [createGasPump] = useCreateGasPumpMutation();

  const computedRows = useMemo(() => {
    if (newRow) {
      return [newRow, ...data];
    }
    return data;
  }, [newRow, data]);

  const gasStationsOptions = useMemo(
    () =>
      gasStations.map((station) => {
        return { value: station.id, label: `${station.name} - ${station.ip}` };
      }),
    [gasStations]
  );

  const {
    data: syncData,
    isLoading: isSyncLoading,
    error: syncError,
    isError: isSyncError,
  } = useGetLastSyncQuery(
    { type: "gas_pumps" },
    { pollingInterval, skip: !canViewSync }
  );

  useEffect(() => {
    if (syncData?.status === "done") {
      setPollingInterval(0);
    }
  }, [syncData]);

  const [syncNow, { isLoading: isLoadingSyncNow }] = useSyncNowMutation();

  const columns = useMemo(
    () => [
      {
        field: "id",
        headerName: "ID",
        flex: 1,
        description: "ID hasheada en UUID4",
      },
      {
        field: "external_id",
        headerName: "ID Externo",
        flex: 1,
        description: "ID que esta en socio smart",
        editable: true,
      },
      {
        field: "gas_station_id",
        headerName: "Estacion",
        type: "singleSelect",
        flex: 1,
        description: "Estacion",
        valueFormatter: ({ value, field, api }) => {
          const colDef = api.getColumn(field);

          const option = colDef.valueOptions.find(
            ({ value: optionValue }) => value === optionValue
          );

          return option?.label;
        },
        valueOptions: [
          { value: "", label: "Selecione opcion..." },
          ...gasStationsOptions,
        ],
        editable: true,
      },
      {
        field: "number",
        headerName: "Numero de bomba",
        flex: 1,
        description: "Numero de bomba en la gasolinera",
        editable: true,
      },
      {
        field: "regular_price",
        headerName: "Precio Gasolina Regular",
        flex: 1,
        type: "number",
        editable: true,
      },
      {
        field: "premium_price",
        headerName: "Precio Gasolina Premium",
        flex: 1,
        type: "number",
        editable: true,
      },
      {
        field: "diesel_price",
        headerName: "Precio Diesel",
        flex: 1,
        type: "number",
        editable: true,
      },
      {
        field: "active",
        type: "boolean",
        headerName: "Estatus",
        description: "Estatus",
        editable: true,
      },
    ],
    [gasStationsOptions]
  );

  const onRowUpdate = async (row) => {
    try {
      setLoading(true);
      if (row.isNew) {
        const { id } = await createGasPump(row).unwrap();
        delete row["isNew"];
        setNewRow(null);
        return { ...row, id };
      }
      await updateGasPump(row).unwrap();
      return row;
    } catch (error) {
      throw new Error(error);
    } finally {
      setLoading(false);
    }
  };

  const onFilterChange = ({ quickFilterValues }) => {
    let search = quickFilterValues.join(" ");

    setSearch(search);
  };

  const addButtonAction = () => {
    setNewRow({
      id: "",
      external_id: "",
      gas_station_id: "",
      regular_price: 0,
      premium_price: 0,
      diesel_price: 0,
      number: "00",
      active: true,
      isNew: true,
    });
  };

  if (isError) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box width="100%">
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
            <Tooltip title="Sincronizar Bombas">
              <LoadingButton
                loading={isLoadingSyncNow || syncData?.status === "running"}
                onClick={() => syncNow({ type: "gas_pumps" })}
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
                : `Ultima actualizacion: ${formatDate(syncData.created_at)} ${syncData.status === "running" || isLoadingSyncNow
                  ? "(Sincronizando...)"
                  : ""
                }`}
          </Typography>
        </Box>
      )}
      <Table
        //disableAddButton={!checkPermissionInUser("add_gas_pump", user)}
        //disableEditButton={!checkPermissionInUser("edit_gas_pump", user)}
        disableAddButton={!user.is_admin}
        disableEditButton={!user.is_admin}
        disableActions={!user.is_admin}
        isLoading={isLoading || loading}
        columns={columns}
        rows={computedRows}
        limit={limit}
        setLimit={setLimit}
        setPage={setPage}
        totalRows={total_rows}
        addButtonText="BOMBA"
        onRowUpdate={onRowUpdate}
        addButtonAction={addButtonAction}
        onFilterChange={onFilterChange}
      />
    </Box>
  );
}
