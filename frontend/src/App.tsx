import React from "react";
import { AppShell } from "@mantine/core";
import Navbar from "./components/global/Navbar";
import Routes from "./components/global/Routes";

function App() {
  return (
    <AppShell navbar={{ width: "15rem", breakpoint: 0 }}>
      <AppShell.Navbar>
        <Navbar />
      </AppShell.Navbar>
      <AppShell.Main>
        <Routes />
      </AppShell.Main>
    </AppShell>
  );
}

export default App;
