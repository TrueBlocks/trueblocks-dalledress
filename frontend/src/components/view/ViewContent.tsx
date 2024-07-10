import React from "react";
import { Text } from "@mantine/core";

export function ViewContent({ children }: { children: React.ReactNode }) {
  return (
    <Text
      style={{
        flex: 1,
      }}
    >
      {children}
    </Text>
  );
}
