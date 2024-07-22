import React, { useState, useEffect } from "react";
import classes from "@/App.module.css";
import lClasses from "../Columns.module.css";
import { GetHistoryPage, GetHistoryCnt } from "@gocode/app/App";
import { app } from "@gocode/models";
import { flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { Stack, Table, Title } from "@mantine/core";
import { txColumns } from "./HistoryTable";
import { CustomMeta } from "../CustomMeta";
import { EditableSelect, View, ViewStatus } from "@components";
import { useKeyboardPaging2 } from "@hooks";

export function HistoryView() {
  const [address, setAddress] = useState<string>("trueblocks.eth");
  const [count, setCount] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [loaded, setLoaded] = useState<boolean>(false);
  const [items, setItems] = useState<app.TransactionEx[]>([]);
  const { curItem, perPage } = useKeyboardPaging2<app.TransactionEx>(items, count, [address]);

  useEffect(() => {
    if (loaded && !loading) {
      const fetch = async (addr: string, currentItem: number, itemsPerPage: number) => {
        const newItems = await GetHistoryPage(addr, currentItem, itemsPerPage);
        setItems(newItems);
      };
      fetch(address, curItem, perPage);
    }
  }, [count, curItem, perPage]);

  useEffect(() => {
    setLoading(true);
    try {
      const fetch = async (addr: string) => {
        const cnt = await GetHistoryCnt(addr);
        setCount(cnt);
      };
      fetch(address);
      setLoaded(true);
    } finally {
      setLoading(false);
    }
  }, [address]);

  const table = useReactTable({
    data: items,
    columns: txColumns,
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
          History of {address}: showing record {curItem + 1}-{curItem + 1 + perPage - 1} of {count}
        </Title>
        <EditableSelect
          value={address}
          onChange={(value) => setAddress(value)}
          label="Select or enter an address or ENS name"
          placeholder="Enter or select an address"
        />
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
