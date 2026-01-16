import { Alert, Box, Button, Chip, IconButton } from "@mui/material";
import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";

import EditIcon from "@mui/icons-material/Edit";
import AddIcon from "@mui/icons-material/Add";

import Table from "../components/Table";
import { useDispatch } from "react-redux";
import { openDialog } from "../store/dialogSlice";
import {
  checkIfBetweenToDates,
  checkPermissionInUser,
  formatDate,
} from "../utils";
import { useGetCampaignsQuery } from "../services/campaigns";
import { useGetMeQuery } from "../services/user";

function ShowGasStations({ values }) {
  return (
    <Box
      sx={{
        display: "flex",
        p: 1,
        maxWidth: 500,
        flexDirection: "row",
        flexWrap: "wrap",
        rowGap: 1.2,
        columnGap: 1.2,
        alignItems: "center",
      }}
    >
      {values.map((v, index) => (
        <Chip label={v.name} key={index} />
      ))}
    </Box>
  );
}

export default function CampaignList() {
  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");

  const navigate = useNavigate();

  const dispatch = useDispatch();

  const { data: user } = useGetMeQuery();

  const canAddCampaign = useMemo(
    () => checkPermissionInUser("add_campaign", user),
    [user],
  );

  const canEditCampaign = useMemo(
    () => checkPermissionInUser("edit_campaign", user),
    [user],
  );

  const {
    isLoading,
    data: { data = [], total_rows } = {},
    error,
    isError,
  } = useGetCampaignsQuery({
    page: page + 1,
    limit: limit,
    search,
  });

  const columns = useMemo(
    () => [
      {
        field: "id",
        headerName: "ID",
        width: 350,
        hide: true,
      },
      {
        field: "name",
        headerName: "Nombre de campaña",
        width: 300,
        valueGetter: ({ value, row }) =>
          `${value} ${
            checkIfBetweenToDates(row.valid_from, row.valid_to)
              ? "(Activa)"
              : ""
          }`,
      },
      {
        field: "discount",
        headerName: "Descuento por litro",
        valueFormatter: ({ value }) => `\$ ${value.toFixed(2)}`,
        width: 150,
      },
      {
        field: "valid_from",
        headerName: "Valido desde",
        valueFormatter: ({ value }) => formatDate(value),
        width: 200,
      },
      {
        field: "valid_to",
        headerName: "Valido hasta",
        valueFormatter: ({ value }) => formatDate(value),
        width: 200,
      },
      {
        field: "gas_stations",
        width: 200,
        headerName: "Aplicado en",
        description: "Que estaciones tienen aplicada esta campaña",
        renderCell: ({ value: gasStations, row }) => (
          <Button
            sx={{
              color: checkIfBetweenToDates(row.valid_from, row.valid_to)
                ? "white"
                : "black",
            }}
            onClick={() => {
              dispatch(
                openDialog({
                  title: `Estaciones asignadas campaña: ${row.name}`,
                  content: <ShowGasStations values={gasStations || []} />,
                }),
              );
            }}
          >
            Ver Estaciones
          </Button>
        ),
      },
      {
        field: "active",
        headerName: "Activo",
        type: "boolean",
        width: 100,
      },
      (user.is_admin || canEditCampaign) && {
        field: "actions",
        type: "actions",
        getActions: ({ id }) => [
          <IconButton
            label="Editar"
            onClick={() => {
              navigate(`/protected/campaigns/${id}/edit`);
            }}
          >
            <EditIcon />
          </IconButton>,
        ],
      },
    ],
    [data],
  );

  const onFilterChange = ({ quickFilterValues }) => {
    let search = quickFilterValues.join(" ");

    setSearch(search);
  };

  if (isError) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box>
      {(canAddCampaign || user.is_admin) && (
        <Box sx={{ my: 2, textAlign: "right" }}>
          <Button
            variant="contained"
            onClick={() => {
              navigate("/protected/campaigns/new");
            }}
          >
            <AddIcon sx={{ mr: 1 }} />
            Nueva Campaña
          </Button>
        </Box>
      )}
      <Table
        isLoading={isLoading}
        disableActions
        rows={data}
        columns={columns}
        onFilterChange={onFilterChange}
        totalRows={total_rows}
        limit={limit}
        setLimit={setLimit}
        page={page}
        setPage={setPage}
        tableProps={{
          getRowClassName: ({ row }) =>
            `row-${
              checkIfBetweenToDates(row.valid_from, row.valid_to)
                ? "active"
                : "not-active"
            }`,
          sx: {
            "& .row-active": {
              backgroundColor: "#1976d2 !important",
              color: "white",
            },
          },
        }}
      />
    </Box>
  );
}
