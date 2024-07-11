import React from "react";
import { Table } from "@mantine/core";

export function VAlignedTd({ children, key }: { children: React.ReactNode; key: string }) {
  return <Table.Td style={{ verticalAlign: "top" }}>{children}</Table.Td>;
}
