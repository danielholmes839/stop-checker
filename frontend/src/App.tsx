import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Provider as UrqlProvider } from "urql";
import { client } from "client";
import { SearchPage, StopPage, StopRoutePage, Travel } from "components";
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
            <Route path="/travel" element={<Travel />} />
            {/* <Route path="/travel/create" element={<ManualOriginInput />} />
          <Route path="/travel/m/:origin" element={<ManualLegInput />} />
          <Route
            path="/travel/p/:origin/:destination"
            element={<AutomaticOutput />}
          />
          <Route path="/travel/r/:encoded" element={<FixedRoute />} />
          <Route path="/dashboard" element={<Dashboard />} /> */}
          </Routes>
        </StorageProvider>
      </BrowserRouter>
    </UrqlProvider>
  );
};

export default App;
