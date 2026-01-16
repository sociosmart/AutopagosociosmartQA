import { useLocation, useNavigate } from "react-router-dom";
import { gql, useMutation, useApolloClient,useQuery } from "@apollo/client";
import { BaseBenefitGenerated } from "../models/base_benefit_generated";
import FormBenefitsGenerated from "../components/FormBenefitsGenerated";


const UPDATE_BENEFIT_GENERATED = gql`
  mutation UpdateBenefitGenerated($id: String!, $body: UpdateGeneratedBenefit!) {
    updateGeneratedBenefit(id: $id, body: $body) {
      ...on BenefitGenerated {
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


export default function BenefitsGeneratedEditCreateScreen() {
    const location = useLocation();
    const navigate = useNavigate();
    const client = useApolloClient();
    
    const level = location.state?.level as BaseBenefitGenerated | undefined;
    const [updateBenefit] = useMutation(UPDATE_BENEFIT_GENERATED);
  

    const { data: levelsData, loading: levelsLoading, error: levelsError } = useQuery(GET_LEVELS);
  
    if (levelsLoading) return <p>Cargando niveles...</p>;
    if (levelsError) return <p>Error al cargar los niveles: {levelsError.message}</p>;
  
    const handleSubmit = async (data: Omit<BaseBenefitGenerated, "createdAt" | "updatedAt" | "benefit">) => {
        try {
          let result;
      
          if (level) {
            const { data: updated } = await updateBenefit({
              variables: {
                id: level.id, 
                body: {
                  stock: data.stock,
                  isActive: data.isActive,
                }
              }
            });
            result = updated?.updateGeneratedBenefit;
          } 

          if (result) {
            if (result.__typename === "BenefitGenerated") {
              await client.resetStore();
              navigate("/protected/benefits-generated");
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
        <FormBenefitsGenerated 
          benefit={level} 
          onSubmit={handleSubmit} 
          levels={levelsData?.levels?.items || []} 
        />
      </div>
    );
  }
  
