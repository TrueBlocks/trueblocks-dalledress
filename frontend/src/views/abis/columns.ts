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

// Column configurations for the Abis data facets

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.DOWNLOADED:
      return getAbisColumns();
    case types.DataFacet.KNOWN:
      return getAbisColumns();
    case types.DataFacet.FUNCTIONS:
      return getFunctionsColumns();
    case types.DataFacet.EVENTS:
      return getFunctionsColumns();
    default:
      return [];
  }
};

const getAbisColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'address',
    header: 'Address',
    sortable: true,
    type: 'text',
  },
  {
    key: 'name',
    header: 'Name',
    sortable: true,
    type: 'text',
  },
  {
    key: 'fileSize',
    header: 'File Size',
    sortable: true,
    type: 'number',
    textAlign: 'right',
  },
  {
    key: 'nFunctions',
    header: 'Functions',
    sortable: true,
    type: 'number',
    textAlign: 'right',
  },
  {
    key: 'nEvents',
    header: 'Events',
    sortable: true,
    type: 'number',
    textAlign: 'right',
  },
  {
    key: 'actions',
    header: 'Actions',
    sortable: false,
    type: 'text',
    width: '120px',
  },
  // EXISTING_CODE
];

const getFunctionsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'encoding',
    header: 'Encoding',
    type: 'text',
    sortable: true,
    width: 'col-encoding',
  },
  {
    key: 'name',
    header: 'Name',
    type: 'text',
    sortable: true,
    width: 'col-name',
  },
  {
    key: 'type',
    header: 'Type',
    type: 'text',
    sortable: true,
    width: 'col-type',
  },
  {
    key: 'signature',
    header: 'Signature',
    type: 'text',
    sortable: true,
    width: 'col-signature',
  },
  // EXISTING_CODE
];

// EXISTING_CODE
// EXISTING_CODE
