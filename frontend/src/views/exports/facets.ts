// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const exportsFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATEMENTS,
    label: toProperCase(types.DataFacet.STATEMENTS),
    isDefault: true,
  },
  {
    id: types.DataFacet.BALANCES,
    label: toProperCase(types.DataFacet.BALANCES),
    isDefault: false,
  },
  {
    id: types.DataFacet.TRANSFERS,
    label: toProperCase(types.DataFacet.TRANSFERS),
    isDefault: false,
  },
  {
    id: types.DataFacet.TRANSACTIONS,
    label: toProperCase(types.DataFacet.TRANSACTIONS),
    isDefault: false,
  },
  {
    id: types.DataFacet.WITHDRAWALS,
    label: toProperCase(types.DataFacet.WITHDRAWALS),
    isDefault: false,
  },
  {
    id: types.DataFacet.ASSETS,
    label: toProperCase(types.DataFacet.ASSETS),
    isDefault: false,
  },
  {
    id: types.DataFacet.LOGS,
    label: toProperCase(types.DataFacet.LOGS),
    isDefault: false,
  },
  {
    id: types.DataFacet.TRACES,
    label: toProperCase(types.DataFacet.TRACES),
    isDefault: false,
  },
  {
    id: types.DataFacet.RECEIPTS,
    label: toProperCase(types.DataFacet.RECEIPTS),
    isDefault: false,
  },
];

export const DEFAULT_FACET = types.DataFacet.STATEMENTS;
export const ROUTE = '/exports' as const;
