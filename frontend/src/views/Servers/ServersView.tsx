import React, { useState, useEffect, Dispatch, SetStateAction } from "react";
import { View, ViewStatus } from "@components";
import { ServerCard, ServerLog } from "./";
import { servers } from "@gocode/models";
import { GetServer, ToggleServer } from "@gocode/app/App";
import { Stack, Title, SimpleGrid } from "@mantine/core";
import { EventsOn, EventsOff } from "@runtime";
import classes from "@/App.module.css";

// TODO: This should be a type from GoLang
type Progress = {
  name: string;
  message: string;
  color: string;
};

export function ServersView() {
  const [scraper, setScraper] = useState<servers.Server>({} as servers.Server);
  const [fileServer, setFileServer] = useState<servers.Server>({} as servers.Server);
  const [monitor, setMonitor] = useState<servers.Server>({} as servers.Server);
  const [ipfs, setIpfs] = useState<servers.Server>({} as servers.Server);
  const [logMessages, setLogMessages] = useState<Progress[]>([]);

  const updateServer = (server: string, setStateFn: Dispatch<SetStateAction<servers.Server>>) => {
    GetServer(server).then((s) => {
      setStateFn(s);
    });
  };

  useEffect(() => {
    updateServer("scraper", setScraper);
    updateServer("fileserver", setFileServer);
    updateServer("monitor", setMonitor);
    updateServer("ipfs", setIpfs);
  }, []);

  const handleMessage = (p: Progress) => {
    switch (p.name) {
      case "scraper":
        updateServer("scraper", setScraper);
        break;
      case "fileserver":
        updateServer("fileserver", setFileServer);
        break;
      case "monitor":
        updateServer("monitor", setMonitor);
        break;
      case "ipfs":
        updateServer("ipfs", setIpfs);
        break;
      default:
        break;
    }
    setLogMessages((prev) => {
      const newLogs = [...prev, p];
      return newLogs.length > 4 ? newLogs.slice(-4) : newLogs;
    });
  };

  useEffect(() => {
    EventsOn("Server", handleMessage);
    return () => {
      EventsOff("Server");
    };
  }, []);

  const toggleServer = (name: string) => {
    ToggleServer(name);
  };

  return (
    <View>
      <Title order={3}>Servers View</Title>
      <Stack className={classes.mainContent}>
        <SimpleGrid cols={2} spacing="lg" style={{ padding: "lg" }}>
          <ServerCard s={scraper} toggle={toggleServer} />
          <ServerCard s={monitor} toggle={toggleServer} />
          <ServerCard s={ipfs} toggle={toggleServer} />
          <ServerCard s={fileServer} toggle={toggleServer} />
        </SimpleGrid>
        <ServerLog logMessages={logMessages} />
      </Stack>
      <ViewStatus />
    </View>
  );
}
