import React from "react";
import { AppShell } from "@mantine/core";
import Aside from "./components/global/Aside";
import Header from "./components/global/Header";
import Navbar from "./components/global/Navbar";
import Routes from "./components/global/Routes";

function App() {
  return (
    <AppShell
      header={{ height: "3rem" }}
      navbar={{ collapsed: { desktop: false }, width: "10rem", breakpoint: 0 }}
      aside={{ collapsed: { desktop: false }, width: "10rem", breakpoint: 0 }}
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
      <AppShell.Aside>
        <Aside />
      </AppShell.Aside>
    </AppShell>
  );
}

export default App;
