import { Alert, Badge, Box, Chip, Typography, Button } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import { useDispatch } from "react-redux";
import Table from "../components/Table";
import { useDoActionMutation, useGetPaymentsQuery } from "../services/payment";
import { useGetMeQuery } from "../services/user";
import { closeDialog, openDialog } from "../store/dialogSlice";
import {
  checkPermissionInUser,
  formatDate,
  getLastPaymentEvent,
} from "../utils";

const status = {
  pending: "Pendiente",
  failed: "Fallido",
  canceled: "Cancelado",
  paid: "Pagado",
};

const paidColor = "#0066ff";
const pendingColor = "#4e5a65";
const failedColor = "#FFB84C";
const canceledColor = "tomato";

const colorsForStatus = [
  {
    color: paidColor,
    text: status["paid"],
  },
  {
    color: pendingColor,
    text: status["pending"],
  },
  {
    color: failedColor,
    text: status["failed"],
  },
  {
    color: canceledColor,
    text: status["canceled"],
  },
];

const paymentEvents = {
  serving: "Despachando",
  serving_paused: "Pausa en el despacho",
  served: "Despachado finalizado",
  pending: "Pendiente de pago",
  failed: "Pago fallido",
  canceled: "Cancelado",
  paid: "Pagado",
  partial_refund: "Reembolso parcial (o completo)",
  funds_reserved: "Fondos reservados",
  pump_ready: "Bomba prefijada",
  internal_cancellation: "Cancelacion Interna",
  manual_action: "Accion manual",
};

export default function Payment() {
  const [limit, setLimit] = useState(10);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");

  const dispatch = useDispatch();

  const { data: user } = useGetMeQuery();

  const canDoPaymentActions = useMemo(
    () => checkPermissionInUser("can_do_payment_actions", user),
    [user],
  );

  const [doActionMutation, { error: errorAction }] = useDoActionMutation();

  useEffect(() => {
    if (!errorAction) return;
    dispatch(
      openDialog({
        content: errorAction,
      }),
    );
  }, [errorAction]);

  const {
    data: { data = [], total_rows = 0 } = {},
    isLoading,
    isError,
    error,
  } = useGetPaymentsQuery(
    { page: page + 1, limit, search },
    {
      pollingInterval: 5000,
    },
  );
  const columns = useMemo(
    () => [
      {
        field: "id",
        headerName: "ID",
        width: 300,
        hide: true,
      },
      user.is_admin && {
        field: "external_transaction_id",
        headerName: "ID Transaccion banco",
        width: 300,
        description: "ID Generado en swit o stripe",
      },
      user.is_admin && {
        field: "payment_provider",
        headerName: "Proveedor de pago",
        width: 150,
        description: "Proveedor de pago usado para esta transaccion",
        valueFormatter: ({ value }) => value.toUpperCase(),
      },
      {
        field: "customer",
        headerName: "Cliente",
        width: 300,
        hideable: false,
        valueFormatter: ({ value }) =>
          `${value.first_name} ${value.first_last_name} ${value.second_last_name}`,
      },
      {
        field: "gas_pump",
        headerName: "Estacion y # Bomba",
        width: 300,
        hideable: false,
        valueFormatter: ({ value }) =>
          `${value.gas_station?.name} - ${value.number}`,
      },
      {
        field: "fuel_type",
        headerName: "Tipo de carga",
        width: 150,
        valueFormatter: ({ value }) =>
          value.charAt(0).toUpperCase() + value.slice(1),
      },
      {
        field: "charge_type",
        headerName: "Tipo de venta",
        width: 150,
        valueFormatter: ({ value }) => {
          if (value === "by_total") return "Por total";
          return "Por litros";
        },
      },
      {
        field: "amount",
        headerName: "Total",
        width: 150,
        valueFormatter: ({ value }) => `$${value.toFixed(2)}`,
      },
      {
        field: "total_liter",
        width: 150,
        headerName: "Litros solicitados",
        description: "Litros que se solicitaron",
        valueFormatter: ({ value }) => `${value.toFixed(2)}`,
      },
      {
        field: "price",
        headerName: "Precio por litro",
        description: "Precio por litro al dia que se vendio",
        width: 150,
        valueFormatter: ({ value }) => `$${value.toFixed(2)}`,
      },
      {
        field: "discount_per_liter",
        headerName: "Descuento por litro",
        width: 200,
        valueFormatter: ({ value }) => `$${value.toFixed(2)}`,
      },
      {
        field: "refunded_amount",
        headerName: "Total regresado",
        description: "Dinero regresado si es que aplica",
        width: 150,
        valueFormatter: ({ value }) => `$${value.toFixed(2)}`,
      },
      {
        field: "real_amount_reported",
        headerName: "Carga real",
        description: "Carga real reportada por la bomba",
        width: 150,
        valueFormatter: ({ value }) => `$${value.toFixed(2)}`,
      },
      {
        field: "charge_fee",
        headerName: "Cargo por servicio",
        description: "Cargo por servicio (minimo 10 pesos)",
        width: 150,
        valueFormatter: ({ value }) => `$${value.toFixed(2)}`,
      },
      {
        field: "gm_points",
        headerName: "Puntos generados",
        description: "Puntos generados en GM",
        width: 150,
        valueFormatter: ({ value }) => `+${value.toFixed(2)}`,
      },
      {
        field: "status",
        headerName: "Estatus",
        width: 150,
        valueFormatter: ({ value }) => status[value],
      },
      {
        field: "created_at",
        headerName: "Fecha de creacion",
        width: 250,
        valueFormatter: ({ value }) => formatDate(value),
      },
      {
        field: "actions",
        type: "actions",
        width: 400,
        getActions: ({ row }) => [
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-between",
              width: 300,
              alignItems: "center",
            }}
          >
            <Badge
              badgeContent={row.events.length}
              color="secondary"
              anchorOrigin={{
                vertical: "top",
                horizontal: "left",
              }}
            >
              <Chip
                sx={{ color: "white" }}
                label="Ver eventos"
                onClick={() => {
                  dispatch(
                    openDialog({
                      title: `Eventos registrados`,
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
                          {row.events.map((data, index) => (
                            <Typography key={index} sx={{ mb: 1 }}>
                              <strong>{paymentEvents[data.type]}, </strong>
                              {formatDate(data.created_at)}
                            </Typography>
                          ))}
                        </Box>
                      ),
                    }),
                  );
                }}
              />
            </Badge>
            {(user.is_admin || canDoPaymentActions) &&
            (getLastPaymentEvent(row.events) === "paid" ||
              getLastPaymentEvent(row.events) === "funds_reserved")
              ? [
                  <Button
                    variant="text"
                    sx={{ color: "white" }}
                    onClick={() => {
                      dispatch(
                        openDialog({
                          content:
                            "¿Estas seguro que deseas intentar hacer un preset manualmente?",
                          title: "Accion",
                          customButtonText: "Si",
                          customButtonAction: () => {
                            doActionMutation({ id: row.id, action: "preset" });
                            dispatch(closeDialog());
                          },
                        }),
                      );
                    }}
                  >
                    Pre set
                  </Button>,
                  <Button
                    variant="text"
                    sx={{ color: "white" }}
                    onClick={() => {
                      dispatch(
                        openDialog({
                          content:
                            "¿Estas seguro que deseas emitir un reembolso?",
                          title: "Accion",
                          customButtonText: "Si",
                          customButtonAction: () => {
                            doActionMutation({ id: row.id, action: "refund" });
                            dispatch(closeDialog());
                          },
                        }),
                      );
                    }}
                  >
                    Reembolso
                  </Button>,
                ]
              : [<Box />, <Box />]}
          </Box>,
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
    <Box width="100%" sx={{ mb: 5 }}>
      <Box
        sx={{
          my: 3,
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-between",
        }}
      >
        {colorsForStatus.map((color) => (
          <div
            key={color.color}
            style={{
              display: "flex",
              flexDirection: "row",
              alignItems: "center",
            }}
          >
            <Box sx={{ backgroundColor: color.color, height: 15, width: 15 }} />
            <Typography sx={{ ml: 2 }} variant="subtitle2">
              {color.text}
            </Typography>
          </div>
        ))}
      </Box>
      <Table
        disableActions
        isLoading={isLoading}
        columns={columns}
        rows={data}
        limit={limit}
        setLimit={setLimit}
        setPage={setPage}
        totalRows={total_rows}
        onFilterChange={onFilterChange}
        tableProps={{
          getRowClassName: ({ row }) => `row-${row.status}`,
          sx: {
            "& .row-pending:hover": {
              backgroundColor: "#b8c0c8",
              color: "white",
            },
            "& .row-pending": {
              backgroundColor: pendingColor,
              color: "white",
            },
            "& .row-canceled": {
              backgroundColor: canceledColor,
              color: "white",
            },
            "& .row-canceled:hover": {
              backgroundColor: "#ffa494",
              color: "white",
            },
            "& .row-failed": {
              backgroundColor: failedColor,
              color: "white",
            },
            "& .row-failed:hover": {
              backgroundColor: "#FFF1DC",
              color: "white",
            },
            "& .row-paid": {
              backgroundColor: paidColor,
              color: "white",
            },
            "& .row-paid:hover": {
              backgroundColor: "#005ce6",
              color: "white",
            },
          },
        }}
      />
    </Box>
  );
}
