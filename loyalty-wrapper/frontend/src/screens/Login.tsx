import FormLogin from "../components/LoginForm"; 
import { gql } from "@apollo/client";

export const POST_LOGIN = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      ... on LoginSuccess {
        accessToken
        refreshToken
      }
      ... on LoginError {
        message
        type
      }
    }
  }
`;




export default function LoginScreen() {
  return (
      <FormLogin />
  );
};

