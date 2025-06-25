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

export const monitorsFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.MONITORS,
    label: toProperCase(types.DataFacet.MONITORS),
    isDefault: true,
  },
];

export const MONITORS_DEFAULT_FACET = types.DataFacet.MONITORS;
export const MONITORS_ROUTE = '/monitors' as const;
