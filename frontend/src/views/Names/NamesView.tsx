import React, { useState, useEffect } from "react";
import classes from "@/App.module.css";
import { types } from "@gocode/models";
import { GetNames, GetNamesCnt } from "@gocode/app/App";
import { createColumnHelper, flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useHotkeys } from "react-hotkeys-hook";
import { Stack, Table, Title } from "@mantine/core";
import { namesColumns, createRowModel } from "./Names";
import { View, ViewStatus } from "@/components/view";

const columnHelper = createColumnHelper<types.Name>();

export function NamesView() {
  const [items, setItems] = useState<types.Name[]>([]);
  const [nItems, setNItems] = useState<number>(0);
  const [curItem, setCurItem] = useState<number>(0);
  const perPage = 20;

  useHotkeys("left", (event) => {
    setCurItem((cur) => Math.max(cur - 1, 0));
    event.preventDefault();
  });
  useHotkeys("up", (event) => {
    setCurItem((cur) => Math.max(cur - perPage, 0));
    event.preventDefault();
  });
  useHotkeys("home", (event) => {
    setCurItem(0);
    event.preventDefault();
  });

  useHotkeys("right", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem((cur) => Math.min(max, cur + 1));
    event.preventDefault();
  });
  useHotkeys("down", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem((cur) => Math.min(max, cur + perPage));
    event.preventDefault();
  });
  useHotkeys("end", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem(max);
    event.preventDefault();
  });

  useEffect(() => {
    const fetchHistory = async () => {
      const fetchedItems = await GetNames(curItem, perPage);
      setItems(fetchedItems);
      const cnt = await GetNamesCnt();
      setNItems(cnt);
    };
    fetchHistory();
  }, [curItem, perPage]);

  useEffect(() => {
    setCurItem(0);
  }, []);

  const table = useReactTable({
    data: items,
    columns: namesColumns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <View>
      <Stack className={classes.mainContent}>
        <Title order={3}>Names Title {curItem}</Title>
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
