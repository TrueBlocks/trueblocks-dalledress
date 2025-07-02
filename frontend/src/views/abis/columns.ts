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
  // EXISTING_CODE
  {
    key: 'address',
    name: 'address',
    header: 'Address',
    label: 'Address',
    sortable: true,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'name',
    name: 'name',
    header: 'Name',
    label: 'Name',
    sortable: true,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'fileSize',
    name: 'fileSize',
    header: 'File Size',
    label: 'File Size',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
  },
  {
    key: 'nFunctions',
    name: 'nFunctions',
    header: 'Functions',
    label: 'Functions',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
  },
  {
    key: 'nEvents',
    name: 'nEvents',
    header: 'Events',
    label: 'Events',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
  },
  {
    key: 'actions',
    name: 'actions',
    header: 'Actions',
    label: 'Actions',
    sortable: false,
    editable: false,
    visible: true,
    type: 'button',
    width: '80px',
  },
];

const getFunctionsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'encoding',
    name: 'encoding',
    header: 'Encoding',
    label: 'Encoding',
    sortable: true,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'name',
    name: 'name',
    header: 'Name',
    label: 'Name',
    sortable: true,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'type',
    name: 'type',
    header: 'Type',
    label: 'Type',
    sortable: true,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'signature',
    name: 'signature',
    header: 'Signature',
    label: 'Signature',
    sortable: true,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'actions',
    name: 'actions',
    header: 'Actions',
    label: 'Actions',
    sortable: false,
    editable: false,
    visible: true,
    type: 'button',
    width: '80px',
  },
];

// EXISTING_CODE
// EXISTING_CODE
