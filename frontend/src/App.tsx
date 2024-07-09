import React from "react";
import { AppShell } from "@mantine/core";
import Header from "./components/global/Header";
import Navbar from "./components/global/Navbar";
import Routes from "./components/global/Routes";

function App() {
  return (
    <AppShell
      header={{ height: "3rem" }}
      navbar={{ width: "15rem", breakpoint: 0 }}
    >
      <AppShell.Header>
        <Header />
      </AppShell.Header>
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
