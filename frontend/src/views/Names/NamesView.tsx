import React, { useState, useEffect } from "react";
import classes from "@/App.module.css";
import { GetNames, GetNamesCnt } from "@gocode/app/App";
import { types } from "@gocode/models";
import { flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useHotkeys } from "react-hotkeys-hook";
import { Stack, Table, Title } from "@mantine/core";
import { nameColumns } from "./NameTable";
import { View, ViewStatus } from "@components";
import { useKeyboardPaging } from "@hooks";

export function NamesView() {
  const { items, nItems, curItem } = useKeyboardPaging<types.Name>(GetNames, GetNamesCnt, undefined, 20);
  const table = useReactTable({
    data: items,
    columns: nameColumns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <View>
      <Stack className={classes.mainContent}>
        <Title order={3}>
          Names {curItem} of {nItems}
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
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}
