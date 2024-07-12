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
  const [names, _setData] = React.useState<types.Name[]>(() => []);
  const [nNames, setNamesCount] = useState<number>(0);
  const [curName, setCurName] = useState<number>(-1);
  const [perPage, setPerPage] = useState<number>(20);

  useHotkeys("left", (event) => {
    setCurName(curName - 1 < 0 ? 0 : curName - 1);
    event.preventDefault();
  });
  useHotkeys("up", (event) => {
    setCurName(curName - 20 < 0 ? 0 : curName - 20);
    event.preventDefault();
  });
  useHotkeys("right", (event) => {
    setCurName(curName + 1 >= nNames ? nNames : curName + 1);
    event.preventDefault();
  });
  useHotkeys("down", (event) => {
    setCurName(curName + 20 >= nNames ? nNames - 20 : curName + 20);
    event.preventDefault();
  });
  useHotkeys("home", (event) => {
    setCurName(0);
    event.preventDefault();
  });
  useHotkeys("end", (event) => {
    setCurName(nNames <= 20 ? 20 : nNames - 20);
    event.preventDefault();
  });

  useEffect(() => {
    GetNames(curName, perPage).then((names: types.Name[]) => {
      _setData(names);
      GetNamesCnt().then((cnt) => {
        setNamesCount(cnt);
      });
    });
  }, [curName]);

  useEffect(() => {
    setCurName(0);
  }, []);

  const table = useReactTable({
    data: names,
    columns: namesColumns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <View>
      <Stack className={classes.mainContent}>
        <Title order={3}>Names Title {curName}</Title>
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
