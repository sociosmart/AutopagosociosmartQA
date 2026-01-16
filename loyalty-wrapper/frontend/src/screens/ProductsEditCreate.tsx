import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { gql, useMutation, useApolloClient } from "@apollo/client";
import FormProducts, { Product } from "../components/FormProducts";

const UPDATE_PRODUCT = gql`
  mutation updateProduct($id: String!, $input: UpdateProductBody!) {
    updateProduct(id: $id, data: $input) {
      ...on Product {
        id
        name
        codename
        isActive
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

const ADD_PRODUCT = gql`
  mutation AddProduct($input: AddProductBody!) {
    addProduct(data: $input) {
      ...on Product {
        id
        name
        codename
        isActive
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


export default function ProductsEditCreateScreen() {
  const location = useLocation();
  const navigate = useNavigate();
  const client = useApolloClient(); 
  const product = location.state?.product as Product | undefined;

  const [updateProduct] = useMutation(UPDATE_PRODUCT);
  const [addProduct] = useMutation(ADD_PRODUCT);


  const handleSubmit = async (data: Product) => {
    try {
      let result;
      if (product) {
        const { data: updated } = await updateProduct({
          variables: { id: product.id, input: { name: data.name, codename: data.codename,  isActive: data.isActive } },
        });
        result = updated?.updateProduct;

      } else {
        const { data: added } = await addProduct({
          variables: { input: { name: data.name, codename: data.codename, isActive: data.isActive } },
        });
        result = added?.addProduct;

      }

      if (result) {
        if (result.__typename === "Product") {
          await client.resetStore();
          navigate("/protected/products");

        } else if (result.__typename === "InputValidationError") {


        } else if (result.__typename === "GeneralError") {


        }
      } 
    } catch (error) {


    }
  };

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">
        {product ? "Editar Producto" : "Crear Producto"}
      </h1>
      <FormProducts product={product} onSubmit={handleSubmit} />
    </div>
  );
}
