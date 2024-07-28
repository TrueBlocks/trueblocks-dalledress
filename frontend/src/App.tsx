import React, { useEffect } from "react";
import { AppShell, Text } from "@mantine/core";
import { Aside, Header, Navbar, Routes } from "@components";
import { EventsOn, EventsOff } from "@runtime";
import { useLocation } from "wouter";
import classes from "@/App.module.css";

function App() {
  const [showHelp, setShowHelp] = React.useState(true);
  const [, setLocation] = useLocation();

  const toggleHelp = () => {
    setShowHelp(!showHelp);
  };

  useEffect(() => {
    const handleNavigation = (route: string) => {
      console.log(`Navigating to ${route}`);
      setLocation(route);
    };

    EventsOn("navigate", handleNavigation);

    return () => {
      EventsOff("navigate");
    };
  }, [setLocation]);

  return (
    <AppShell
      header={{ height: "3rem" }}
      navbar={{ collapsed: { desktop: false }, width: "10rem", breakpoint: 0 }}
      aside={{ collapsed: { desktop: showHelp }, width: "10rem", breakpoint: 0 }}
      footer={{ height: "2rem" }}
    >
      <AppShell.Header>
        <Header title="ApplicationTitle" />
        <button onClick={toggleHelp}>Toggle Help</button>
      </AppShell.Header>
      <AppShell.Navbar>
        <Navbar />
      </AppShell.Navbar>
      <AppShell.Main className={classes.mainContent}>
        <Routes />
      </AppShell.Main>
      <AppShell.Aside>
        <Aside />
      </AppShell.Aside>
      <AppShell.Footer>
        <Text size={"sm"}>time / date / currently opened file</Text>
      </AppShell.Footer>
    </AppShell>
  );
}

export default App;
