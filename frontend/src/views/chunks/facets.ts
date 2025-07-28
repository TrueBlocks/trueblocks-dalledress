// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated. Do not edit.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const chunksFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATS,
    label: toProperCase(types.DataFacet.STATS),
  },
  {
    id: types.DataFacet.INDEX,
    label: toProperCase(types.DataFacet.INDEX),
  },
  {
    id: types.DataFacet.BLOOMS,
    label: toProperCase(types.DataFacet.BLOOMS),
  },
  {
    id: types.DataFacet.MANIFEST,
    label: toProperCase(types.DataFacet.MANIFEST),
  },
];
