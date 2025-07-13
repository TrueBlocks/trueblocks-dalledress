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

// Column configurations for the Contracts data facets

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.DASHBOARD:
      return getContractsColumns();
    case types.DataFacet.DYNAMIC:
      return getContractsColumns();
    case types.DataFacet.EVENTS:
      return getLogsColumns();
    default:
      return [];
  }
};

const getContractsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'address',
    name: 'address',
    header: 'Address',
    label: 'Address',
    type: 'address',
    width: '340px',
    readOnly: true,
  },
  {
    key: 'name',
    name: 'name',
    header: 'Name',
    label: 'Name',
    type: 'text',
    width: '200px',
  },
  {
    key: 'lastUpdated',
    name: 'lastUpdated',
    header: 'Last Updated',
    label: 'Last Updated',
    type: 'timestamp',
    width: '120px',
  },
  {
    key: 'date',
    name: 'date',
    header: 'Date',
    label: 'Date',
    type: 'datetime',
    width: '120px',
  },
  {
    key: 'errorCount',
    name: 'errorCount',
    header: 'Error Count',
    label: 'Error Count',
    type: 'number',
    width: '80px',
  },
];

const getLogsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'date',
    name: 'date',
    header: 'Date',
    label: 'Date',
    type: 'datetime',
    width: '120px',
    render: renderDate,
  },
  {
    key: 'address',
    name: 'address',
    header: 'Address',
    label: 'Address',
    type: 'address',
    width: '340px',
    readOnly: true,
  },
  {
    key: 'topics',
    name: 'topics',
    header: 'Topics',
    label: 'Topics',
    type: 'topic',
    width: '120px',
  },
  {
    key: 'data',
    name: 'data',
    header: 'Data',
    label: 'Data',
    type: 'bytes',
    width: '120px',
  },
  {
    key: 'compressedLog',
    name: 'compressedLog',
    header: 'Compressed Log',
    label: 'Compressed Log',
    type: 'text',
    width: '200px',
    render: renderCompressedLog,
  },
  {
    key: 'isNFT',
    name: 'isNFT',
    header: 'N F T',
    label: 'N F T',
    type: 'checkbox',
    width: '80px',
  },
];

export function renderCompressedLog(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    // EXISTING_CODE
  }
  return '';
}

export function renderDate(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    // EXISTING_CODE
  }
  return '';
}

// EXISTING_CODE
// EXISTING_CODE
