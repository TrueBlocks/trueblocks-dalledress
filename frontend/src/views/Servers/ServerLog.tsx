import React from "react";
import { Card, Text, Title } from "@mantine/core";
import { messages } from "@gocode/models";

export const ServerLog = ({ logMessages }: { logMessages: messages.ServerMsg[] }) => {
  return (
    <>
      <Title order={4}>Logs:</Title>
      <Card shadow="xl" style={{ maxHeight: "16rem", overflowY: "auto" }}>
        {logMessages.map((log, index) => (
          <Text c={log.color} key={index}>
            {log.message}
          </Text>
        ))}
      </Card>
    </>
  );
};
