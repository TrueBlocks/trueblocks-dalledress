import React, { useState, useEffect } from "react";
import classes from "@/App.module.css";
import { GetHistory, GetHistoryCnt } from "@gocode/app/App";
import { app } from "@gocode/models";
import { flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useHotkeys } from "react-hotkeys-hook";
import { Stack, Table, Title } from "@mantine/core";
import { txColumns } from "./HistoryTable";
import { View, ViewStatus } from "@components";
import { useKeyboardPaging2 } from "@hooks";

export function HistoryView() {
  // const [address, setAddress] = useState<string>("0x5088d623ba0fcf0131e0897a91734a4d83596aa0");
  // const [address, setAddress] = useState<string>("0xee906d7d5f1748258174be4cbc38930302ab7b42");
  // const [address, setAddress] = useState<string>("0x8bae48f227d978d084b009b775222baaf61ed9fe");
  // const [address, setAddress] = useState<string>("0x3fb1cd2cd96c6d5c0b5eb3322d807b34482481d4");
  // const [address, setAddress] = useState<string>("0x5ed8cee6b63b1c6afce3ad7c92f4fd7e1b8fad9f");
  // const [address, setAddress] = useState<string>("0x1db3439a222c519ab44bb1144fc28167b4fa6ee6");
  // const [address, setAddress] = useState<string>("0xd8da6bf26964af9d7eed9e03e53415d37aa96045");
  // const [address, setAddress] = useState<string>("0x9531c059098e3d194ff87febb587ab07b30b1306");
  const [address, setAddress] = useState<string>("0xf503017d7baf7fbc0fff7492b751025c6a78179b");
  const [nItems, setNItems] = useState<number>(0);
  const { items, curItem } = useKeyboardPaging2<app.TransactionEx>(
    (curItem, perPage) => GetHistory(address, curItem, perPage),
    nItems,
    address,
    20
  );

  useEffect(() => {
    const fetch = async () => {
      const cnt = await GetHistoryCnt(address);
      setNItems(cnt);
    };
    fetch();
  }, [address]);

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
