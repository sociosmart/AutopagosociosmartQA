import { useState, useMemo, useCallback } from "react";
import { gql, useQuery, useApolloClient } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import LoadingSpinner from "../components/LoadingSpinner";
import Pagination from "../components/Pagination";
import { useNavigate } from "react-router-dom";
import { BaseMargin } from "../models/base_margin";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import { FaPencilAlt } from "react-icons/fa";
import InputSearch from "../components/InputSearch"; 


const GET_MARGINS = gql`
  query GasStationMargins($nextCursor: String, $prevCursor: String, $search: String) {
    gasStationsMargin(
      pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor },
      search: $search) {
      ...on InputValidationError {
        errors {
          field,
          type,
          message
        }
      }
      ... on GasStationMarginPagination {
        nextCursor
        prevCursor
        total
        items {
          id
          marginType
          margin
          points
          product {
            name  
          }
          gasStation {
            name  
          }
        }
      }
      ...on GeneralError {
        code,
        message
      }
    }
  }
`;

const GET_MARGIN_BY_ID = gql`
  query getGasStationMarginById($id: String!) {
    getGasStationMarginById(id: $id) {
      ... on GasStationMargin {
        id
        marginType
        margin
        points
        product {
          name
          id
        }
        gasStation {
          name
          id
        }
      }
      ... on GeneralError {
        code
        message
      }
    }
  }
`;

const columnHelper = createColumnHelper<BaseMargin>();

export default function MarginsScreen() {
  const navigate = useNavigate();
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

    const handleSearchChange = useCallback((newSearch: string) => {
      setSearch(newSearch);
    }, []);


  const { data, loading, error } = useQuery<{
    gasStationsMargin: BasePagination<BaseMargin> | InputValidationError[] | GeneralError;
  }>(GET_MARGINS, {
    variables: { ...cursor,
      search
     },
    onCompleted: () => {
      client.resetStore();

    },
  });


  const handleNextPage = () => {
    if (data?.gasStationsMargin && "__typename" in data?.gasStationsMargin) {
      const marginsPaginated = data.gasStationsMargin as BasePagination<BaseMargin>;
      setCursor({ nextCursor: marginsPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.gasStationsMargin && "__typename" in data?.gasStationsMargin) {
      const marginsPaginated = data.gasStationsMargin as BasePagination<BaseMargin>;
      setCursor({ prevCursor: marginsPaginated.prevCursor ?? null });
    }
  };

  const [loadingProduct, setLoadingProduct] = useState(false);

  const handleEdit = useCallback(async (marginId: string) => {
    setLoadingProduct(true);
    try {
      const { data} = await client.query({
        query: GET_MARGIN_BY_ID,
        variables: { id: marginId },
        fetchPolicy: "network-only",
      });

      if (data?.getGasStationMarginById?.__typename === "GasStationMargin") {

         navigate("/protected/margins/create", { state: { margin: data.getGasStationMarginById } });
      } else if (data?.getGasStationMarginById?.__typename === "GeneralError") {
      }
    } catch (error) {

    } finally {
      setLoadingProduct(false);
    }
  }, [navigate, client]);

  const handleCreate = useCallback(() => {
    navigate("/protected/margins/create");
  }, [navigate]);

  const columns = useMemo(
    () => [
      columnHelper.accessor("marginType", {
        cell: (info) => {
          const value = info.getValue();
          return value === "by_liter" ? "Por Litro" : "Por Margen";
        },
        header: () => <span>Tipo de Margen</span>,
      }),
      columnHelper.accessor("margin", {
        cell: (info) => info.getValue(),
        header: () => <span>Margen (%)</span>,
      }),
      columnHelper.accessor("points", {
        cell: (info) => info.getValue(),
        header: () => <span>Puntos</span>,
      }),
      columnHelper.accessor("product.name", { 
        cell: (info) => info.getValue(),
        header: () => <span>Producto</span>,
      }),
      columnHelper.accessor("gasStation.name", {  
        cell: (info) => info.getValue(),
        header: () => <span>Estación</span>,
      }),
      columnHelper.display({
        id: "actions",
        header: "Acciones",
        cell: ({ row }) => (
          <div className="flex space-x-2">
            <button className="btn btn-sm btn-outline" onClick={() => handleEdit(row.original.id)}>
              <FaPencilAlt />
            </button>
          </div>
        ),
      }),
    ],
    [handleEdit]
  );
  

  if (loading) return <LoadingSpinner />;

  if (error) return <p>Error al cargar los datos</p>;

  if (data?.gasStationsMargin && "__typename" in data?.gasStationsMargin) {

    if (data.gasStationsMargin.__typename === "InputValidationError") {

      return <p>Error en la validación de los campos</p>;
    } else if (data.gasStationsMargin.__typename === "GeneralError") {

      return <p>Error: {(data.gasStationsMargin as GeneralError).message}</p>;
    } else {
      const margins = (data.gasStationsMargin as BasePagination<BaseMargin>)?.items ?? [];
      const marginsPaginated = data.gasStationsMargin as BasePagination<BaseMargin>;
      const totalMargins = marginsPaginated.total ?? 0;

      return (
        <div>

          <div className="flex flex-col sm:flex-row items-center sm:justify-between gap-2 w-full">
            <InputSearch onSearchChange={handleSearchChange} />
            <button
              className="btn btn-primary w-full sm:w-auto ml-4 mr-4 sm:ml-0 sm:mr-4"
              onClick={handleCreate}
            >
            Agregar Margen
            </button>
          </div>
          

          <div className="relative">
            {loadingProduct ? (
              <div className="absolute inset-0 flex items-center justify-center bg-white opacity-75 z-10">
                <LoadingSpinner />
              </div>
            ) : (
              <Table
                data={margins}
                columns={columns}
                total={totalMargins}
                paginationComponent={
                  <Pagination
                    onNextPage={handleNextPage}
                    onPreviousPage={handlePreviousPage}
                    hasNext={!!marginsPaginated.nextCursor}
                    hasPrevious={!!marginsPaginated.prevCursor}
                  />
                }
              />
            )}
          </div>

        </div>
      );
    }
  }

  return <p>No hay registros de márgenes disponibles</p>;
}
