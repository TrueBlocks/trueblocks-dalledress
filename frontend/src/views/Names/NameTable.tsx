import React from "react";
import { IconCircleCheck } from "@tabler/icons-react";
import { types } from "@gocode/models";
import { createColumnHelper, ColumnDef } from "@tanstack/react-table";
import { CustomMeta } from "../CustomMeta";

type CustomColumnDef<TData, TValue> = ColumnDef<TData, TValue> & {
  meta?: CustomMeta;
};

const nameColumnHelper = createColumnHelper<types.Name>();

export const nameColumns: CustomColumnDef<types.Name, any>[] = [
  nameColumnHelper.accessor("address", {
    header: () => "Address",
    cell: (info) => info.renderValue(),
    size: 100,
    meta: { className: "wide" },
  }),
  nameColumnHelper.accessor("tags", {
    header: () => "Tags",
    cell: (info) => info.renderValue(),
    size: 100,
    meta: { className: "wide" },
  }),
  nameColumnHelper.accessor("name", {
    header: () => "Name",
    cell: (info) => info.renderValue(),
    size: 100,
    meta: { className: "wide" },
  }),
  nameColumnHelper.accessor("symbol", {
    header: () => "Symbol",
    cell: (info) => info.renderValue(),
    size: 100,
    meta: { className: "small" },
  }),
  nameColumnHelper.accessor("decimals", {
    header: () => "Decimals",
    cell: (info) => info.renderValue(),
    size: 100,
    meta: { className: "small" },
  }),
  nameColumnHelper.accessor("isContract", {
    header: () => "isContract",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="white" fill="green" /> : ""),
    size: 100,
    meta: { className: "small-centered" },
  }),
  nameColumnHelper.accessor("isErc20", {
    header: () => "isErc20",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="white" fill="green" /> : ""),
    size: 100,
    meta: { className: "small-centered" },
  }),
  nameColumnHelper.accessor("isErc721", {
    header: () => "isErc721",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="white" fill="green" /> : ""),
    size: 100,
    meta: { className: "small-centered" },
  }),
];
