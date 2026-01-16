import { gql, useQuery } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import ActiveCheckbox from "../components/ActiveCheckbox"; 
import { BaseUser } from "../models/base_user";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import LoadingSpinner from "../components/LoadingSpinner";
import { useMemo } from "react";
import Pagination from "../components/Pagination";
import { useApolloClient } from "@apollo/client";
import { FaPencilAlt } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import InputSearch from "../components/InputSearch"; 
import { useSearchParams } from "react-router-dom";
import React, { useState, useEffect, useCallback } from "react";


const GET_USERS = gql`
  query Users($nextCursor: String, $prevCursor: String, $search: String) {
    users(
        pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor },
      search: $search
    ) {
      ...on InputValidationError {
        errors {
          field
          type
          message
        }
      }
      ...on UserPagination {
        nextCursor
        prevCursor
        total
        items {
          id,
          firstName,
          lastName,
          email,
          isActive
        }
      }
      ...on GeneralError {
        code
        message
      }
    }
  }
`;


const GET_USER_BY_ID = gql`
  query getUserById($id: String!) {
      getUserById(id: $id) {
      ... on User {
        id,
        firstName,
        lastName,
        email,
        isActive
      }
      ...on GeneralError {
        code,
        message
      }
    }
  }
`;

const columnHelper = createColumnHelper<BaseUser>();

export default function UsersScreen() {
  const navigate = useNavigate();
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

  const handleSearchChange = useCallback((newSearch: string) => {
    setSearch(newSearch);
  }, []);

  
  const { data, loading, error } = useQuery<{
    users: BasePagination<BaseUser> | InputValidationError[] | GeneralError;
  }>(GET_USERS, {
    variables: {
      ...cursor,
      search
    },
    onCompleted: (data) => {
      client.resetStore();
    },
  });

const handleNextPage = () => {
  if (data?.users && '__typename' in data?.users) {
    const customersPaginated = data.users as BasePagination<BaseUser>;
    setCursor({ nextCursor: customersPaginated.nextCursor ?? null });
  }
};

const handlePreviousPage = () => {
  if (data?.users && '__typename' in data?.users) {
    const customersPaginated = data.users as BasePagination<BaseUser>;
    setCursor({ prevCursor: customersPaginated.prevCursor ?? null });
  }
};


  const [loadingProduct, setLoadingProduct] = useState(false);

  const handleEdit = useCallback(async (userId: string) => {
    setLoadingProduct(true);
    try {
      const { data } = await client.query({
        query: GET_USER_BY_ID,
        variables: { id: userId },
        fetchPolicy: "network-only",
      });

      if (data?.getUserById?.__typename === "User") {
        navigate("/protected/users/create", { state: { product: data.getUserById } });
      } else if (data?.getUserById?.__typename === "GeneralError") {

      }
    } catch (error) {

    } finally {
      setLoadingProduct(false);
    }
  }, [navigate, client]);

  const handleCreate = useCallback(() => {
    navigate("/protected/users/create");
  }, [navigate]);

  const columns = useMemo(
    () => [
  columnHelper.accessor("firstName", {
    cell: (info) => info.getValue(),
    header: () => <span>Nombres</span>,
  }),
  columnHelper.accessor("lastName", {
    cell: (info) => info.getValue(),
    header: () => <span>Apellidos</span>,
  }),
  columnHelper.accessor("email", {
    cell: (info) => info.getValue(),
    header: () => <span>Email</span>,
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
        <button className="btn btn-sm btn-outline" onClick={() => handleEdit(row.original.id ?? "")}>
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


if (data?.users && '__typename' in data?.users) {
  if (data.users.__typename === 'InputValidationError') {
    return <p>Error en la validación de los campos</p>;
  } else if (data.users.__typename === 'GeneralError') {
    return <p>Error: {(data.users as GeneralError).message}</p>;
  } else {
    const users = (data.users as BasePagination<BaseUser>)?.items ?? [];
    const usersPaginated = data.users as BasePagination<BaseUser>;
    const usersTotal = usersPaginated.total ?? 0;

    return (
      <div>

  <div className="flex flex-col sm:flex-row items-center sm:justify-between gap-2 w-full">
    <InputSearch onSearchChange={handleSearchChange} />
    <button
      className="btn btn-primary w-full sm:w-auto ml-4 mr-4 sm:ml-0 sm:mr-4"
      onClick={handleCreate}
    >
      Agregar Usuario
    </button>
  </div>


        <div className="relative">
          {loadingProduct ? (
            <div className="absolute inset-0 flex items-center justify-center bg-white opacity-75 z-10">
          <LoadingSpinner />
        </div>
              ) : (
        <Table
          data={users}
          columns={columns}
          total={usersTotal}
          paginationComponent={<Pagination
            onNextPage={handleNextPage}
            onPreviousPage={handlePreviousPage}
            hasNext={!!usersPaginated.nextCursor}
            hasPrevious={!!usersPaginated.prevCursor}
            />
          }
        />
      )}
      </div>


      </div>
    );
  }
}

return <p>No hay registros de clientes disponibles</p>;
}



