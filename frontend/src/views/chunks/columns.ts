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
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'blkrange',
    width: '150px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'magic',
    name: 'magic',
    header: 'Magic',
    label: 'Magic',
    sortable: true,
    type: 'text',
    width: '150px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'hash',
    name: 'hash',
    header: 'Hash',
    label: 'Hash',
    sortable: true,
    type: 'hash',
    width: '150px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'nBlooms',
    name: 'nBlooms',
    header: 'Blooms',
    label: 'Blooms',
    sortable: true,
    type: 'number',
    width: '150px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'nInserted',
    name: 'nInserted',
    header: 'Inserted',
    label: 'Inserted',
    sortable: true,
    type: 'number',
    width: '150px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'size',
    name: 'size',
    header: 'Size',
    label: 'Size',
    sortable: true,
    type: 'number',
    width: '150px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'byteWidth',
    name: 'byteWidth',
    header: 'Byte Width',
    label: 'Byte Width',
    sortable: true,
    type: 'number',
    width: '150px',
    textAlign: 'right',
    readOnly: true,
  },
];

const getIndexColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'blkrange',
    width: '150px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'magic',
    name: 'magic',
    header: 'Magic',
    label: 'Magic',
    sortable: true,
    type: 'text',
    width: '150px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'hash',
    name: 'hash',
    header: 'Hash',
    label: 'Hash',
    sortable: true,
    type: 'hash',
    width: '150px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'nAddresses',
    name: 'nAddresses',
    header: 'Addresses',
    label: 'Addresses',
    sortable: true,
    type: 'number',
    width: '150px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'nAppearances',
    name: 'nAppearances',
    header: 'Appearances',
    label: 'Appearances',
    sortable: true,
    type: 'number',
    width: '150px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'size',
    name: 'size',
    header: 'Size',
    label: 'Size',
    sortable: true,
    type: 'number',
    width: '150px',
    textAlign: 'right',
    readOnly: true,
  },
];

const getManifestColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'version',
    name: 'version',
    header: 'Version',
    label: 'Version',
    sortable: true,
    type: 'text',
    width: '100px',
    textAlign: 'left',
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
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'specification',
    name: 'specification',
    header: 'Specification',
    label: 'Specification',
    sortable: true,
    type: 'ipfshash',
    width: '100px',
    textAlign: 'left',
    readOnly: true,
  },
];

const getStatsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'blkrange',
    width: '150px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'nAddrs',
    name: 'nAddrs',
    header: 'Addrs',
    label: 'Addrs',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'nApps',
    name: 'nApps',
    header: 'Apps',
    label: 'Apps',
    sortable: true,
    type: 'number',
    width: '100px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'nBlocks',
    name: 'nBlocks',
    header: 'Blocks',
    label: 'Blocks',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'nBlooms',
    name: 'nBlooms',
    header: 'Blooms',
    label: 'Blooms',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'recWid',
    name: 'recWid',
    header: 'Rec Wid',
    label: 'Rec Wid',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'bloomSz',
    name: 'bloomSz',
    header: 'Bloom Sz',
    label: 'Bloom Sz',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'chunkSz',
    name: 'chunkSz',
    header: 'Chunk Sz',
    label: 'Chunk Sz',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
    readOnly: true,
  },
  {
    key: 'addrsPerBlock',
    name: 'addrsPerBlock',
    header: 'Addrs Per Block',
    label: 'Addrs Per Block',
    sortable: true,
    type: 'float64',
    width: '100px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'appsPerBlock',
    name: 'appsPerBlock',
    header: 'Apps Per Block',
    label: 'Apps Per Block',
    sortable: true,
    type: 'float64',
    width: '100px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'appsPerAddr',
    name: 'appsPerAddr',
    header: 'Apps Per Addr',
    label: 'Apps Per Addr',
    sortable: true,
    type: 'float64',
    width: '100px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'ratio',
    name: 'ratio',
    header: 'Ratio',
    label: 'Ratio',
    sortable: true,
    type: 'float64',
    width: '100px',
    textAlign: 'left',
    readOnly: true,
  },
];

// EXISTING_CODE
// EXISTING_CODE
