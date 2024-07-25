import React from "react";
import { Card, Text, Title } from "@mantine/core";

// TODO: This should be a type from GoLang
type Progress = {
  name: string;
  message: string;
  color: string;
};

export const ServerLog = ({ logMessages }: { logMessages: Progress[] }) => {
  return (
    <>
      <Title order={4} >Logs:</Title>
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
