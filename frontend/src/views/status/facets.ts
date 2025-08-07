// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated. Do not edit.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const statusFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATUS,
    label: toProperCase(types.DataFacet.STATUS),
  },
  {
    id: types.DataFacet.CACHES,
    label: toProperCase(types.DataFacet.CACHES),
  },
  {
    id: types.DataFacet.CHAINS,
    label: toProperCase(types.DataFacet.CHAINS),
  },
];
