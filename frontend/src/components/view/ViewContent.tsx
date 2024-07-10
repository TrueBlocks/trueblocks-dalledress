import React, { ReactNode } from "react";
import { Text } from "@mantine/core";

export function ViewContent({ children }: { children: ReactNode }) {
  return (
    <Text
      style={{
        flexGrow: 1,
      }}
    >
      {children}
    </Text>
  );
}
