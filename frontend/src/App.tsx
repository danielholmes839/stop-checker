import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Provider as UrqlProvider } from "urql";
import { client } from "client";
import { SearchPage, StopPage, StopRoutePage } from "pages";
import { Nav } from "components";
import { AutomaticInput, AutomaticOutput } from "pages/travel";

const App: React.FC = () => {
  return (
    <UrqlProvider value={client}>
      <BrowserRouter>
        <Nav />
        <Routes>
          <Route path="/" element={<SearchPage />} />
          <Route path="/stop/:id" element={<StopPage />} />
          <Route path="/stop/:stop/route/:route" element={<StopRoutePage />} />
          <Route path="/travel/automatic" element={<AutomaticInput />} />
          <Route
            path="/travel/:origin/:destination"
            element={<AutomaticOutput />}
          />
        </Routes>
      </BrowserRouter>
    </UrqlProvider>
  );
};

export default App;
