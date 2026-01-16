import { useState, useMemo, useEffect } from "react";
import {
  useCreateGasStationMutation,
  useGetGasStationsQuery,
  useGetSMStationsQuery,
  useUpdateGasStationMutation,
} from "../services/gasStations";
import { Alert, Tooltip, Typography } from "@mui/material";
import Table from "../components/Table";
import { Box } from "@mui/system";
import {
  useGetLastSyncQuery,
  useSyncNowMutation,
} from "../services/synchronization";
import { checkPermissionInUser, formatDate } from "../utils";
import { LoadingButton } from "@mui/lab";
import { CloudSync } from "@mui/icons-material";
import { useGetMeQuery } from "../services/user";

export default function GasStation() {
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
  } = useGetGasStationsQuery({ page: page + 1, limit, search });

  const { data: smStations = [] } = useGetSMStationsQuery();

  const { data: user } = useGetMeQuery();

  const canViewSync = useMemo(
    () => checkPermissionInUser("view_synchronizations", user),
    [user]
  );
  const canAddSync = useMemo(
    () => checkPermissionInUser("add_synchronization", user),
    [user]
  );

  const {
    data: syncData,
    isLoading: isSyncLoading,
    error: syncError,
    isError: isSyncError,
  } = useGetLastSyncQuery(
    { type: "gas_stations" },
    { pollingInterval, skip: !canViewSync }
  );

  useEffect(() => {
    if (syncData?.status === "done") {
      setPollingInterval(0);
    }
  }, [syncData]);

  const smStationsOptions = useMemo(
    () =>
      smStations.map((station) => {
        return {
          value: station["Cve_PuntoDeVenta"],
          label: `${station["NombreComercial"]} - ${station["Num_PermisoCRE"]}`,
        };
      }),
    [smStations]
  );

  const [syncNow, { isLoading: isLoadingSyncNow }] = useSyncNowMutation();

  const columns = [
    {
      field: "id",
      headerName: "ID",
      flex: 1,
      description: "ID hasheada en UUID4",
    },
    {
      field: "external_id",
      headerName: "ID Externo",
      type: "singleSelect",
      flex: 1,
      description:
        "Este es el ID que se encuentra en el sistema de gasomarshal",
      editable: true,
      valueOptions: [
        { value: "", label: "Selecione opcion..." },
        ...smStationsOptions,
      ],
    },
    {
      field: "cre_permission",
      headerName: "CRE",
      flex: 1,
      description: "CRE",
    },
    {
      field: "name",
      headerName: "Estacion",
      flex: 1,
      description: "Nombre de la estacion",
    },
    {
      field: "ip",
      headerName: "IP",
      flex: 0.5,
      description: "IP en la vpn",
      editable: true,
    },
    {
      field: "active",
      type: "boolean",
      headerName: "Estatus",
      description: "Estatus",
      editable: true,
    },
  ];

  const [updateGasStation] = useUpdateGasStationMutation();

  const [createGasStation] = useCreateGasStationMutation();

  const computedRows = useMemo(() => {
    if (newRow) {
      return [newRow, ...data];
    }
    return data;
  }, [newRow, data]);

  const onRowUpdate = async (row) => {
    const smStation = smStations.find(
      (station) => station["Cve_PuntoDeVenta"] === row.external_id
    );

    if (!smStation) {
      throw new Error(
        "Por favor seleccione la estacion a la cual quiere hacer referencia"
      );
    }

    row = {
      ...row,
      name: smStation["NombreComercial"],
      cre_permission: smStation["Num_PermisoCRE"],
    };

    try {
      setLoading(true);
      if (row.isNew) {
        const { id } = await createGasStation(row).unwrap();
        delete row["isNew"];
        setNewRow(null);
        return { ...row, id };
      }
      await updateGasStation(row).unwrap();
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
      name: "",
      external_id: "",
      ip: "0.0.0.0",
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
            <Tooltip title="Sincronizar Estaciones">
              <LoadingButton
                loading={isLoadingSyncNow || syncData?.status === "running"}
                onClick={() => syncNow({ type: "gas_stations" })}
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
        //disableAddButton={!checkPermissionInUser("add_gas_station", user)}
        //disableEditButton={!checkPermissionInUser("edit_gas_station", user)}
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
        addButtonText="ESTACION"
        onRowUpdate={onRowUpdate}
        addButtonAction={addButtonAction}
        onFilterChange={onFilterChange}
      />
    </Box>
  );
}
