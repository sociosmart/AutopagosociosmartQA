import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { ApolloProvider } from "@apollo/client";
import { client } from "./client";
import { store } from './store'
import { Provider } from 'react-redux'


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Provider store={store}>
      <ApolloProvider client={client}>
        <App />
      </ApolloProvider>
    </Provider>,
  </StrictMode>
)


