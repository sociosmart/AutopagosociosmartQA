import { gql, useQuery } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import { BaseCustomerLevel } from "../models/base_customer_level";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import LoadingSpinner from "../components/LoadingSpinner";
import { useState } from "react";
import Pagination from "../components/Pagination";
import { useApolloClient } from "@apollo/client";
import { format } from "date-fns";

const GET_CUSTOMERS_LEVEL = gql`
  query CustomerLevels($nextCursor: String, $prevCursor: String) {
    customerLevels(pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor }) {
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
    ... on CustomerLevelPagination {
      nextCursor
      prevCursor
      total
      items {
        id
        customer {
          name
          lastName
          phoneNumber
        }
        level {
          name
        }
        startDate
        endDate
      }
    }
    }
  }
`;

const columnHelper = createColumnHelper<BaseCustomerLevel>();

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
  columnHelper.accessor("level.name", {
    cell: (info) => info.getValue(),
    header: () => <span>Nivel</span>,
  }),
  columnHelper.accessor("startDate", {
    cell: (info) =>
        format(new Date(info.getValue()), "dd/MM/yyyy hh:mm aaaa"),
    header: () => <span>Fecha de Inicio</span>,
  }),
  columnHelper.accessor("endDate", {
    cell: (info) =>
        format(new Date(info.getValue()), "dd/MM/yyyy hh:mm aaaa"),
    header: () => <span>Fecha Final</span>,
  }),
];

export default function CustomerLevelsScreen() {
    const [cursor, setCursor] = useState<object>({});
    const client = useApolloClient();
  
    const { data, loading, error } = useQuery<{
        customerLevels: BasePagination<BaseCustomerLevel> | InputValidationError[] | GeneralError;
    }>(GET_CUSTOMERS_LEVEL, {
      variables: {
        ...cursor,
      },
      onCompleted: (data) => {
        client.resetStore();
      },
    });

  const handleNextPage = () => {
    if (data?.customerLevels && '__typename' in data?.customerLevels) {
      const customersPaginated = data.customerLevels as BasePagination<BaseCustomerLevel>;
      setCursor({ nextCursor: customersPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.customerLevels && '__typename' in data?.customerLevels) {
      const customersPaginated = data.customerLevels as BasePagination<BaseCustomerLevel>;
      setCursor({ prevCursor: customersPaginated.prevCursor ?? null });
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <p>Error al cargar los datos</p>;


  if (data?.customerLevels && '__typename' in data?.customerLevels) {
    if (data.customerLevels.__typename === 'InputValidationError') {
      return <p>Error en la validación de los campos</p>;
    } else if (data.customerLevels.__typename === 'GeneralError') {
      return <p>Error: {(data.customerLevels as GeneralError).message}</p>;
    } else {
      const users = (data.customerLevels as BasePagination<BaseCustomerLevel>)?.items ?? [];
      const usersPaginated = data.customerLevels as BasePagination<BaseCustomerLevel>;
      const usersTotal = usersPaginated.total ?? 0;

      return (
        <div>
          <Table
            data={users}
            columns={columns}
            total={usersTotal}
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

