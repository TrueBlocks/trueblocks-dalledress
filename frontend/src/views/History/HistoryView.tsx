import React, { useState, useEffect } from "react";
import classes from "@/App.module.css";
import { GetHistory, GetHistoryCnt } from "@gocode/app/App";
import { types } from "@gocode/models";
import { flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useHotkeys } from "react-hotkeys-hook";
import { Stack, Table, Title } from "@mantine/core";
import { txColumns } from "./HistoryTable";
import { View, ViewStatus } from "@components";
import { useKeyboardPaging } from "@hooks";

export function HistoryView() {
  const [address, setAddress] = useState<string>("0x9531c059098e3d194ff87febb587ab07b30b1306");
  // const [address, setAddress] = useState<string>("0xf503017d7baf7fbc0fff7492b751025c6a78179b");
  const { items, nItems, curItem } = useKeyboardPaging<types.Transaction>(
    (curItem, perPage) => GetHistory(address, curItem, perPage),
    () => GetHistoryCnt(address),
    address,
    20
  );

  const table = useReactTable({
    data: items,
    columns: txColumns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <View>
      <Stack className={classes.mainContent}>
        <Title order={3}>
          History {curItem} of {nItems}
        </Title>
        <Table>
          <Table.Thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <Table.Tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <Table.Th key={header.id}>
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                  </Table.Th>
                ))}
              </Table.Tr>
            ))}
          </Table.Thead>
          <Table.Tbody>
            {table.getRowModel().rows.map((row) => (
              <Table.Tr key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <Table.Td style={{ verticalAlign: "top" }} key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </Table.Td>
                ))}
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      </Stack>
      <ViewStatus />
    </View>
  );
}
