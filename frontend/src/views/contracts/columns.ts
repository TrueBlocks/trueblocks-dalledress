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
    case types.DataFacet.EXECUTE:
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
    key: 'blockNumber',
    name: 'blockNumber',
    header: 'Block',
    label: 'Block',
    type: 'number',
    width: '120px',
  },
  {
    key: 'articulatedLog.name',
    name: 'articulatedLog.name',
    header: 'Event',
    label: 'Event',
    type: 'text',
    width: '150px',
  },
  {
    key: 'transactionHash',
    name: 'transactionHash',
    header: 'Transaction',
    label: 'Transaction',
    type: 'hash',
    width: '200px',
  },
  {
    key: 'address',
    name: 'address',
    header: 'Contract',
    label: 'Contract',
    type: 'address',
    width: '200px',
  },
  {
    key: 'contractName',
    name: 'contractName',
    header: 'Name',
    label: 'Name',
    type: 'text',
    width: '150px',
  },
  {
    key: 'compressedLog',
    name: 'compressedLog',
    header: 'Compressed Log',
    label: 'Compressed Log',
    type: 'text',
    width: '300px',
    render: renderCompressedLog,
  },
];

export function renderCompressedLog(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    // For SDK Log type, the CompressedLog would be provided by the backend
    // This function can access the compressedLog field or fall back to articulatedLog
    if (row.compressedLog && typeof row.compressedLog === 'string') {
      return row.compressedLog;
    }

    // Fallback: try to render from articulatedLog if available
    if (row.articulatedLog && typeof row.articulatedLog === 'object') {
      const articulated = row.articulatedLog as Record<string, unknown>;
      if (articulated.name) {
        return `${articulated.name}()`;
      }
    }
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
