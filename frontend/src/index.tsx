import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import "./css/custom.css";
import "./css/tailwind.css";
import "react-datetime-picker/dist/DateTimePicker.css";
import "react-calendar/dist/Calendar.css";
import "react-clock/dist/Clock.css";

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root")
);
