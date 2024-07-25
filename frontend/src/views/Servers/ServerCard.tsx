import React from "react";
import { Card, Text, Group, Badge, Title, Button } from "@mantine/core";
import { servers } from "@gocode/models";

export const ServerCard = ({ s, toggle }: { s: servers.Server; toggle: (name: string) => void }) => {
  const { name, sleep, started, runs, state } = s;
  const status = state === servers.State.RUNNING ? "green" : "red";
  var stateStr = "";
  switch (state) {
    case servers.State.STOPPED:
      stateStr = "Stopped";
      break;
    case servers.State.PAUSED:
      stateStr = "Paused";
      break;
    case servers.State.RUNNING:
      stateStr = "Running";
      break;
  }

  const handleToggle = () => {
    toggle(name);
  };

  return (
    <Card shadow="xl" style={{ margin: "1rem" }}>
      <Group style={{ justifyContent: "space-between", marginBottom: 5 }}>
        <Title order={4}>{name}</Title>
        <Badge color={status}>
          <button onClick={handleToggle}>{stateStr}</button>
        </Badge>
      </Group>
      <Text size="sm" style={{ lineHeight: 1.5 }}>
        Sleep Duration: {sleep}
      </Text>
      <Text size="sm" style={{ lineHeight: 1.5 }}>
        Started At: {new Date(started).toLocaleString()}
      </Text>
      <Text size="sm" style={{ lineHeight: 1.5 }}>
        Runs: {runs}
      </Text>
      <pre>{JSON.stringify(s, null, 2)}</pre>
    </Card>
  );
};
