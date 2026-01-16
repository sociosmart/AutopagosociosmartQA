import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { gql, useMutation, useApolloClient } from "@apollo/client";
import FormProducts, { Product } from "../components/FormProducts";
import { BaseUser } from "../models/base_user";
import FormUsers from "../components/FromUsers";

const UPDATE_USER = gql`
  mutation updateUser($id: String!, $input: UpdateUserBody!) {
    updateUser(id: $id, body: $input) {
      ...on User {
        id
        firstName
        lastName
        email
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


const ADD_USER = gql`
  mutation AddUser($input: AddUserBody!) {
    addUser(data: $input) {
      ...on User {
        id
        firstName
        lastName
        email
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


export default function UsersEditCreateScreen() {
  const location = useLocation();
  const navigate = useNavigate();
  const client = useApolloClient(); 
  const product = location.state?.product as BaseUser | undefined;

  const [updateUser] = useMutation(UPDATE_USER);
  const [addUser] = useMutation(ADD_USER);


  const handleSubmit = async (data: BaseUser) => {
    try {
      let result;
      if (product) {
        const { data: updated } = await updateUser({
          variables: { 
            id: product.id, 
            input: { 
                firstName: data.firstName,
                lastName: data.lastName,
                email: data.email,
                password: data.password,
                isActive: data.isActive
            } },
        });
        result = updated?.updateUser;

      } else {
        const { data: added } = await addUser({
          variables: { 
            input: { 
                firstName: data.firstName,
                lastName: data.lastName,
                email: data.email,
                password: data.password,
                isActive: data.isActive
            } },
        });
        result = added?.addUser;

      }

      if (result) {
        if (result.__typename === "User") {
          await client.resetStore();
          navigate("/protected/users");

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
        {product ? "Editar Usuario" : "Crear Usuario"}
      </h1>
      <FormUsers user={product} onSubmit={handleSubmit} />
    </div>
  );
}
