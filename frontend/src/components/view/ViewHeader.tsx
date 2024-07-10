import React, { ReactNode } from "react";
import { Title } from "@mantine/core";

export function ViewHeader({ children }: { children: ReactNode }) {
  return <Title order={3}>{children}</Title>;
}
