import React, { useState, useEffect } from "react";
import { View, ViewStatus } from "@components";
import { servers } from "@gocode/models";
import { GetServer /* , ToggleServer */ } from "@gocode/app/App";
import { Stack, Title, SimpleGrid } from "@mantine/core";
import { EventsOn, EventsOff } from "@runtime";
import classes from "@/App.module.css";

export function ServersView() {
  const [scraper, setScraper] = useState<servers.Server>({} as servers.Server);

  const updateServer = (server: string, update: React.Dispatch<React.SetStateAction<servers.Server>>) => {
    GetServer(server).then((s) => {
      update(s);
    });
  };

  useEffect(() => {
    updateServer("scraper", setScraper);
  }, []);

  return (
    <View>
      <Title order={3}>Servers View</Title>
      <Stack className={classes.mainContent}>
        <SimpleGrid cols={2} spacing="lg" style={{ padding: "lg" }}>
          <div>{JSON.stringify(scraper)}</div>
        </SimpleGrid>
      </Stack>
      <ViewStatus />
    </View>
  );
}
