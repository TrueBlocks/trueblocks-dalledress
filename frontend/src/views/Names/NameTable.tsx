import React from "react";
import { IconCircleCheck } from "@tabler/icons-react";
import { types } from "@gocode/models";
import { createColumnHelper } from "@tanstack/react-table";

// export function createRowModel(name: types.Name): types.Name {
//   var cV = function (a: any, classs: any, asMap: boolean = false): any {
//     return a;
//   };
//   return {
//     address: name.address,
//     decimals: name.decimals,
//     isContract: name.isContract,
//     isErc20: name.isErc20,
//     isErc721: name.isErc721,
//     name: name.name,
//     source: name.source,
//     symbol: name.symbol,
//     tags: name.tags,
//     convertValues: cV,
//   };
// }

const nameColumnHelper = createColumnHelper<types.Name>();

export const nameColumns = [
  nameColumnHelper.accessor("address", {
    header: () => "Address",
    cell: (info) => info.getValue(),
    size: 100,
  }),
  nameColumnHelper.accessor("tags", {
    header: () => "Tags",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  nameColumnHelper.accessor("name", {
    header: () => "Name",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  nameColumnHelper.accessor("symbol", {
    header: () => "Symbol",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  nameColumnHelper.accessor("decimals", {
    header: () => "Decimals",
    cell: (info) => info.renderValue(),
    size: 100,
  }),
  nameColumnHelper.accessor("isContract", {
    header: () => "isContract",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="white" fill="green" /> : ""),
    size: 100,
  }),
  nameColumnHelper.accessor("isErc20", {
    header: () => "isErc20",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="white" fill="green" /> : ""),
    size: 100,
  }),
  nameColumnHelper.accessor("isErc721", {
    header: () => "isErc721",
    cell: (info) => (info.getValue() ? <IconCircleCheck size={20} color="white" fill="green" /> : ""),
    size: 100,
  }),
];
