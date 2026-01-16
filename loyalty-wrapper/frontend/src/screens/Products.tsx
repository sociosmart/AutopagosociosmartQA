import { useState, useMemo, useCallback } from "react";
import { gql, useQuery, useApolloClient } from "@apollo/client";
import { createColumnHelper } from "@tanstack/react-table";
import Table from "../components/Table";
import ActiveCheckbox from "../components/ActiveCheckbox";
import LoadingSpinner from "../components/LoadingSpinner";
import Pagination from "../components/Pagination";
import { useNavigate } from "react-router-dom";
import { BaseProduct } from "../models/base_product";
import { BasePagination } from "../models/base_pagination";
import { InputValidationError } from "../models/input_validation_error";
import { GeneralError } from "../models/general_error";
import { FaPencilAlt } from "react-icons/fa";
import InputSearch from "../components/InputSearch"; 



const GET_PRODUCTS = gql`
  query Products($nextCursor: String, $prevCursor: String, $search: String) {
    products(
      pagination: { limit: 100, nextCursor: $nextCursor, prevCursor: $prevCursor },
      search: $search) {
      ...on InputValidationError {
        errors {
          field,
          type,
          message
        }
      }
      ...on ProductsPagination {
        nextCursor,
        prevCursor,
        total,
        items {
          id,
          name,
          codename,
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

const GET_PRODUCT_BY_ID = gql`
  query GetProductById($id: String!) {
    getProductById(id: $id) {
      ...on Product {
        id
        name
        codename
        isActive
      }
      ...on GeneralError {
        code,
        message
      }
    }
  }
`;

const columnHelper = createColumnHelper<BaseProduct>();

export default function ProductsScreen() {
  const navigate = useNavigate();
  const [cursor, setCursor] = useState<object>({});
  const client = useApolloClient();
  const [search, setSearch] = useState<string>("");

    const handleSearchChange = useCallback((newSearch: string) => {
      setSearch(newSearch);
    }, []);

  const { data, loading, error } = useQuery<{
    products: BasePagination<BaseProduct> | InputValidationError[] | GeneralError;
  }>(GET_PRODUCTS, {
    variables: { 
      ...cursor,
      search
     },
    onCompleted: () => {
      client.resetStore();
    },
  });

  const handleNextPage = () => {
    if (data?.products && "__typename" in data?.products) {
      const productsPaginated = data.products as BasePagination<BaseProduct>;
      setCursor({ nextCursor: productsPaginated.nextCursor ?? null });
    }
  };

  const handlePreviousPage = () => {
    if (data?.products && "__typename" in data?.products) {
      const productsPaginated = data.products as BasePagination<BaseProduct>;
      setCursor({ prevCursor: productsPaginated.prevCursor ?? null });
    }
  };

  const [loadingProduct, setLoadingProduct] = useState(false);

  const handleEdit = useCallback(async (productId: string) => {
    setLoadingProduct(true);
    try {
      const { data } = await client.query({
        query: GET_PRODUCT_BY_ID,
        variables: { id: productId },
        fetchPolicy: "network-only",
      });

      if (data?.getProductById?.__typename === "Product") {
        navigate("/protected/products/create", { state: { product: data.getProductById } });
      } else if (data?.getProductById?.__typename === "GeneralError") {

      }
    } catch (error) {

    } finally {
      setLoadingProduct(false);
    }
  }, [navigate, client]);

  const handleCreate = useCallback(() => {
    navigate("/protected/products/create");
  }, [navigate]);

  const columns = useMemo(
    () => [
      columnHelper.accessor("name", {
        cell: (info) => info.getValue(),
        header: () => <span>Nombre</span>,
      }),
      columnHelper.accessor("codename", {
        cell: (info) => info.getValue(),
        header: () => <span>Codigo</span>,
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
  
  if (data?.products && "__typename" in data?.products) {
    if (data.products.__typename === "InputValidationError") {
      return <p>Error en la validación de los campos</p>;
    } else if (data.products.__typename === "GeneralError") {
      return <p>Error: {(data.products as GeneralError).message}</p>;
    } else {
      const products = (data.products as BasePagination<BaseProduct>)?.items ?? [];
      const productsPaginated = data.products as BasePagination<BaseProduct>;
      const totalProducts = productsPaginated.total ?? 0;


      return (
        <div>

          <div className="flex flex-col sm:flex-row items-center sm:justify-between gap-2 w-full">
            <InputSearch onSearchChange={handleSearchChange} />
            <button
              className="btn btn-primary w-full sm:w-auto ml-4 mr-4 sm:ml-0 sm:mr-4"
              onClick={handleCreate}
            >
            Agregar Producto
            </button>
          </div>
          

          <div className="relative">
            {loadingProduct ? (
              <div className="absolute inset-0 flex items-center justify-center bg-white opacity-75 z-10">
            <LoadingSpinner />
              </div>
              ) : (
                <Table
                  data={products}
                  columns={columns}
                  total={totalProducts}
                  paginationComponent={
                    <Pagination
                      onNextPage={handleNextPage}
                      onPreviousPage={handlePreviousPage}
                      hasNext={!!productsPaginated.nextCursor}
                      hasPrevious={!!productsPaginated.prevCursor}
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
