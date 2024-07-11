import React, { useState, useEffect } from "react";
import classes from "@/App.module.css";
import { types } from "@gocode/models";
import { GetNames, GetNamesCnt } from "@gocode/app/App";
import { createColumnHelper, flexRender, getCoreRowModel, useReactTable } from "@tanstack/react-table";
import { useHotkeys } from "react-hotkeys-hook";
import { Stack, Table, Title } from "@mantine/core";
import { type Name, defaultData, columns, createRowModel } from "./Names";
import { View, ViewStatus } from "@/components/view";
import { VAlignedTd } from "./VAlignedTd";

const columnHelper = createColumnHelper<Name>();

export function NamesView() {
  const [data, _setData] = React.useState(() => [...defaultData]);
  const [names, setName] = useState<types.Name[]>();
  const [nNames, setNamesCount] = useState<number>(0);
  const [curName, setCurName] = useState<number>(-1);

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
    GetNames(curName, 20).then((names: types.Name[]) => {
      _setData(names?.map((name) => createRowModel(name)) ?? []);
      setName(names);
      GetNamesCnt().then((cnt) => {
        setNamesCount(cnt);
      });
    });
  }, [curName]);

  useEffect(() => {
    setCurName(0);
  }, []);

  const table = useReactTable({
    data,
    columns,
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
                  <VAlignedTd key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</VAlignedTd>
                ))}
              </Table.Tr>
            ))}
          </Table.Tbody>
          {/* <Table.Tfoot>
        {table.getFooterGroups().map((footerGroup) => (
          <Table.Tr key={footerGroup.id}>
            {footerGroup.headers.map((header) => (
              <Table.Th key={header.id}>
                {header.isPlaceholder ? null : flexRender(header.column.columnDef.footer, header.getContext())}
              </Table.Th>
            ))}
          </Table.Tr>
        ))}
      </Table.Tfoot> */}
        </Table>
      </Stack>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}
