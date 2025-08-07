// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated. Do not edit.
 */
import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const abisFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.DOWNLOADED,
    label: toProperCase(types.DataFacet.DOWNLOADED),
  },
  {
    id: types.DataFacet.KNOWN,
    label: toProperCase(types.DataFacet.KNOWN),
  },
  {
    id: types.DataFacet.FUNCTIONS,
    label: toProperCase(types.DataFacet.FUNCTIONS),
  },
  {
    id: types.DataFacet.EVENTS,
    label: toProperCase(types.DataFacet.EVENTS),
  },
];
