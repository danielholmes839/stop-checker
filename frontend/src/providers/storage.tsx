import { TravelLegInput } from "client/types";
import React, { useEffect, useState } from "react";

type Route = TravelLegInput[];

type StorageValue = { [key: string]: Route };

type StorageContextValue = {
  routes: Route[];
  add: (route: Route) => void;
  remove: (route: Route) => void;
  has: (route: Route) => boolean;
};

const getKey = (legs: Route): string => {
  return legs
    .map(({ origin, destination, route }, i) => {
      let last = legs.length === i + 1;
      return `${origin}:${route ? route : "W"}:${last ? destination : ""}`;
    })
    .join("");
};

const StorageContext = React.createContext<StorageContextValue>({
  add: (route) => {},
  has: (route) => false,
  remove: (route) => {},
  routes: [],
});

export const useStorage = () => React.useContext(StorageContext);

export const StorageProvider: React.FC = ({ children }) => {
  // read the routes
  const [routes, setRoutes] = useState<StorageValue>(() => {
    const data = localStorage.getItem("storage-value");
    if (data === null) {
      return {};
    }
    return JSON.parse(data);
  });

  // update the local storage when routes change
  useEffect(() => {
    localStorage.setItem("storage-value", JSON.stringify(routes, undefined, 0));
  }, [routes]);

  const has = (route: Route) => {
    return routes[getKey(route)] !== undefined;
  };

  const add = (route: Route) => {
    if (has(route)) {
      return route;
    }
    let key = getKey(route);
    let copy = { ...routes };
    copy[key] = route;
    setRoutes(copy);
  };

  const remove = (route: Route) => {
    let copy = { ...routes };
    delete copy[getKey(route)];
    setRoutes(copy);
  };

  return (
    <StorageContext.Provider
      value={{
        add: add,
        has: has,
        remove: remove,
        routes: Object.values(routes),
      }}
    >
      {children}
    </StorageContext.Provider>
  );
};
