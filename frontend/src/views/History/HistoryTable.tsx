import React from "react";
import { IconCircleCheck } from "@tabler/icons-react";
import { types } from "@gocode/models";
import { createColumnHelper } from "@tanstack/react-table";

export function createRowModel(tx: types.Transaction): types.Transaction {
  var cV = function (a: any, classs: any, asMap: boolean = false): any {
    return a;
  };
  return {
    // articulatedTx?: Function;
    blockHash: tx.blockHash,
    blockNumber: tx.blockNumber,
    from: tx.from,
    gas: tx.gas,
    gasPrice: tx.gasPrice,
    gasUsed: tx.gasUsed,
    hasToken: tx.hasToken,
    hash: tx.hash,
    input: tx.input,
    isError: tx.isError,
    maxFeePerGas: tx.maxFeePerGas,
    maxPriorityFeePerGas: tx.maxPriorityFeePerGas,
    nonce: tx.nonce,
    // receipt?: Receipt,
    timestamp: tx.timestamp,
    to: tx.to,
    traces: tx.traces,
    transactionIndex: tx.transactionIndex,
    type: tx.type,
    value: tx.value,
    // statements?: Statement[],
    convertValues: cV,
  };
}

const txColumnHelper = createColumnHelper<types.Transaction>();

export const txColumns = [
  txColumnHelper.accessor("blockNumber", {
    header: () => "Block No.",
    cell: (info) => info.getValue(),
    size: 100,
  }),
  txColumnHelper.accessor("transactionIndex", {
    header: () => "Txid",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  txColumnHelper.accessor("timestamp", {
    header: () => "Timestamp",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  txColumnHelper.accessor("from", {
    header: () => "From",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  txColumnHelper.accessor("to", {
    header: () => "To",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  txColumnHelper.accessor("value", {
    header: () => "Value",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  txColumnHelper.accessor("hasToken", {
    header: () => "hasToken",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="white" fill="green" /> : ""),
    size: 100,
  }),
  txColumnHelper.accessor("isError", {
    header: () => "isError",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="green" fill="red" /> : ""),
    size: 100,
  }),
];
