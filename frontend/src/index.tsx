import { client } from "client";
import React from "react";
import ReactDOM from "react-dom";
import { Provider } from "urql";
import App from "./App";
import "./css/custom.css";
import "./css/tailwind.css";

ReactDOM.render(
  <React.StrictMode>
    <Provider value={client}>
      <App />
    </Provider>
  </React.StrictMode>,
  document.getElementById("root")
);
