import { gql, useQuery, useApolloClient } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import { BaseBenefitTicket } from "../models/base_benefit_ticket";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import Table from "../components/Table";
import LoadingSpinner from "../components/LoadingSpinner";
import Pagination from "../components/Pagination";
import InputSearch from "../components/InputSearch"; 
import { useState, useCallback } from "react";
import { format } from "date-fns";
import ActiveCheckbox from "../components/ActiveCheckbox";

const GET_BENEFITS_TICKETS = gql`
  query BenefitsTickets($nextCursor: String, $prevCursor: String, $search: String) {
    benefitsTickets(
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
    ... on BenefitTicketPagination {
      nextCursor
      prevCursor
      total
      items {
      	id
        createdAt
      	updatedAt
        customer{
          name
          lastName
        }
        benefitGenerated{
            name
        }
        startDate
        endDate
        redeemed
      }
    }
      ... on GeneralError {
        code
        message
      }
    }
  }
`;


const columnHelper = createColumnHelper<BaseBenefitTicket>();

const columns = [
  columnHelper.accessor("customer.name", {
    cell: (info) => info.getValue(),
    header: () => <span>Nombre</span>,
  }),
  columnHelper.accessor("customer.lastName", {
    cell: (info) => info.getValue(),
    header: () => <span>Apellidos</span>,
  }),
  columnHelper.accessor("benefitGenerated.name", {
    cell: (info) => info.getValue(),
    header: () => <span>Beneficio</span>,
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
  columnHelper.accessor("redeemed", {
    cell: (info) => <ActiveCheckbox isActive={info.getValue()} />,
    header: () => <span>Activo</span>,
  }),
];

export default function BenefitsTicketsScreen() {
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

  const handleSearchChange = useCallback((newSearch: string) => {
    setSearch(newSearch);
  }, []);


  const { data, loading, error } = useQuery<{
    benefitsTickets: BasePagination<BaseBenefitTicket> | InputValidationError[] | GeneralError;
  }>(GET_BENEFITS_TICKETS, {
    variables: {
      ...cursor,
      search
    },
    onCompleted: (data) => {
      client.resetStore();
    },
  });

  const handleNextPage = () => {
    if (data?.benefitsTickets && '__typename' in data?.benefitsTickets) {
      const benefitsTicketsPaginated = data.benefitsTickets as BasePagination<BaseBenefitTicket>;
      setCursor({ nextCursor: benefitsTicketsPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.benefitsTickets && '__typename' in data?.benefitsTickets) {
      const benefitsTicketsPaginated = data.benefitsTickets as BasePagination<BaseBenefitTicket>;
      setCursor({ prevCursor: benefitsTicketsPaginated.prevCursor ?? null });
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <p>Error al cargar los datos</p>;

  if (data?.benefitsTickets && '__typename' in data?.benefitsTickets) {
    if (data.benefitsTickets.__typename === 'InputValidationError') {
      return <p>Error en la validación de los campos</p>;
    } else if (data.benefitsTickets.__typename === 'GeneralError') {
      return <p>Error: {(data.benefitsTickets as GeneralError).message}</p>;
    } else {
      const benefitsTicketsPaginated = data.benefitsTickets as BasePagination<BaseBenefitTicket>;
      const benefitsTickets = benefitsTicketsPaginated?.items ?? [];
      const totalBenefitsTickets = benefitsTicketsPaginated.total ?? 0; //  total

      return (
        <div className="w-full">
          <InputSearch onSearchChange={handleSearchChange} />
          <Table
            data={benefitsTickets}
            columns={columns}
            total={totalBenefitsTickets} 
            paginationComponent={
              <Pagination
                onNextPage={handleNextPage}
                onPreviousPage={handlePreviousPage}
                hasNext={!!benefitsTicketsPaginated.nextCursor}
                hasPrevious={!!benefitsTicketsPaginated.prevCursor}
              />
            }
          />
        </div>
      );
    }
  }

  return <p>No hay registros de clientes disponibles</p>;
}

