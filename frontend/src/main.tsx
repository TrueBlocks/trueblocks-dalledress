import React from "react";
import { MantineProvider } from "@mantine/core";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./style.css";
import "@mantine/core/styles.css";

const container = document.getElementById("root");
const root = createRoot(container!);
root.render(
  <React.StrictMode>
    <MantineProvider defaultColorScheme="dark">
      <App />
    </MantineProvider>
  </React.StrictMode>
);
