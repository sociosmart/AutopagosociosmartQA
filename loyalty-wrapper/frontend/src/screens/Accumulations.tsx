import { useState, useMemo, useCallback } from "react";
import { gql, useQuery, useApolloClient } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import LoadingSpinner from "../components/LoadingSpinner";
import Pagination from "../components/Pagination";
import { format } from "date-fns";
import { BaseAccumulation } from "../models/base_accumulation";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import InputSearch from "../components/InputSearch"; 


const GET_ACCUMULATIONS = gql`
  query accumulations($nextCursor: String, $prevCursor: String, $search: String) {
    accumulations(
      pagination: {
        limit: 100
        nextCursor: $nextCursor
        prevCursor: $prevCursor
      },
      search: $search
    ) {
      ... on InputValidationError {
        errors {
          field
          type
          message
        }
      }
      ... on AccumulationPagination {
        nextCursor
        prevCursor
        total
        items {
          id
          margin
          amount
          points
          marginType
          gasPrice
          generatedPoints
          usedPoints
          createdAt
          gasStation {
            name
          }
          customer {
            name
            lastName
            phoneNumber
          }
          product {
            name
          }
        }
      }
      ... on GeneralError {
        code
        message
      }
    }
  }
`;

const columnHelper = createColumnHelper<BaseAccumulation>();

export default function AcumulationsScreen() {
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

  const handleSearchChange = useCallback((newSearch: string) => {
    setSearch(newSearch);
  }, []);



  const { data, loading, error } = useQuery<{
    accumulations:
      | BasePagination<BaseAccumulation>
      | InputValidationError[]
      | GeneralError;
  }>(GET_ACCUMULATIONS, {
    variables: { ...cursor, search },
    onCompleted: () => {
      client.resetStore();
    },
  });

  const handleNextPage = () => {
    if (data?.accumulations && "__typename" in data?.accumulations) {
      const accumulationsPaginated =
        data.accumulations as BasePagination<BaseAccumulation>;
      setCursor({ nextCursor: accumulationsPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.accumulations && "__typename" in data?.accumulations) {
      const accumulationsPaginated =
        data.accumulations as BasePagination<BaseAccumulation>;
      setCursor({ prevCursor: accumulationsPaginated.prevCursor ?? null });
    }
  };

  const columns = useMemo(
    () => [
      columnHelper.accessor("customer", {
        cell: (info) => {
          let { name, lastName } = info.getValue();
          return `${name} ${lastName}`;
        },
        header: () => <span>Cliente</span>,
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
      columnHelper.accessor("gasStation.name", {
        cell: (info) => info.getValue(),
        header: () => <span>Estación</span>,
      }),
      columnHelper.accessor("marginType", {
        cell: (info) => {
          const value = info.getValue();
          return value === "by_liter" ? "Por Litro" : "Por Margen";
        },
        header: () => <span>Tipo de Margen</span>,
      }),
      columnHelper.accessor("product.name", {
        cell: (info) => info.getValue(),
        header: () => <span>Producto</span>,
      }),
      columnHelper.accessor("amount", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Monto</span>,
      }),
      columnHelper.accessor("margin", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Margen (%)</span>,
      }),
      columnHelper.accessor("points", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Puntos por peso</span>,
      }),
      columnHelper.accessor("generatedPoints", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Puntos Generados</span>,
      }),
      columnHelper.accessor("usedPoints", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Puntos Utilizados</span>,
      }),
      columnHelper.accessor("gasPrice", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Precio Gasolina</span>,
      }),
      columnHelper.accessor("createdAt", {
        cell: (info) =>
          format(new Date(info.getValue()), "dd/MM/yyyy hh:mm aaaa"),
        header: () => <span>Fecha</span>,
      }),
    ],
    [],
  );

  if (loading) return <LoadingSpinner />;

  if (error) return <p>Error al cargar los datos</p>;

  if (data?.accumulations && "__typename" in data?.accumulations) {
    if (data.accumulations.__typename === "InputValidationError") {
      return <p>Error en la validación de los campos</p>;
    } else if (data.accumulations.__typename === "GeneralError") {
      return <p>Error: {(data.accumulations as GeneralError).message}</p>;
    } else {
      const accumulations =
        (data.accumulations as BasePagination<BaseAccumulation>)?.items ?? [];
      const accumulationsPaginated =
        data.accumulations as BasePagination<BaseAccumulation>;
      const accumulationsTotal = accumulationsPaginated.total ?? 0;

      return (
        <div>

          <div className="flex flex-col sm:flex-row items-center sm:justify-between gap-2 w-full">
            <InputSearch onSearchChange={handleSearchChange} />
          </div>


          <Table
            data={accumulations}
            columns={columns}
            total={accumulationsTotal}
            paginationComponent={
              <Pagination
                onNextPage={handleNextPage}
                onPreviousPage={handlePreviousPage}
                hasNext={!!accumulationsPaginated.nextCursor}
                hasPrevious={!!accumulationsPaginated.prevCursor}
              />
            }
          />
        </div>
      );
    }
  }

  return <p>No hay registros de márgenes disponibles</p>;
}
