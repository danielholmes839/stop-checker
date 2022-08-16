import React, { useEffect, useState } from "react";
import { NavLink, useLocation } from "react-router-dom";
import { Container } from "./helper";

const links = [
  { to: "/", text: "Search" },
  { to: "/travel", text: "Travel Planner" },
  { to: "/dashboard", text: "Dashboard" },
];

const linkClassName = "hover:text-primary-800 text-primary-600";
const linkClassNameActive = "hover:text-primary-800 text-primary-600 underline";
const linkClassNameFunc: any = ({ isActive }: { isActive: boolean }) =>
  isActive ? linkClassNameActive : linkClassName;

const NavStandard: React.FC = () => {
  return (
    <div className="hidden w-full md:block md:w-auto">
      <ul className="flex flex-col p-4 mt-4 md:flex-row">
        {links.map(({ to, text }) => (
          <li className="ml-5">
            <NavLink to={to} className={linkClassNameFunc}>
              {text}
            </NavLink>
          </li>
        ))}
      </ul>
    </div>
  );
};

const NavDropDown: React.FC = () => {
  return (
    <div className="md:hidden p-3">
      <ul>
        {links.map(({ to, text }) => (
          <li className="mt-1">
            <NavLink to={to} className={linkClassNameFunc}>
              {`${text}`}
            </NavLink>
          </li>
        ))}
      </ul>
    </div>
  );
};

export const Nav: React.FC = () => {
  const [dropDown, setDropDown] = useState(false);
  const location = useLocation();
  useEffect(() => {
    setDropDown(false);
  }, [location]);
  return (
    <nav className="bg-gray-50 py-3">
      <Container>
        <div className="flex flex-wrap justify-between items-center mx-auto">
          <h1 className="text-3xl lg:text-4xl text-primary-600 font-semibold">
            stop-checker.com
          </h1>
          <button
            onClick={() => setDropDown(!dropDown)}
            className="inline-flex items-center ml-3 text-sm text-primary-500 rounded md:hidden hover:bg-primary-100 focus:outline-none"
          >
            <svg
              className="w-6 h-6"
              aria-hidden="true"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill-rule="evenodd"
                d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                clip-rule="evenodd"
              ></path>
            </svg>
          </button>
          <NavStandard />
        </div>
        {dropDown && <NavDropDown />}
      </Container>
    </nav>
  );
};
