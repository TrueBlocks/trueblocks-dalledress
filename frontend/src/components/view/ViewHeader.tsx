import React from "react";
import { Title } from "@mantine/core";

export function ViewHeader({ children }: { children: React.ReactNode }) {
  return (
    <div>
      <Title order={3}>{children}</Title>
    </div>
  );
}
