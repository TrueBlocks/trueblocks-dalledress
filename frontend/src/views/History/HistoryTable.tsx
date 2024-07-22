import React from "react";
import { IconCircleCheck } from "@tabler/icons-react";
import { app } from "@gocode/models";
import { createColumnHelper, ColumnDef } from "@tanstack/react-table";
import { CustomMeta } from "../CustomMeta";
import lClasses from "../Columns.module.css";

type CustomColumnDef<TData, TValue> = ColumnDef<TData, TValue> & {
  meta?: CustomMeta;
};

const txColumnHelper = createColumnHelper<app.TransactionEx>();

export const txColumns: CustomColumnDef<app.TransactionEx, any>[] = [
  txColumnHelper.accessor((row) => `${row.blockNumber}.${row.transactionIndex}`, {
    id: "blockTx",
    header: () => "Id",
    cell: (info) => info.getValue(),
    meta: { className: lClasses.small },
  }),
  txColumnHelper.accessor("date", {
    header: () => "Date",
    cell: (info) => info.renderValue(),
    meta: { className: lClasses.medium },
  }),
  txColumnHelper.accessor("fromName", {
    header: () => "From",
    cell: (info) => info.renderValue(),
    meta: { className: lClasses.wide },
  }),
  txColumnHelper.accessor("toName", {
    header: () => "To",
    cell: (info) => info.renderValue(),
    meta: { className: "wide" },
  }),
  txColumnHelper.accessor("logCount", {
    header: () => "nEvents",
    cell: (info) => info.renderValue(),
    meta: { className: lClasses.medium },
  }),
  txColumnHelper.accessor("ether", {
    header: () => "Ether",
    cell: (info) => info.renderValue(),
    meta: { className: lClasses.medium },
  }),
  txColumnHelper.accessor("hasToken", {
    header: () => "hasToken",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="white" fill="green" /> : ""),
    meta: { className: lClasses.centered },
  }),
  txColumnHelper.accessor("isError", {
    header: () => "isError",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="green" fill="red" /> : ""),
    meta: { className: lClasses.centered },
  }),
];
