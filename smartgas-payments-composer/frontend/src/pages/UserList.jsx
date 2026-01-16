import { Alert, Box, Button, Chip, IconButton } from "@mui/material";
import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";

import EditIcon from "@mui/icons-material/Edit";
import AddIcon from "@mui/icons-material/Add";
import { useGetUsersQuery } from "../services/user";

import Table from "../components/Table";
import { useDispatch } from "react-redux";
import { openDialog } from "../store/dialogSlice";
import { formatDate } from "../utils";

function ShowPermissionsOrGroups({ values }) {
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
        <Chip label={v} key={index} />
      ))}
    </Box>
  );
}

export default function UserList() {
  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");

  const navigate = useNavigate();

  const dispatch = useDispatch();

  const {
    isLoading,
    data: { data = [], total_rows } = {},
    error,
    isError,
  } = useGetUsersQuery({
    page: page + 1,
    limit: limit,
    search,
  });

  const onFilterChange = ({ quickFilterValues }) => {
    let search = quickFilterValues.join(" ");

    setSearch(search);
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
        field: "email",
        headerName: "Email",
        flex: 1,
      },
      {
        field: "first_name",
        headerName: "Nombres",
        flex: 1,
      },
      {
        field: "last_name",
        headerName: "Apellidos",
        flex: 1,
      },
      {
        field: "groups",
        flex: 1,
        headerName: "Grupos",
        renderCell: ({ value: groups, row }) => (
          <Button
            onClick={() => {
              dispatch(
                openDialog({
                  title: `Grupos assignados al usuario ${row.first_name} ${row.last_name}`,
                  content: <ShowPermissionsOrGroups values={groups} />,
                }),
              );
            }}
          >
            Ver grupos
          </Button>
        ),
      },
      {
        field: "permissions",
        flex: 1,
        headerName: "Permisos",
        renderCell: ({ value: groups, row }) => (
          <Button
            onClick={() => {
              dispatch(
                openDialog({
                  title: `Permisos assignados al usuario ${row.first_name} ${row.last_name}`,
                  content: <ShowPermissionsOrGroups values={groups} />,
                }),
              );
            }}
          >
            Ver Permisos
          </Button>
        ),
      },
      {
        field: "gas_stations",
        flex: 1,
        headerName: "Estaciones",
        description: "Estaciones asignadas",
        renderCell: ({ value: gas_stations, row }) => (
          <Button
            onClick={() => {
              dispatch(
                openDialog({
                  title: `Estaciones assignadas al usuario ${row.first_name} ${row.last_name}`,
                  content: (
                    <ShowPermissionsOrGroups
                      values={gas_stations.map((g) => g.name)}
                    />
                  ),
                }),
              );
            }}
          >
            Ver Estaciones
          </Button>
        ),
      },
      {
        field: "created_at",
        headerName: "Miembro desde",
        valueFormatter: ({ value }) => formatDate(value),
      },
      {
        field: "active",
        headerName: "Activo",
        type: "boolean",
        flex: 1,
      },
      {
        field: "is_admin",
        headerName: "Admin",
        type: "boolean",
        flex: 1,
      },
      {
        field: "actions",
        type: "actions",
        getActions: ({ id }) => [
          <IconButton
            label="Editar"
            onClick={() => {
              navigate(`/protected/users/${id}/edit`);
            }}
          >
            <EditIcon />
          </IconButton>,
        ],
      },
    ],
    [data],
  );

  if (isError) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box>
      <Box sx={{ my: 2, textAlign: "right" }}>
        <Button
          variant="contained"
          onClick={() => {
            navigate("/protected/users/new");
          }}
        >
          <AddIcon sx={{ mr: 1 }} />
          Nuevo usuario
        </Button>
      </Box>
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
      />
    </Box>
  );
}
