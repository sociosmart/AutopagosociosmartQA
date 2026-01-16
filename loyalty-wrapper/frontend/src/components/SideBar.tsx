import React from "react";
import { Link, useLocation } from "react-router";

export interface SideBarNavigationItem {
  text: string;
  icon?: React.ReactNode;
  redirectTo: string;
}

export interface SideBarNavigationGroup {
  groupTitle: string;
  items: SideBarNavigationItem[];
}

export type NavigationEntry = SideBarNavigationItem | SideBarNavigationGroup;

interface SideBarProps {
  items: NavigationEntry[];
}

export default function SideBar({ items }: SideBarProps) {
  const location = useLocation();

  return (
    <div className="h-screen flex flex-row w-56 max-md:hidden">
      <ul className="menu  rounded-box w-56 p-0">
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
                >
                  {entry.icon}
                  {entry.text}
                </Link>
              </li>
            );
          }
        })}
      </ul>
      <div className="flex">
        <div className="divider divider-horizontal"></div>
      </div>
    </div>
  );
}
