import { LoadingButton } from "@mui/lab";
import {
  Alert,
  Box,
  Button,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Typography,
} from "@mui/material";
import { useMemo, useState } from "react";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import Table from "../components/Table";
import { useGetSynchronizationsQuery } from "../services/synchronization";
import { openDialog } from "../store/dialogSlice";
import { formatDate } from "../utils";

const statusMap = {
  running: "Sincronizando",
  done: "Terminado",
};

const typeToHuman = {
  customer_levels: "Niveles de clientes",
  gas_pumps: "Bombas",
  gas_stations: "Estaciones",
};

function ErrorsCell({ value, row }) {
  const dispatch = useDispatch();

  const onClickErrors = () => {
    dispatch(
      openDialog({
        title: `Errores registrados - ${formatDate(row.created_at)}`,
        scroll: "paper",
        content: (
          <Box
            sx={{
              backgroundColor: "black",
              color: "white",
              borderRadius: 2,
              p: 2,
            }}
          >
            {value.map((data, index) => (
              <Typography key={index} sx={{ mb: 1 }}>
                {index + 1} - {data.text}
              </Typography>
            ))}
          </Box>
        ),
      }),
    );
  };
  return (
    <Box>
      {value.length}
      {value.length > 0 && (
        <Button
          variant="contained"
          size="small"
          sx={{ ml: 1 }}
          onClick={onClickErrors}
        >
          Detalles
        </Button>
      )}
    </Box>
  );
}

export default function SynchronizationEvents() {
  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);
  const [type, setType] = useState("");

  const navigate = useNavigate();

  const {
    isLoading,
    data: { data = [], total_rows } = {},
    error,
    isError,
  } = useGetSynchronizationsQuery({
    page: page + 1,
    limit: limit,
    type: type,
  });

  const columns = useMemo(
    () => [
      {
        field: "id",
        headerName: "ID",
        flex: 1,
        hide: true,
      },
      {
        field: "type",
        headerName: "Tipo de sincronizacion",
        flex: 1,
        valueFormatter: ({ value }) => typeToHuman[value],
      },
      {
        field: "status",
        headerName: "Estatus",
        valueFormatter: ({ value }) => statusMap[value],
        flex: 1,
      },
      {
        field: "errors",
        headerName: "Errores",
        renderCell: ErrorsCell,
        flex: 0.5,
      },
      {
        field: "created_at",
        headerName: "Fecha de creacion",
        valueFormatter: ({ value }) => formatDate(value),
        flex: 1,
      },
      {
        field: "actions",
        type: "actions",
        flex: 0.5,
        getActions: ({ row }) => [
          <LoadingButton
            loading={row.status === "running"}
            disabled={row.type === "customer_levels"}
            size="small"
            onClick={() =>
              navigate(`/protected/synchronizations/${row.id}/details`)
            }
          >
            Ver Detalles
          </LoadingButton>,
        ],
      },
    ],
    [data],
  );

  const handleChange = (event) => {
    setType(event.target.value);
  };

  if (isError) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box sx={{ width: "100%" }}>
      <Box
        sx={{
          p: 2,
          display: "flex",
          flexDirection: "row-reverse",
        }}
      >
        <FormControl sx={{ m: 1, minWidth: 120 }}>
          <InputLabel id="demo-simple-select-helper-label">Filtrar</InputLabel>
          <Select
            labelId="demo-simple-select-helper-label"
            id="demo-simple-select-helper"
            value={type}
            label="Tipo"
            onChange={handleChange}
          >
            <MenuItem value="">Todos</MenuItem>
            <MenuItem value="gas_pumps">Bombas</MenuItem>
            <MenuItem value="gas_stations">Estaciones</MenuItem>
            <MenuItem value="customer_levels">Niveles de clientes</MenuItem>
          </Select>
        </FormControl>
      </Box>
      <Table
        isLoading={isLoading}
        disableActions
        disableQuickFilter
        rows={data}
        totalRows={total_rows}
        columns={columns}
        limit={limit}
        setLimit={setLimit}
        page={page}
        setPage={setPage}
      />
    </Box>
  );
}
