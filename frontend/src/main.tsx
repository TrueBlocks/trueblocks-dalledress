import React from "react";
import { createTheme, MantineProvider } from "@mantine/core";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./style.css";

const theme = createTheme({
  // scale: 0.8
  // white: "red",
  // black: "blue"
});

const container = document.getElementById("root");
const root = createRoot(container!);
root.render(
  <React.StrictMode>
    <MantineProvider theme={theme}>
      <App />
    </MantineProvider>
  </React.StrictMode>
);
