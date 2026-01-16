import {
  ApolloClient,
  InMemoryCache,
  HttpLink,
  ApolloLink,
} from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import { store } from "./store";
import { RootState } from "./store";

const apiUrl = import.meta.env.VITE_API_URL;

const httpLink = new HttpLink({
  uri: `${apiUrl}/graphql`,
});

const authLink = setContext((_, { headers }) => {
  const state: RootState = store.getState();
  const accessToken = state.auth.accessToken;

  return {
    headers: {
      ...headers,
      Authorization: accessToken ? `Bearer ${accessToken}` : "",
    },
  };
});

export const client = new ApolloClient({
  link: ApolloLink.from([authLink, httpLink]),
  cache: new InMemoryCache(),
});
