import { createBrowserRouter } from "react-router-dom";
import Layout from "./components/Layout";
import Protected from "./components/Protected";
import Unprotected from "./components/Unprotected";
import GasPump from "./pages/GasPump";
import GasStation from "./pages/GasStation";
import Login from "./pages/Login";
import Payment from "./pages/Payment";
import SettingPage from "./pages/Setting";
import SynchronizationEvents from "./pages/Synchronization";
import SynchronizationDetails from "./pages/SynchronizationDetails";
import UserDetail from "./pages/UserDetail";
import UserList from "./pages/UserList";
import CampaignList from "./pages/Campaign";
import CampaignDetail from "./pages/CampaignDetail";
import ElegibilityLevelPage from "./pages/Elegibility";
import CustomerLevelPage from "./pages/CustomerLevel";

export const routes = [
  {
    path: "/login",
    element: (
      <Unprotected>
        <Login />
      </Unprotected>
    ),
  },
  {
    path: "/protected",
    element: (
      <Protected>
        <Layout />
      </Protected>
    ),
    children: [
      {
        path: "/protected",
        element: <div>Dashboard</div>,
      },
      {
        path: "/protected/gas-stations",
        element: <GasStation />,
      },
      {
        path: "/protected/gas-pumps",
        element: <GasPump />,
      },
      {
        path: "/protected/payments",
        element: <Payment />,
      },
      {
        path: "/protected/synchronizations",
        element: <SynchronizationEvents />,
      },
      {
        path: "/protected/synchronizations/:id/details",
        element: <SynchronizationDetails />,
      },
      {
        path: "/protected/users",
        element: <UserList />,
      },
      {
        path: "/protected/users/new",
        element: <UserDetail />,
      },
      {
        path: "/protected/users/:id/edit",
        element: <UserDetail />,
      },
      {
        path: "/protected/settings",
        element: <SettingPage />,
      },
      {
        path: "/protected/campaigns",
        element: <CampaignList />,
      },
      {
        path: "/protected/campaigns/new",
        element: <CampaignDetail />,
      },
      {
        path: "/protected/campaigns/:id/edit",
        element: <CampaignDetail />,
      },
      {
        path: "/protected/elegibility/levels",
        element: <ElegibilityLevelPage />,
      },
      {
        path: "/protected/elegibility/customers/levels",
        element: <CustomerLevelPage />,
      },
    ],
  },
];

export default createBrowserRouter(routes, {
  basename: import.meta.env.VITE_BASE_URL,
});
