import React from "react";
import { NavLink } from "react-router-dom";
import { Container } from "./page";

export const Nav: React.FC = () => {
  const linkClassName = "text-indigo-600";
  const linkClassNameActive = "text-indigo-600 underline";
  const linkClassNameFunc: any = ({ isActive }: { isActive: boolean }) =>
    isActive ? linkClassNameActive : linkClassName;

  return (
    <nav className="bg-gray-50 py-3 mb-3">
      <Container>
        <div className="flex flex-wrap justify-between items-center mx-auto">
          <h1 className="text-4xl text-indigo-600 font-semibold">
            stop-checker.com
          </h1>
          {/* <Link className="text-blue-500" to={"/planner"}>
              Travel Planner
            </Link> */}
          <div>
            <ul className="flex flex-col p-4 mt-4 md:flex-row">
              <li className="ml-5">
                <NavLink to="/" className={linkClassNameFunc}>
                  Search
                </NavLink>
              </li>
              <li className="ml-5">
                <NavLink to="/planner" className={linkClassNameFunc}>
                  Travel Planner
                </NavLink>
              </li>
            </ul>
          </div>
        </div>
      </Container>
    </nav>
  );
};
