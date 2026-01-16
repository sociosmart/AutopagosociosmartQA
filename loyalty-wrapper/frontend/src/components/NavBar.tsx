import { useDispatch } from "react-redux";
import { logOut } from "../store/authSlice";
import { useNavigate } from "react-router-dom";
import Drawer from "./Drawer"; 
import { FaGasPump, FaUsers } from "react-icons/fa";
import { GiCardExchange } from "react-icons/gi";
import  {NavigationEntry } from "./SideBar";
import { useMemo, useEffect } from "react";

interface NavBarProps {
  firstName: string;
  lastName: string;
}

export default function NavBar({ firstName, lastName }: NavBarProps) {
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogout = () => {
    dispatch(logOut());
    navigate("/login");
  };

  const menuItems = useMemo<NavigationEntry[]>(() => [
    {
      groupTitle: "Catalogo",
      items: [
        { text: "Usuarios", redirectTo: "/protected/users", icon: <FaUsers /> },
        { text: "Clientes", redirectTo: "/protected/customers", icon: <FaUsers /> },
        { text: "Estaciones", redirectTo: "/protected/gas-stations", icon: <FaGasPump /> },
        { text: "Productos", redirectTo: "/protected/products", icon: <GiCardExchange /> },
        { text: "Niveles", redirectTo: "/protected/levels", icon: <GiCardExchange /> },
      ],
    },
    {
      groupTitle: "Procesos",
      items: [
        { text: "Beneficios", redirectTo: "/protected/benefits", icon: <GiCardExchange /> },
        { text: "Margenes", redirectTo: "/protected/margins", icon: <GiCardExchange /> },
      ],
    },
    {
      groupTitle: "Reportes",
      items: [
        { text: "Acumulaciones", redirectTo: "/protected/acumulations", icon: <GiCardExchange /> },
        { text: "Reporte de Transacciones", redirectTo: "/protected/report", icon: <GiCardExchange /> },
        { text: "Niveles de los Clientes", redirectTo: "/protected/customer-levels", icon: <FaUsers /> },
        { text: "Beneficios Generados", redirectTo: "/protected/benefits-generated", icon: <GiCardExchange /> },
        { text: "Beneficios Asignados", redirectTo: "/protected/benefits-tickets", icon: <GiCardExchange /> },
      ],
    },
  ], []);
  

  return (
    <div className="navbar bg-base-100 sticky top-0 left-0 right-0 z-50">
      <div className="flex-none">
        < Drawer  items={menuItems}  />
      </div>
      <div className="flex-1">
        <a className="btn btn-ghost text-xl p-0">SmartGas</a>
      </div>

      <div className="">
        <ul className="menu menu-horizontal px-1">
          <li>
            <details>
              <summary>{firstName} {lastName}</summary>
              <ul className="bg-base-100 rounded-t-none p-2">
                <li>
                  <button onClick={handleLogout}>Log Out</button>
                </li>
              </ul>
            </details>
          </li>
        </ul>
      </div>
    </div>
  );
}




