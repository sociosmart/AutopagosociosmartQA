
import { useState, useMemo, useCallback } from "react";
import { gql, useQuery, useApolloClient } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import ActiveCheckbox from "../components/ActiveCheckbox";
import LoadingSpinner from "../components/LoadingSpinner";
import Pagination from "../components/Pagination";
import { useNavigate } from "react-router-dom";
import { BaseBenefitGenerated } from "../models/base_benefit_generated";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import { FaPencilAlt } from "react-icons/fa";
import { format } from "date-fns";
import InputSearch from "../components/InputSearch"; 

const GET_BENEFITS = gql`
  query BenefitsGenerated($nextCursor: String, $prevCursor: String, $search: String) {
    benefitsGenerated(
      pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor },
      search: $search) {
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
    ... on BenefitGeneratedPagination {
      nextCursor
      prevCursor
      total
      items {
      id
      name
      level {
        name
      }
      stock
      discount
      numTimes
      frequency
      type
      createdAt
      updatedAt
      isActive
      externalProductId
        startDate
        endDate
        benefit{
          name
        }
        dependency
        minAmount
      }
    }
    }
  }
`;

const GET_BENEFIT_BY_ID = gql`
  query BenefitGeneratedById($id: String!) {
    getBenefitGeneratedById(id: $id) {
    ... on BenefitGenerated {
      id
      name
      level {
        name
        id
      }
      stock
      discount
      numTimes
      frequency
      type
      createdAt
      updatedAt
      isActive
      externalProductId
        startDate
        endDate
        benefit{
          name
        }
      dependency
      minAmount  
    }
    }
  }
`;



const columnHelper = createColumnHelper<BaseBenefitGenerated>();

export default function BenefitsGeneratedScreen() {
  const navigate = useNavigate();
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

  const handleSearchChange = useCallback((newSearch: string) => {
    setSearch(newSearch);
  }, []);

  const { data, loading, error } = useQuery<{
    benefitsGenerated: BasePagination<BaseBenefitGenerated> | InputValidationError[] | GeneralError;
  }>(GET_BENEFITS, {
    variables: { ...cursor, search },
    onCompleted: () => {
      client.resetStore();
    },
  });

  const handleNextPage = () => {
    if (data?.benefitsGenerated && "__typename" in data?.benefitsGenerated) {
      const levelsPaginated = data.benefitsGenerated as BasePagination<BaseBenefitGenerated>;
      setCursor({ nextCursor: levelsPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.benefitsGenerated && "__typename" in data?.benefitsGenerated) {
      const levelsPaginated = data.benefitsGenerated as BasePagination<BaseBenefitGenerated>;
      setCursor({ prevCursor: levelsPaginated.prevCursor ?? null });
    }
  };

  const [loadingProduct, setLoadingProduct] = useState(false);

  const handleEdit = useCallback(async (Id: string) => {
    setLoadingProduct(true);
    try {
      const { data } = await client.query({
        query: GET_BENEFIT_BY_ID,
        variables: { id: Id },
        fetchPolicy: "network-only",
      });

      if (data?.getBenefitGeneratedById?.__typename === "BenefitGenerated") {
        navigate("/protected/benefits-generated/create", { state: { level: data.getBenefitGeneratedById } });
      } else if (data?.getBenefitGeneratedById?.__typename === "GeneralError") {

      }
    } catch (error) {

    } finally {
      setLoadingProduct(false);
    }
  }, [navigate, client]);

  // const handleCreate = useCallback(() => {
  //   navigate("/protected/benefits-generated/create");
  // }, [navigate]);

  const columns = useMemo(
    () => [
      columnHelper.accessor("name", {
        cell: (info) => info.getValue(),
        header: () => <span>Nombre</span>,
      }),
      columnHelper.accessor("level.name", {
        cell: (info) => info.getValue(),
        header: () => <span>Nivel</span>,
      }),
      columnHelper.accessor("type", {
        cell: (info) => {
          const typeValue = info.getValue();
          const typeMapping: Record<string, string> = {
            physical: "Fisico",
            digital: "Digital",
            gas: "Combustible",
            periferics: "Perifericos",
          };
          return typeMapping[typeValue] || typeValue;
        },
        header: () => <span>Tipo</span>,
      }),
      columnHelper.accessor("frequency", {
        cell: (info) => {
          const freqValue = info.getValue();
          const frequencyMapping: Record<string, string> = {
            n_times: "N veces",
            hourly: "Cada Hora",
            daily: "Diariamente",
            weekly: "Semanalmente",
            monthly: "Mensualmente",
            always: "Siempre",
          };
          return frequencyMapping[freqValue] || freqValue;
        },
        header: () => <span>Frecuencia</span>,
      }),
      columnHelper.accessor("externalProductId", {
        cell: (info) => info.getValue(),
        header: () => <span>Producto Ext.</span>,
      }),
      columnHelper.accessor("discount", {
        cell: (info) => info.getValue(),
        header: () => <span>Descuento</span>,
      }),
      columnHelper.accessor("numTimes", {
        cell: (info) => info.getValue(),
        header: () => <span>N Veces</span>,
      }),
      columnHelper.accessor("stock", {
        cell: (info) => info.getValue(),
        header: () => <span>Stock</span>,
      }),
      columnHelper.accessor("dependency", {
        cell: (info) => <ActiveCheckbox isActive={info.getValue()} />,
        header: () => <span>Es Dependiente</span>,
      }),
      columnHelper.accessor("minAmount", {
        cell: (info) => info.getValue(),
        header: () => <span>Cantidad Minima</span>,
      }),
      columnHelper.accessor("startDate", {
        cell: (info) => {
          const dateValue = info.getValue();
          return dateValue
            ? format(new Date(dateValue), "dd/MM/yyyy hh:mm aaaa")
            : "N/A";
        },
        header: () => <span>Fecha de Inicio</span>,
      }),
      columnHelper.accessor("endDate", {
        cell: (info) => {
          const dateValue = info.getValue();
          return dateValue
            ? format(new Date(dateValue), "dd/MM/yyyy hh:mm aaaa")
            : "N/A";
        },
        header: () => <span>Fecha de Termino</span>,
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
            <button
              className="btn btn-sm btn-outline"
              onClick={() => handleEdit(row.original.id ?? "")}
            >
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
  
  if (data?.benefitsGenerated && "__typename" in data?.benefitsGenerated) {
    if (data.benefitsGenerated.__typename === "InputValidationError") {
      return <p>Error en la validación de los campos</p>;
    } else if (data.benefitsGenerated.__typename === "GeneralError") {
      return <p>Error: {(data.benefitsGenerated as GeneralError).message}</p>;
    } else {
      const levels = (data.benefitsGenerated as BasePagination<BaseBenefitGenerated>)?.items ?? [];
      const levelsPaginated = data.benefitsGenerated as BasePagination<BaseBenefitGenerated>;
      const totalLevels = levelsPaginated.total ?? 0;


      return (
        <div>
          <InputSearch onSearchChange={handleSearchChange} />
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
