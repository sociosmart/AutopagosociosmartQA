import { gql, useQuery } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import LoadingSpinner from "../components/LoadingSpinner";
import { useState } from "react";
import Pagination from "../components/Pagination";
import { useApolloClient } from "@apollo/client";
import { BaseReport } from "../models/base_report";

const GET_REPORT = gql`
    query AccumulationsReport($nextCursor: String, $prevCursor: String) {
        accumulationsReport(pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor }) {
        ... on InputValidationError {
        errors {
            field
            type
            message
        }
        }
        ... on GeneralError {
        code
        message
        }
        ... on AccumulationReportPagination {
        nextCursor
        prevCursor
        total
        items {
            id
            totalTransactions
            avgAmount
            totalAmount
            totalGeneratedPoints
            totalUsedPoints
            totalPoints
            customer {
            name
            lastName
            phoneNumber
            }
        }
    }
}
}
`;






const columnHelper = createColumnHelper<BaseReport>();

const columns = [
  columnHelper.accessor("customer.name", {
    cell: (info) => info.getValue(),
    header: () => <span>Nombres</span>,
  }),
  columnHelper.accessor("customer.lastName", {
    cell: (info) => info.getValue(),
    header: () => <span>Apellidos</span>,
  }),
  columnHelper.accessor("customer.phoneNumber", {
    cell: (info) => {
      const phoneNumber = info.getValue();
      return phoneNumber
        ? phoneNumber.replace(/^(\d{3})(\d{3})(\d{4})$/, "$1 $2 $3")
        : "";
    },
    header: () => <span>Numero Telefonico</span>,
  }),
  columnHelper.accessor("totalTransactions", {
    cell: (info) => info.getValue(),
    header: () => <span>Transacciones totales</span>,
  }),
  columnHelper.accessor("avgAmount", {
    cell: (info) => info.getValue().toFixed(2),
    header: () => <span>Ticket promedio</span>,
  }),
    columnHelper.accessor("totalAmount", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Monto Acumulado</span>,
    }),
    columnHelper.accessor("totalGeneratedPoints", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Puntos generados</span>,
    }),
    columnHelper.accessor("totalUsedPoints", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Puntos usados</span>,
    }),
    columnHelper.accessor("totalPoints", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Puntos totales</span>,
    }),

];

export default function ReportScreen() {
    const [cursor, setCursor] = useState<object>({});
    const client = useApolloClient();
  
    const { data, loading, error } = useQuery<{
        accumulationsReport: BasePagination<BaseReport> | InputValidationError[] | GeneralError;
    }>(GET_REPORT, {
      variables: {
        ...cursor,
      },
      onCompleted: (data) => {
        client.resetStore();
      },
    });


  const handleNextPage = () => {
    if (data?.accumulationsReport && '__typename' in data?.accumulationsReport) {
      const customersPaginated = data.accumulationsReport as BasePagination<BaseReport>;
      setCursor({ nextCursor: customersPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.accumulationsReport && '__typename' in data?.accumulationsReport) {
      const customersPaginated = data.accumulationsReport as BasePagination<BaseReport>;
      setCursor({ prevCursor: customersPaginated.prevCursor ?? null });
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <p>Error al cargar los datos</p>;


  if (data?.accumulationsReport && '__typename' in data?.accumulationsReport) {
    if (data.accumulationsReport.__typename === 'InputValidationError') {
      return <p>Error en la validación de los campos</p>;
    } else if (data.accumulationsReport.__typename === 'GeneralError') {
      return <p>Error: {(data.accumulationsReport as GeneralError).message}</p>;
    } else {
      const users = (data.accumulationsReport as BasePagination<BaseReport>)?.items ?? [];
      const usersPaginated = data.accumulationsReport as BasePagination<BaseReport>;
      const totalReport = usersPaginated.total ?? 0;


      return (
        <div>
          <Table
            data={users}
            columns={columns}
            total={totalReport}
            paginationComponent={<Pagination
              onNextPage={handleNextPage}
              onPreviousPage={handlePreviousPage}
              hasNext={!!usersPaginated.nextCursor}
              hasPrevious={!!usersPaginated.prevCursor}
              />}
          />
        </div>
      );
    }
  }

  return <p>No hay registros de clientes disponibles</p>;
}




