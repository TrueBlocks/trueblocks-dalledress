import React, { useState, useEffect } from "react";
import classes from "@/App.module.css";
import lClasses from "../Columns.module.css";
import { GetNames, GetNamesCnt } from "@gocode/app/App";
import { types } from "@gocode/models";
import { flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { Stack, Table, Title } from "@mantine/core";
import { nameColumns } from "./NameTable";
import { CustomMeta } from "../CustomMeta";
import { View, ViewStatus } from "@components";
import { useKeyboardPaging } from "@hooks";

export function NamesView() {
  const [count, setCount] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [loaded, setLoaded] = useState<boolean>(false);
  const [items, setItems] = useState<types.Name[]>([]);
  const { curItem, perPage } = useKeyboardPaging<types.Name>(items, count);

  useEffect(() => {
    if (loaded && !loading) {
      const fetch = async (currentItem: number, itemsPerPage: number) => {
        const newItems = await GetNames(currentItem, itemsPerPage);
        setItems(newItems);
      };
      fetch(curItem, perPage);
    }
  }, [count, curItem, perPage]);

  useEffect(() => {
    setLoading(true);
    try {
      const fetch = async () => {
        const cnt = await GetNamesCnt();
        setCount(cnt);
      };
      fetch();
      setLoaded(true);
    } finally {
      setLoading(false);
    }
  }, []);

  const table = useReactTable({
    data: items,
    columns: nameColumns,
    getCoreRowModel: getCoreRowModel(),
  });

  if (loading) {
    return (
      <View>
        <Stack className={classes.mainContent}>
          <Title order={3}>Loading...</Title>
        </Stack>
      </View>
    );
  }

  return (
    <View>
      <Stack className={classes.mainContent}>
        <Title order={3}>
          Names: showing record {curItem + 1}-{curItem + 1 + perPage - 1} of {count}
        </Title>
        <Table>
          <Table.Thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <Table.Tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <Table.Th key={header.id} className={lClasses.centered}>
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                  </Table.Th>
                ))}
              </Table.Tr>
            ))}
          </Table.Thead>
          <Table.Tbody>
            {table.getRowModel().rows.map((row) => (
              <Table.Tr key={row.id}>
                {row.getVisibleCells().map((cell) => {
                  var meta = cell.column.columnDef.meta as CustomMeta;
                  return (
                    <Table.Td key={cell.id} className={lClasses[meta?.className || ""]}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </Table.Td>
                  );
                })}
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      </Stack>
      <ViewStatus />
    </View>
  );
}
