import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { gql, useMutation, useApolloClient } from "@apollo/client";
import FormProducts, { Product } from "../components/FormProducts";
import { BaseLevel } from "../models/base_level";
import FormLevels from "../components/FormLevels";

const UPDATE_LEVEL = gql`
  mutation updateLevel($id: String!, $input: UpdateLevelBody!) {
    updateLevel(id: $id, body: $input) {
      ...on Level {
        id
        minPoints
        name
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


const ADD_LEVEL = gql`
  mutation addLevel($input: AddlevelBody!) {
      addLevel(body: $input) {
        ...on Level {
          id
          minPoints
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



export default function LevelsEditCreateScreen() {
  const location = useLocation();
  const navigate = useNavigate();
  const client = useApolloClient();
  
  
  const level = location.state?.level as BaseLevel | undefined;

  const [updateLevel] = useMutation(UPDATE_LEVEL);
  const [addLevel] = useMutation(ADD_LEVEL); 

  const handleSubmit = async (data: BaseLevel) => {
    try {
      let result;
      if (level) {
        const { data: updated } = await updateLevel({
          variables: { 
            id: level.id, 
            input: { name: data.name, minPoints: data.minPoints, isActive: data.isActive } 
          },
        });
        result = updated?.updateLevel;
      } else {
        const { data: added } = await addLevel({
          variables: { input: { name: data.name, minPoints: data.minPoints, isActive: data.isActive } },
        });
        result = added?.addLevel;
      }

      if (result) {
        if (result.__typename === "Level") {
          await client.resetStore();
          navigate("/protected/levels");
        } else if (result.__typename === "InputValidationError") {

        } else if (result.__typename === "GeneralError") {

        }
      }
    } catch (error) {

    }
  };

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold">{level ? "Editar Nivel" : "Crear Nivel"}</h1>
      <FormLevels level={level} onSubmit={handleSubmit} />
    </div>
  );
}
