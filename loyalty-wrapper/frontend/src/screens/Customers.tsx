import { gql, useQuery, useApolloClient } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import { BaseCustomer } from "../models/base_customer";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import Table from "../components/Table";
import LoadingSpinner from "../components/LoadingSpinner";
import Pagination from "../components/Pagination";
import InputSearch from "../components/InputSearch"; 
import { useSearchParams } from "react-router-dom";
import { useState, useEffect, useCallback } from "react";

const GET_CUSTOMERS = gql`
  query Customers($nextCursor: String, $prevCursor: String, $search: String) {
    customers(
      pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor },
      search: $search
    ) {
      ... on InputValidationError {
        errors {
          field
          type
          message
        }
      }
      ... on CustomerPagination {
        nextCursor
        prevCursor
        total
        items {
          id
          externalId
          name
          lastName
          phoneNumber
          email
        }
      }
      ... on GeneralError {
        code
        message
      }
    }
  }
`;


const columnHelper = createColumnHelper<BaseCustomer>();

const columns = [
  columnHelper.accessor("externalId", {
    cell: (info) => info.getValue(),
    header: () => <span>Id Externo</span>,
  }),
  columnHelper.accessor("name", {
    cell: (info) => info.getValue(),
    header: () => <span>Nombre</span>,
  }),
  columnHelper.accessor("lastName", {
    cell: (info) => info.getValue(),
    header: () => <span>Apellidos</span>,
  }),
  columnHelper.accessor("phoneNumber", {
    cell: (info) => {
      const phoneNumber = info.getValue();
      return phoneNumber
        ? phoneNumber.replace(/^(\d{3})(\d{3})(\d{4})$/, "$1 $2 $3")
        : "";
    },
    header: () => <span>Numero Telefonico</span>,
  }),
  columnHelper.accessor("email", {
    cell: (info) => info.getValue(),
    header: () => <span>Correo</span>,
  }),
];

export default function CustomersScreen() {
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

  const handleSearchChange = useCallback((newSearch: string) => {
    setSearch(newSearch);
  }, []);


  const { data, loading, error } = useQuery<{
    customers: BasePagination<BaseCustomer> | InputValidationError[] | GeneralError;
  }>(GET_CUSTOMERS, {
    variables: {
      ...cursor,
      search
    },
    onCompleted: (data) => {
      client.resetStore();
    },
  });

  const handleNextPage = () => {
    if (data?.customers && '__typename' in data?.customers) {
      const customersPaginated = data.customers as BasePagination<BaseCustomer>;
      setCursor({ nextCursor: customersPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.customers && '__typename' in data?.customers) {
      const customersPaginated = data.customers as BasePagination<BaseCustomer>;
      setCursor({ prevCursor: customersPaginated.prevCursor ?? null });
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <p>Error al cargar los datos</p>;

  if (data?.customers && '__typename' in data?.customers) {
    if (data.customers.__typename === 'InputValidationError') {
      return <p>Error en la validación de los campos</p>;
    } else if (data.customers.__typename === 'GeneralError') {
      return <p>Error: {(data.customers as GeneralError).message}</p>;
    } else {
      const customersPaginated = data.customers as BasePagination<BaseCustomer>;
      const customers = customersPaginated?.items ?? [];
      const totalCustomers = customersPaginated.total ?? 0; //  total

      return (
        <div className="w-full">
          <InputSearch onSearchChange={handleSearchChange} />
          <Table
            data={customers}
            columns={columns}
            total={totalCustomers} 
            paginationComponent={
              <Pagination
                onNextPage={handleNextPage}
                onPreviousPage={handlePreviousPage}
                hasNext={!!customersPaginated.nextCursor}
                hasPrevious={!!customersPaginated.prevCursor}
              />
            }
          />
        </div>
      );
    }
  }

  return <p>No hay registros de clientes disponibles</p>;
}

