import { gql, useQuery } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import { BaseGasStation } from "../models/base_gas_station";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import LoadingSpinner from "../components/LoadingSpinner";
import { useMemo } from "react";
import Pagination from "../components/Pagination";
import { useApolloClient } from "@apollo/client";
import { FaLocationArrow } from "react-icons/fa";
import React, { useState, useEffect, useCallback } from "react";
import InputSearch from "../components/InputSearch"; 
import { useSearchParams } from "react-router-dom";

const GET_GAS_STATIONS = gql`
  query GasStation($nextCursor: String, $prevCursor: String, $search: String) {
  gasStations(
    pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor },
        search: $search) {
    ...on InputValidationError {
      errors {
        field,
        type,
        message
      }
    },
    ...on  GasStationPagination {
      nextCursor,
      prevCursor,
      total,
      items {
        name,
        regularPrice,
        premiumPrice,
        dieselPrice,
        externalId,
        crePermission,
        latitude,
        longitude,
        city
      }
    },
    ...on GeneralError {
      code,
      message
    }
  }
}
`;

const handleMap = (latitude: string, longitude: string) => {
  const mapUrl = `https://maps.google.com/?q=${latitude},${longitude}`;
  window.open(mapUrl, "_blank"); 
};

export default function GasStationsScreen() {
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

  const handleSearchChange = useCallback((newSearch: string) => {
    setSearch(newSearch);
  }, []);

  const { data, loading, error } = useQuery<{
    gasStations: BasePagination<BaseGasStation> | InputValidationError[] | GeneralError;
  }>(GET_GAS_STATIONS, {
    variables: {
      ...cursor,
      search
    },
    onCompleted: (data) => {
      client.resetStore();
    },
  });

  const handleNextPage = () => {
    if (data?.gasStations && '__typename' in data?.gasStations) {
      const customersPaginated = data.gasStations as BasePagination<BaseGasStation>;
      setCursor({ nextCursor: customersPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.gasStations && '__typename' in data?.gasStations) {
      const customersPaginated = data.gasStations as BasePagination<BaseGasStation>;
      setCursor({ prevCursor: customersPaginated.prevCursor ?? null });
    }
  };

  const columnHelper = createColumnHelper<BaseGasStation>();

  const columns = useMemo(
    () => [
      columnHelper.accessor("name", {
        cell: (info) => info.getValue(),
        header: () => <span>Nombre</span>,
      }),
      columnHelper.accessor("externalId", {
        cell: (info) => info.getValue(),
        header: () => <span>Id Externo</span>,
      }),
      columnHelper.accessor("crePermission", {
        cell: (info) => info.getValue(),
        header: () => <span>Permiso CRE</span>,
      }),
      columnHelper.accessor("city", {
        cell: (info) => info.getValue(),
        header: () => <span>Ciudad</span>,
      }),
      columnHelper.accessor("regularPrice", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Regular</span>,
      }),
      columnHelper.accessor("premiumPrice", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Premier</span>,
      }),
      columnHelper.accessor("dieselPrice", {
        cell: (info) => info.getValue().toFixed(2),
        header: () => <span>Diesel</span>,
      }),
      columnHelper.display({
        id: "actions",
        header: "Acciones",
        cell: ({ row }) => (
          <div className="flex space-x-2">
            <button
              className="btn btn-sm btn-outline"
              onClick={() => handleMap(row.original.latitude, row.original.longitude)}
            >
              <FaLocationArrow />
            </button>
          </div>
        ),
      }),
    ],
    [] 
  );

  if (loading) return <LoadingSpinner />;
  if (error) return <p>Error al cargar los datos</p>;

  if (data?.gasStations && '__typename' in data?.gasStations) {
    if (data.gasStations.__typename === 'InputValidationError') {
      return <p>Error en la validación de los campos</p>;
    } else if (data.gasStations.__typename === 'GeneralError') {
      return <p>Error: {(data.gasStations as GeneralError).message}</p>;
    } else {
      const gasStations = (data.gasStations as BasePagination<BaseGasStation>)?.items ?? [];
      const gasStationsPaginated = data.gasStations as BasePagination<BaseGasStation>;
      const totalGasStations = gasStationsPaginated.total ?? 0;

      return (
        <div>
          <InputSearch onSearchChange={handleSearchChange} />
          <Table
            data={gasStations}
            columns={columns}
            total={totalGasStations}
            paginationComponent={<Pagination
              onNextPage={handleNextPage}
              onPreviousPage={handlePreviousPage}
              hasNext={!!gasStationsPaginated.nextCursor}
              hasPrevious={!!gasStationsPaginated.prevCursor}
              />}
          />
        </div>
      );
    }
  }

  return <p>No hay registros de clientes disponibles</p>;
}
