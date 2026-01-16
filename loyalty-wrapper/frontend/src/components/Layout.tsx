import { Outlet } from "react-router";
import NavBar from "./NavBar";
import SideBar, { SideBarNavigationItem,NavigationEntry } from "./SideBar";
import { useMemo, useEffect } from "react";
import { FaGasPump, FaUsers, FaSortAmountUpAlt, FaCreditCard, FaAlignJustify  } from "react-icons/fa";
import { GiCardExchange } from "react-icons/gi";
import { gql, useQuery, useApolloClient } from "@apollo/client";
import { FaPoll,FaCheckDouble  } from "react-icons/fa";

const GET_ME_USER = gql`
  query MeUser {
    meUser {
      firstName
      lastName
      isActive
      id
      email
    }
  }
`;

export default function Layout() {
  const { data, loading, error } = useQuery(GET_ME_USER);
  const client = useApolloClient(); 

  useEffect(() => {
    return () => {
      client.cache.evict({ fieldName: 'meUser' });  
    };
  }, [client]);

 const menuItems = useMemo<NavigationEntry[]>(() => [
    {
      groupTitle: "Catalogo",
      items: [
        { text: "Usuarios", redirectTo: "/protected/users", icon: <FaUsers /> },
        { text: "Clientes", redirectTo: "/protected/customers", icon: <FaUsers /> },
        { text: "Estaciones", redirectTo: "/protected/gas-stations", icon: <FaGasPump /> },
        { text: "Productos", redirectTo: "/protected/products", icon: <GiCardExchange /> },
        { text: "Niveles", redirectTo: "/protected/levels", icon: <FaSortAmountUpAlt />},
      ],
    },
    {
      groupTitle: "Procesos",
      items: [
        { text: "Beneficios", redirectTo: "/protected/benefits", icon: <FaCheckDouble /> },
        { text: "Margenes", redirectTo: "/protected/margins", icon: <FaPoll /> },
      ],
    },
    {
      groupTitle: "Reportes",
      items: [
        { text: "Acumulaciones", redirectTo: "/protected/acumulations", icon: <FaAlignJustify /> },
        { text: "Reporte de Transacciones", redirectTo: "/protected/report", icon: <FaCreditCard /> },
        { text: "Niveles de los Clientes", redirectTo: "/protected/customer-levels", icon: <FaUsers /> },
        { text: "Beneficios Generados", redirectTo: "/protected/benefits-generated", icon: <FaCheckDouble /> },
        { text: "Beneficios Asignados", redirectTo: "/protected/benefits-tickets", icon: <FaCheckDouble /> },
      ],
    },
  ], []);
  

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <NavBar 
        firstName={data.meUser.firstName} 
        lastName={data.meUser.lastName} />
      <div className="flex flex-row">
        <SideBar items={menuItems} />
        <div className="w-full">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
