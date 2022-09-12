import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Provider as UrqlProvider } from "urql";
import { client } from "client";
import { SearchPage, StopPage, StopRoutePage } from "pages";
import { Nav } from "components";
import {
  AutomaticInput,
  AutomaticOutput,
  FixedRoute,
  ManualLegInput,
  ManualOriginInput,
} from "pages/travel";
import { StorageProvider } from "providers/storage";
import { Dashboard } from "pages/dashboard";

const App: React.FC = () => {
  return (
    <UrqlProvider value={client}>
      <StorageProvider>
        <BrowserRouter>
          <Nav />
          <Routes>
            <Route path="/" element={<SearchPage />} />
            <Route path="/stop/:id" element={<StopPage />} />
            <Route
              path="/stop/:stop/route/:route"
              element={<StopRoutePage />}
            />
            <Route path="/travel" element={<AutomaticInput />} />
            <Route path="/travel/create" element={<ManualOriginInput />} />
            <Route path="/travel/m/:origin" element={<ManualLegInput />} />
            <Route
              path="/travel/p/:origin/:destination"
              element={<AutomaticOutput />}
            />
            <Route path="/travel/r/:encoded" element={<FixedRoute />} />
            <Route path="/dashboard" element={<Dashboard />} />
          </Routes>
        </BrowserRouter>
      </StorageProvider>
    </UrqlProvider>
  );
};

export default App;
