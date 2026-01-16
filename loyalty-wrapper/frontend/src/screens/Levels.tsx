import { useState, useMemo, useCallback } from "react";
import { gql, useQuery, useApolloClient } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import ActiveCheckbox from "../components/ActiveCheckbox";
import LoadingSpinner from "../components/LoadingSpinner";
import Pagination from "../components/Pagination";
import { useNavigate } from "react-router-dom";
import { BaseLevel } from "../models/base_level";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import { FaPencilAlt } from "react-icons/fa";
import React, { useEffect } from "react";
import InputSearch from "../components/InputSearch"; 
import { useSearchParams } from "react-router-dom";


const GET_LEVELS = gql`
  query Levels($nextCursor: String, $prevCursor: String, $search: String) {
    levels(
      pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor },
      search: $search) {
      ...on InputValidationError {
        errors {
          field,
          type,
          message
        }
      }
        ... on LevelPagination {
        nextCursor
        prevCursor
        total
        items {
            id
            name
            minPoints
            isActive
        }
        }
      ...on GeneralError {
        code,
        message
      }
    }
  }
`;

const GET_LEVEL_BY_ID = gql`
  query getLevelById($id: String!) {
    getLevelById(id: $id) {
      ...on Level {
        id
        name
        minPoints
        isActive
      }
      ...on GeneralError {
        code,
        message
      }
    }
  }
`;



const columnHelper = createColumnHelper<BaseLevel>();

export default function LevelsScreen() {
  const navigate = useNavigate();
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

    const handleSearchChange = useCallback((newSearch: string) => {
      setSearch(newSearch);
    }, []);


  const { data, loading, error } = useQuery<{
    levels: BasePagination<BaseLevel> | InputValidationError[] | GeneralError;
  }>(GET_LEVELS, {
    variables: { ...cursor,
      search
     },
    onCompleted: () => {
      client.resetStore();
    },
  });

  const handleNextPage = () => {
    if (data?.levels && "__typename" in data?.levels) {
      const levelsPaginated = data.levels as BasePagination<BaseLevel>;
      setCursor({ nextCursor: levelsPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.levels && "__typename" in data?.levels) {
      const levelsPaginated = data.levels as BasePagination<BaseLevel>;
      setCursor({ prevCursor: levelsPaginated.prevCursor ?? null });
    }
  };

  const [loadingProduct, setLoadingProduct] = useState(false);

  const handleEdit = useCallback(async (Id: string) => {
    setLoadingProduct(true);
    try {
      const { data } = await client.query({
        query: GET_LEVEL_BY_ID,
        variables: { id: Id },
        fetchPolicy: "network-only",
      });

      if (data?.getLevelById?.__typename === "Level") {
        navigate("/protected/levels/create", { state: { level: data.getLevelById } });
      } else if (data?.getProductById?.__typename === "GeneralError") {

      }
    } catch (error) {

    } finally {
      setLoadingProduct(false);
    }
  }, [navigate, client]);

  const handleCreate = useCallback(() => {
    navigate("/protected/levels/create");
  }, [navigate]);

  const columns = useMemo(
    () => [
      columnHelper.accessor("name", {
        cell: (info) => info.getValue(),
        header: () => <span>Nombre</span>,
      }),
        columnHelper.accessor("minPoints", {
            cell: (info) => info.getValue(),
            header: () => <span>Puntos Minimos</span>,
        }),
      columnHelper.accessor("isActive", {
        cell: (info) => <ActiveCheckbox isActive={info.getValue()} />,
        header: () => <span>Activo</span>,
      }),
      columnHelper.display({
        id: "actions",
        header: "Acciones",
        cell: ({ row }) => (
          <div className="flex space-x-2">
            <button className="btn btn-sm btn-outline" onClick={() => handleEdit(row.original.id?? "")}>
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
  
  if (data?.levels && "__typename" in data?.levels) {
    if (data.levels.__typename === "InputValidationError") {
      return <p>Error en la validación de los campos</p>;
    } else if (data.levels.__typename === "GeneralError") {
      return <p>Error: {(data.levels as GeneralError).message}</p>;
    } else {
      const levels = (data.levels as BasePagination<BaseLevel>)?.items ?? [];
      const levelsPaginated = data.levels as BasePagination<BaseLevel>;
      const totalLevels = levelsPaginated.total ?? 0;


      return (
        <div>

            <div className="flex flex-col sm:flex-row items-center sm:justify-between gap-2 w-full">
              <InputSearch onSearchChange={handleSearchChange} />
              <button
                className="btn btn-primary w-full sm:w-auto ml-4 mr-4 sm:ml-0 sm:mr-4"
                onClick={handleCreate}
              >
              Agregar Nivel
              </button>
            </div>

          <div className="relative">
            {loadingProduct ? (
              <div className="absolute inset-0 flex items-center justify-center bg-white opacity-75 z-10">
            <LoadingSpinner />
              </div>
              ) : (
                <Table
                  data={levels}
                  columns={columns}
                  total={totalLevels}
                  paginationComponent={
                    <Pagination
                      onNextPage={handleNextPage}
                      onPreviousPage={handlePreviousPage}
                      hasNext={!!levelsPaginated.nextCursor}
                      hasPrevious={!!levelsPaginated.prevCursor}
                    />
                  }
                />
              )}
          </div>

        </div>
      );
    }
  }
  
  return <p>No hay registros de productos disponibles</p>;
  
}
