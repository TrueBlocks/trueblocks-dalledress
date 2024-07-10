import React from "react";
import { Text } from "@mantine/core";

export function ViewStatus({ children }: { children: React.ReactNode }) {
  return <Text size="xs">{children}</Text>;
}
