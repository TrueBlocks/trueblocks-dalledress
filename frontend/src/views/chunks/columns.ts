// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
import { FormField } from '@components';
import { types } from '@models';

// EXISTING_CODE
// EXISTING_CODE

// Column configurations for the Chunks data facets

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.STATS:
      return getStatsColumns();
    case types.DataFacet.INDEX:
      return getIndexColumns();
    case types.DataFacet.BLOOMS:
      return getBloomsColumns();
    case types.DataFacet.MANIFEST:
      return getManifestColumns();
    default:
      return [];
  }
};

const getBloomsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'text',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'magic',
    name: 'magic',
    header: 'Magic',
    label: 'Magic',
    sortable: true,
    type: 'text',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'hash',
    name: 'hash',
    header: 'Hash',
    label: 'Hash',
    sortable: true,
    type: 'text',
    width: '300px',
    readOnly: true,
  },
  {
    key: 'nBlooms',
    name: 'nBlooms',
    header: 'Blooms',
    label: 'Blooms',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nInserted',
    name: 'nInserted',
    header: 'Inserted',
    label: 'Inserted',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'size',
    name: 'size',
    header: 'Size',
    label: 'Size',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'byteWidth',
    name: 'byteWidth',
    header: 'Byte Width',
    label: 'Byte Width',
    sortable: true,
    type: 'number',
    width: '110px',
    readOnly: true,
    textAlign: 'right',
  },
  // EXISTING_CODE
];

const getIndexColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'text',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'magic',
    name: 'magic',
    header: 'Magic',
    label: 'Magic',
    sortable: true,
    type: 'text',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'hash',
    name: 'hash',
    header: 'Hash',
    label: 'Hash',
    sortable: true,
    type: 'text',
    width: '300px',
    readOnly: true,
  },
  {
    key: 'nAddresses',
    name: 'nAddresses',
    header: 'Addresses',
    label: 'Addresses',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nAppearances',
    name: 'nAppearances',
    header: 'Appearances',
    label: 'Appearances',
    sortable: true,
    type: 'number',
    width: '130px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'size',
    name: 'size',
    header: 'Size',
    label: 'Size',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  // EXISTING_CODE
];

const getManifestColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'version',
    name: 'version',
    header: 'Version',
    label: 'Version',
    sortable: true,
    type: 'text',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'chain',
    name: 'chain',
    header: 'Chain',
    label: 'Chain',
    sortable: true,
    type: 'text',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'specification',
    name: 'specification',
    header: 'Specification',
    label: 'Specification',
    sortable: true,
    type: 'text',
    width: '200px',
    readOnly: true,
  },
  // EXISTING_CODE
];

const getStatsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'text',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'nAddrs',
    name: 'nAddrs',
    header: 'Addresses',
    label: 'Addresses',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nApps',
    name: 'nApps',
    header: 'Apps',
    label: 'Apps',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nBlocks',
    name: 'nBlocks',
    header: 'Blocks',
    label: 'Blocks',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nBloomsHit',
    name: 'nBloomsHit',
    header: 'Blooms Hit',
    label: 'Blooms Hit',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nBloomsMiss',
    name: 'nBloomsMiss',
    header: 'Blooms Miss',
    label: 'Blooms Miss',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'ratio',
    name: 'ratio',
    header: 'Hit Ratio',
    label: 'Hit Ratio',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  // EXISTING_CODE
];

// EXISTING_CODE
// EXISTING_CODE
