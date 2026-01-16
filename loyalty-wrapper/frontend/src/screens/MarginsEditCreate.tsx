import { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { gql, useMutation, useQuery, useApolloClient } from "@apollo/client";
import FormMargins from "../components/FormMargins";
import { BaseMargin } from "../models/base_margin";


const UPDATE_MARGIN = gql`
  mutation updateMargin($id: String!, $input: UpdateGasStationMargin!) {
    updateGasStationMargin(id: $id, data: $input) {
      ...on GasStationMargin {
        id
        marginType
        margin
        points
        gasStation {
          id
        }
        product {
          id
        }
      }
      ...on InputValidationError {
        errors {
          field
          type
          message
        }
      }
      ...on GeneralError {
        code
        message
      }
    }
  }
`;

const ADD_MARGIN = gql`
  mutation addMargin($input: AddGasStationMargin!) {
    addGasStationMargin(data: $input) {
      ...on GasStationMargin {
        marginType
        margin
        points
        product {
          id
        }
        gasStation {
          id
        }
      }
      ...on InputValidationError {
        errors {
          field
          type
          message
        }
      }
      ...on GeneralError {
        code
        message
      }
    }
  }
`;

const GET_STATIONS = gql`
  query GasStation {
    gasStations(pagination: {limit: 100}) {
      ...on InputValidationError {
        errors {
          field,
          type,
          message
        }
      },
      ...on GasStationPagination {
        total,
        items {
          id,
          name,
          externalId,
          crePermission,
          latitude,
          longitude
        }
      },
      ...on GeneralError {
        code,
        message
      }
    }
  }
`;

const GET_PRODUCTS = gql`
  query Products {
    products(pagination: {limit: 100}) {
      ...on InputValidationError {
        errors {
          field,
          type,
          message
        }
      },
      ...on ProductsPagination {
        total
        items {
          id,
          name,
          isActive
        }
      },
      ...on GeneralError {
        code,
        message
      }
    }
  }
`;



export default function MarginsEditCreateScreen() {
  const location = useLocation();
  const navigate = useNavigate();
  const client = useApolloClient();
  const margin = location.state?.margin as BaseMargin | undefined;

  const [updateMargin] = useMutation(UPDATE_MARGIN);
  const [addMargin] = useMutation(ADD_MARGIN);
  const [error, setError] = useState(""); 

  const { data: productsData, loading: productsLoading } = useQuery(GET_PRODUCTS);
  const { data: stationsData, loading: stationsLoading } = useQuery(GET_STATIONS);

  const handleSubmit = async (data: BaseMargin) => {
    try {
      let result;
      if (margin) {
        const { data: updated } = await updateMargin({
          variables: { id: margin.id, input: { marginType: data.marginType, margin: data.margin, points: data.points, product: data.product, gasStation: data.gasStation } },

        });
        result = updated?.updateGasStationMargin;

      } else {
        const { data: added } = await addMargin({
          variables: { input: { marginType: data.marginType, margin: data.margin, points: data.points, product: data.product, gasStation: data.gasStation } },
        });

        result = added?.addGasStationMargin;

      }

      if (result) {
        if (result.__typename === "GasStationMargin") {

          await client.resetStore();
          navigate("/protected/margins");
        } else if (result.__typename === "InputValidationError") {

        } else if (result.__typename === "GeneralError") {
            setError("La estacion y el producto seleccionados ya existen ");
        }
      }
    } catch (error) {

    }
  };

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">
        {margin ? "Editar Margen" : "Crear Margen"}
      </h1>


      {error && (
        <div role="alert"   className="alert alert-warning absolute top-1/2 left-1/2 transform -translate-x-24 -translate-y-38 z-50 shadow-lg"        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <span>{error}</span>
        </div>
      )}


      {!productsLoading && !stationsLoading ? (
        <FormMargins
          margin={margin}
          products={productsData?.products?.items || []}
          gasStations={stationsData?.gasStations?.items || []}
          onSubmit={handleSubmit}
        />
      ) : (
        <div>Loading...</div>
      )}
    </div>
  );
}