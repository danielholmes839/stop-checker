import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Provider as UrqlProvider } from "urql";
import { client } from "client";
import {
  Dashboard,
  DashboardAddFavourite,
  DashboardEditFavourite,
  SearchPage,
  StopPage,
  StopRoutePage,
  TravelDestinationInput,
  TravelOriginInput,
  TravelSchedule,
} from "components";
import { Nav } from "components/util";
import { StorageProvider } from "core";

const App: React.FC = () => {
  return (
    <UrqlProvider value={client}>
      <BrowserRouter>
        <StorageProvider>
          <Nav />
          <Routes>
            <Route path="/" element={<SearchPage />} />
            <Route path="/stop/:id" element={<StopPage />} />
            <Route
              path="/stop/:stop/route/:route"
              element={<StopRoutePage />}
            />
            <Route path="/travel" element={<TravelDestinationInput />} />
            <Route
              path="/travel/p/:destinationId"
              element={<TravelOriginInput />}
            />
            <Route
              path="/travel/p/:destinationId/:originId"
              element={<TravelSchedule />}
            />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route
              path="/dashboard/favourite/add"
              element={<DashboardAddFavourite />}
            />
            <Route
              path="/dashboard/favourite/:id"
              element={<DashboardEditFavourite />}
            />
          </Routes>
        </StorageProvider>
      </BrowserRouter>
    </UrqlProvider>
  );
};

export default App;
