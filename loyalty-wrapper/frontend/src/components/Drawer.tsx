import React from "react";
import { Link, useLocation } from "react-router-dom";
import { FaAlignJustify } from "react-icons/fa";
import { SideBarNavigationItem, NavigationEntry } from "./SideBar"; // Ajusta la ruta según corresponda

interface DrawerProps {
  items: NavigationEntry[];
}

export default function Drawer({ items }: DrawerProps) {
  const location = useLocation();

  const handleLinkClick = () => {
    const drawerToggle = document.getElementById("my-drawer") as HTMLInputElement;
    if (drawerToggle) {
      drawerToggle.checked = false;
    }
  };

  return (
    <div className="drawer">
      <input id="my-drawer" type="checkbox" className="drawer-toggle" />
      <div className="drawer-content mr-3 ml-3 md:hidden">
        <label htmlFor="my-drawer" className="drawer-button">
          <FaAlignJustify />
        </label>
      </div>
      <div className="drawer-side">
        <label htmlFor="my-drawer" aria-label="close sidebar" className="drawer-overlay"></label>
        <ul className="menu bg-base-200 text-base-content min-h-full w-full md:w-80 p-4">
          {items.map((entry, index) => {
            if ("groupTitle" in entry) {
              return (
                <li key={index}>
                  <h2 className="menu-title">{entry.groupTitle}</h2>
                  <ul>
                    {entry.items.map((item: SideBarNavigationItem, idx) => (
                      <li key={`${index}-${idx}`}>
                        <Link
                          to={item.redirectTo}
                          className={item.redirectTo === location.pathname ? "menu-active" : ""}
                          onClick={handleLinkClick}
                        >
                          {item.icon}
                          {item.text}
                        </Link>
                      </li>
                    ))}
                  </ul>
                </li>
              );
            } else {
              return (
                <li key={index}>
                  <Link
                    to={entry.redirectTo}
                    className={entry.redirectTo === location.pathname ? "menu-active" : ""}
                    onClick={handleLinkClick}
                  >
                    {entry.icon}
                    {entry.text}
                  </Link>
                </li>
              );
            }
          })}
        </ul>
      </div>
    </div>
  );
}
