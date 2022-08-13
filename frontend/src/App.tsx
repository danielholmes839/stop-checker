import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Provider as UrqlProvider } from "urql";
import { client } from "client";
import { PlannerPage, SearchPage, StopPage } from "pages";
import { Nav } from "components";

const App: React.FC = () => {
  return (
    <UrqlProvider value={client}>
      <BrowserRouter>
        <Nav />
        <Routes>
          <Route path="/" element={<SearchPage />} />
          <Route path="/stop/:id" element={<StopPage />} />
          <Route path="/planner" element={<PlannerPage />} />
        </Routes>
      </BrowserRouter>
    </UrqlProvider>
  );
};

export default App;