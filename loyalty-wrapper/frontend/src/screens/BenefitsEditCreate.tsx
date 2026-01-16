import { useLocation, useNavigate } from "react-router-dom";
import { gql, useMutation, useApolloClient,useQuery } from "@apollo/client";
import { BaseBenefit } from "../models/base_benefit";
import FormBenefits from "../components/FormBenefits";
import { useState } from "react";


const UPDATE_BENEFIT = gql`
  mutation updateBenefit($id: String!, $body: UpdateBenefit!) {
    updateBenefit(id: $id, body: $body) {
      ...on Benefit {
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
        externalProductId
        createdAt
        updatedAt
        isActive
        dependency
        minAmount
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




const ADD_BENEFIT = gql`
  mutation addBenefit($input: AddBenefit!) {
    addBenefit(body: $input) {
      ...on Benefit {
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
        externalProductId
        createdAt
        updatedAt
        isActive
        dependency
        minAmount
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


const GET_LEVELS = gql`
  query Levels {
    levels(pagination: { limit: 100}) {
      ...on InputValidationError {
        errors {
          field,
          type,
          message
        }
      }
        ... on LevelPagination {
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


export default function BenefitsEditCreateScreen() {
    const location = useLocation();
    const navigate = useNavigate();
    const client = useApolloClient();
    
    const level = location.state?.level as BaseBenefit | undefined;
  
    const [updateBenefit] = useMutation(UPDATE_BENEFIT);
    const [addBenefit] = useMutation(ADD_BENEFIT); 
  

    const { data: levelsData, loading: levelsLoading, error: levelsError } = useQuery(GET_LEVELS);
  
    if (levelsLoading) return <p>Cargando niveles...</p>;
    if (levelsError) return <p>Error al cargar los niveles: {levelsError.message}</p>;
  
    const handleSubmit = async (data: Omit<BaseBenefit, "createdAt" | "updatedAt">) => {
        try {
          let result;
      
          if (level) {
            const { data: updated } = await updateBenefit({
              variables: {
                id: level.id, 
                body: {
                  name: data.name,
                  type: data.type,
                  externalProductId: data.externalProductId,
                  frequency: data.frequency,
                  discount: data.discount,
                  stock: data.stock,
                  numTimes: data.numTimes,
                  isActive: data.isActive,
                  level: data.level.id,
                  dependency: data.dependency,
                  minAmount: data.minAmount
                }
              }
            });
            result = updated?.updateBenefit;
          } else {
            const { data: added } = await addBenefit({
              variables: {
                input: {
                  name: data.name,
                  type: data.type,
                  externalProductId: data.externalProductId,
                  frequency: data.frequency,
                  discount: data.discount,
                  stock: data.stock,
                  numTimes: data.numTimes,
                  isActive: data.isActive,
                  level: data.level.id,
                  dependency: data.dependency,
                  minAmount: data.minAmount
                }
              }
            });
            result = added?.addBenefit;
          }
      
          if (result) {
            if (result.__typename === "Benefit") {
              await client.resetStore();
              navigate("/protected/benefits");
            } else if (result.__typename === "InputValidationError") {

            } else if (result.__typename === "GeneralError") {

            }
          }
        } catch (error) {
          console.error("Error al guardar el beneficio:", error);
        }
      };
      
    
    return (
      <div className="p-4">
        <h1 className="text-xl font-bold mb-4">{level ? "Editar Beneficio" : "Crear Beneficio"}</h1>
        <FormBenefits 
          benefit={level} 
          onSubmit={handleSubmit} 
          levels={levelsData?.levels?.items || []} 
        />
      </div>
    );
  }
  
