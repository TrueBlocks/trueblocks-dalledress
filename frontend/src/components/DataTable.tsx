import React from "react";
import lClasses from "./DataTable.module.css";
import { flexRender, Table as ReactTable } from "@tanstack/react-table";
import { Table, Title } from "@mantine/core";
import { CustomMeta } from "./";

export function DataTable<T>({ table, loading }: { table: ReactTable<T>; loading: boolean }) {
  if (loading) {
    return <Title order={3}>Loading...</Title>;
  } else {
    return (
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
                const meta = cell.column.columnDef.meta as CustomMeta;
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
    );
  }
}
