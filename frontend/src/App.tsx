import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Provider as UrqlProvider } from "urql";
import { client } from "client";
import { SearchPage, StopPage, StopRoutePage } from "pages";
import { Nav } from "components";
import { AutomaticInput, AutomaticOutput } from "pages/travel";
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
            <Route path="/travel/automatic" element={<AutomaticInput />} />
            <Route
              path="/travel/:origin/:destination"
              element={<AutomaticOutput />}
            />
            <Route path="/dashboard" element={<Dashboard />} />
          </Routes>
        </BrowserRouter>
      </StorageProvider>
    </UrqlProvider>
  );
};

export default App;
