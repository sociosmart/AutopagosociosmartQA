import { Alert, Button, Typography } from "@mui/material";
import { Box } from "@mui/system";
import { useState } from "react";
import { useDispatch } from "react-redux";
import { useParams } from "react-router-dom";
import Table from "../components/Table";
import { useGetSynchronizationDetailsQuery } from "../services/synchronization";
import { openDialog } from "../store/dialogSlice";
import { formatDate } from "../utils";

function DataCell({ value, row }) {
  const dispatch = useDispatch();

  var jsonPrettier = JSON.stringify(JSON.parse(value), null, 2);
  const onClick = () => {
    dispatch(
      openDialog({
        title: `Data recibida de la estacion ${row.external_id} en socio smart`,
        scroll: "paper",
        content: (
          <Box
            sx={{
              backgroundColor: "black",
              color: "white",
              borderRadius: 2,
              minHeight: 200,
              p: 2,
            }}
          >
            <Typography sx={{ mb: 1 }}>
              <pre>{jsonPrettier}</pre>
            </Typography>
          </Box>
        ),
      })
    );
  };
  return (
    <Box>
      {value.length > 0 && (
        <Button
          variant="contained"
          size="small"
          sx={{ ml: 1 }}
          onClick={onClick}
        >
          Ver data
        </Button>
      )}
    </Box>
  );
}

const columns = [
  {
    field: "id",
    headerName: "ID",
    flex: 1,
  },
  {
    field: "external_id",
    headerName: "ID Estacion (externo)",
    flex: 1,
  },
  {
    field: "action",
    headerName: "Accion",
    flex: 0.5,
  },
  {
    field: "data",
    headerName: "Data",
    flex: 1,
    renderCell: DataCell,
  },
  {
    field: "error_text",
    headerName: "Error",
    valueFormatter: ({ value }) => (value ? value : "Ninguno"),
    flex: 1,
  },
  {
    field: "created_at",
    headerName: "Fecha de creacion",
    valueFormatter: ({ value }) => (value ? formatDate(value) : "Invalid date"),
    flex: 1,
  },
];

export default function SynchronizationDetails() {
  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);

  const { id } = useParams();

  const {
    isLoading,
    data: { data = [], total_rows } = {},
    error,
    isError,
  } = useGetSynchronizationDetailsQuery({
    page: page + 1,
    limit: limit,
    id,
  });

  if (isError) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box>
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
