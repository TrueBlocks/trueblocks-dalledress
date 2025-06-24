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
import { formatWeiToEther, formatWeiToGigawei } from '../../utils/ether';

// EXISTING_CODE

// Column configurations for the Exports data facets

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.STATEMENTS:
      return getStatementsColumns();
    case types.DataFacet.BALANCES:
      return getBalancesColumns();
    case types.DataFacet.TRANSFERS:
      return getTransfersColumns();
    case types.DataFacet.TRANSACTIONS:
      return getTransactionsColumns();
    case types.DataFacet.WITHDRAWALS:
      return getWithdrawalsColumns();
    case types.DataFacet.ASSETS:
      return getAssetsColumns();
    case types.DataFacet.LOGS:
      return getLogsColumns();
    case types.DataFacet.TRACES:
      return getTracesColumns();
    case types.DataFacet.RECEIPTS:
      return getReceiptsColumns();
    default:
      return [];
  }
};

const getAssetsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'address',
    header: 'Address',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'name',
    header: 'Name',
    type: 'text',
    sortable: true,
    width: 'col-name',
  },
  {
    key: 'symbol',
    header: 'Symbol',
    type: 'text',
    sortable: true,
    width: 'col-symbol',
  },
  {
    key: 'decimals',
    header: 'Decimals',
    type: 'number',
    sortable: true,
    width: 'col-decimals',
    textAlign: 'right',
  },
  {
    key: 'source',
    header: 'Source',
    type: 'text',
    sortable: true,
    width: 'col-source',
  },
  {
    key: 'tags',
    header: 'Tags',
    type: 'text',
    sortable: true,
    width: 'col-tags',
  },
  // EXISTING_CODE
];

const getBalancesColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'blockNumber',
    header: 'Block',
    type: 'number',
    sortable: true,
    width: 'col-block',
    textAlign: 'right',
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'holder',
    header: 'Holder',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'address',
    header: 'Token Address',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'symbol',
    header: 'Symbol',
    type: 'text',
    sortable: true,
    width: 'col-symbol',
  },
  {
    key: 'name',
    header: 'Token Name',
    type: 'text',
    sortable: true,
    width: 'col-name',
  },
  {
    key: 'balance',
    header: 'Balance',
    type: 'ether',
    sortable: true,
    width: 'col-balance',
    textAlign: 'right',
  },
  {
    key: 'priorBalance',
    header: 'Prior Balance',
    type: 'ether',
    sortable: true,
    width: 'col-balance',
    textAlign: 'right',
    render: renderPriorBalance,
  },
  {
    key: 'decimals',
    header: 'Decimals',
    type: 'number',
    sortable: true,
    width: 'col-decimals',
    textAlign: 'right',
  },
  {
    key: 'actions',
    header: 'Actions',
    type: 'text',
    sortable: false,
    width: '120px',
  },
  // EXISTING_CODE
];

const getLogsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'blockNumber',
    header: 'Block',
    type: 'number',
    sortable: true,
    width: 'col-block',
    textAlign: 'right',
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'logIndex',
    header: 'Log Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'blockHash',
    header: 'Block Hash',
    type: 'text',
    sortable: true,
    width: 'col-hash',
  },
  {
    key: 'transactionHash',
    header: 'Tx Hash',
    type: 'text',
    sortable: true,
    width: 'col-hash',
  },
  {
    key: 'timestamp',
    header: 'Timestamp',
    type: 'timestamp',
    sortable: true,
    width: 'col-timestamp',
  },
  {
    key: 'date',
    header: 'Date',
    type: 'text',
    sortable: true,
    width: 'col-date',
  },
  {
    key: 'address',
    header: 'Address',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'topic0',
    header: 'Topic 0',
    type: 'text',
    sortable: true,
    width: 'col-topic',
  },
  {
    key: 'topic1',
    header: 'Topic 1',
    type: 'text',
    sortable: true,
    width: 'col-topic',
  },
  {
    key: 'topic2',
    header: 'Topic 2',
    type: 'text',
    sortable: true,
    width: 'col-topic',
  },
  {
    key: 'topic3',
    header: 'Topic 3',
    type: 'text',
    sortable: true,
    width: 'col-topic',
  },
  {
    key: 'data',
    header: 'Data',
    type: 'text',
    sortable: true,
    width: 'col-data',
  },
  // EXISTING_CODE
];

const getReceiptsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'blockNumber',
    header: 'Block',
    type: 'number',
    sortable: true,
    width: 'col-block',
    textAlign: 'right',
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'transactionHash',
    header: 'Tx Hash',
    type: 'text',
    sortable: true,
    width: 'col-hash',
  },
  {
    key: 'status',
    header: 'Status',
    type: 'text',
    sortable: true,
    width: 'col-status',
  },
  {
    key: 'gasUsed',
    header: 'Gas Used',
    type: 'number',
    sortable: true,
    width: 'col-gas',
    textAlign: 'right',
  },
  {
    key: 'contractAddress',
    header: 'Contract Address',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'from',
    header: 'From',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'to',
    header: 'To',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  // EXISTING_CODE
];

const getStatementsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'blockNumber',
    header: 'Block',
    type: 'number',
    sortable: true,
    width: 'col-block',
    textAlign: 'right',
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'logIndex',
    header: 'Log Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'timestamp',
    header: 'Timestamp',
    type: 'text',
    sortable: true,
    width: 'col-timestamp',
  },
  {
    key: 'accountedFor',
    header: 'Account',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'asset',
    header: 'Asset',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'prevBal',
    header: 'Prev Balance',
    type: 'ether',
    sortable: true,
    width: 'col-balance',
    textAlign: 'right',
  },
  {
    key: 'amountNet',
    header: 'Net Amount',
    type: 'ether',
    sortable: true,
    width: 'col-amount',
    textAlign: 'right',
    render: renderStatementNetAmount,
  },
  {
    key: 'endBal',
    header: 'End Balance',
    type: 'ether',
    sortable: true,
    width: 'col-balance',
    textAlign: 'right',
  },
  {
    key: 'actions',
    header: 'Actions',
    type: 'text',
    sortable: false,
    width: '120px',
  },
  // EXISTING_CODE
];

const getTracesColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'blockNumber',
    header: 'Block',
    type: 'number',
    sortable: true,
    width: 'col-block',
    textAlign: 'right',
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'timestamp',
    header: 'Timestamp',
    type: 'timestamp',
    sortable: true,
    width: 'col-timestamp',
  },
  {
    key: 'date',
    header: 'Date',
    type: 'text',
    sortable: true,
    width: 'col-date',
  },
  {
    key: 'error',
    header: 'Error',
    type: 'text',
    sortable: true,
    width: 'col-error',
  },
  {
    key: 'action::callType',
    header: 'Call Type',
    type: 'text',
    sortable: true,
    width: 'col-calltype',
  },
  {
    key: 'action::from',
    header: 'From',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'action::to',
    header: 'To',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'action::value',
    header: 'Value',
    type: 'ether',
    sortable: true,
    width: 'col-value',
    textAlign: 'right',
  },
  {
    key: 'action::ether',
    header: 'Ether',
    type: 'ether',
    sortable: true,
    width: 'col-ether',
    textAlign: 'right',
  },
  {
    key: 'action::gas',
    header: 'Gas',
    type: 'gas',
    sortable: true,
    width: 'col-gas',
    textAlign: 'right',
  },
  {
    key: 'result::gasUsed',
    header: 'Gas Used',
    type: 'number',
    sortable: true,
    width: 'col-gas',
    textAlign: 'right',
  },
  {
    key: 'action::input',
    header: 'Input',
    type: 'text',
    sortable: true,
    width: 'col-input',
  },
  {
    key: 'result::output',
    header: 'Output',
    type: 'text',
    sortable: true,
    width: 'col-output',
  },
  // EXISTING_CODE
];

const getTransactionsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'hash',
    header: 'Hash',
    type: 'text',
    sortable: true,
    width: 'col-hash',
  },
  {
    key: 'blockNumber',
    header: 'Block',
    type: 'number',
    sortable: true,
    width: 'col-block',
    textAlign: 'right',
  },
  {
    key: 'transactionIndex',
    header: 'Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'from',
    header: 'From',
    type: 'text',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'to',
    header: 'To',
    type: 'text',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'value',
    header: 'Value',
    type: 'ether',
    sortable: true,
    width: 'col-value',
    textAlign: 'right',
  },
  {
    key: 'gas',
    header: 'Gas',
    type: 'gas',
    sortable: true,
    width: 'col-gas',
    textAlign: 'right',
  },
  {
    key: 'gasPrice',
    header: 'Gas Price',
    type: 'ether',
    sortable: true,
    width: 'col-gas-price',
    textAlign: 'right',
  },
  {
    key: 'gasCost',
    header: 'Gas Cost',
    type: 'gas',
    sortable: true,
    width: 'col-gas-price',
    textAlign: 'right',
    render: renderGasCost,
  },
  {
    key: 'actions',
    header: 'Actions',
    type: 'text',
    sortable: false,
    width: '120px',
  },
  // EXISTING_CODE
];

const getTransfersColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'blockNumber',
    header: 'Block',
    type: 'number',
    sortable: true,
    width: 'col-block',
    textAlign: 'right',
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'logIndex',
    header: 'Log Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'sender',
    header: 'Sender',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'recipient',
    header: 'Recipient',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'asset',
    header: 'Asset',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'assetSymbol',
    header: 'Symbol',
    type: 'text',
    sortable: true,
    width: 'col-symbol',
  },
  {
    key: 'amount',
    header: 'Amount',
    type: 'ether',
    sortable: true,
    width: 'col-amount',
    textAlign: 'right',
  },
  {
    key: 'actions',
    header: 'Actions',
    type: 'text',
    sortable: false,
    width: '120px',
  },
  // EXISTING_CODE
];

const getWithdrawalsColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'blockNumber',
    header: 'Block',
    type: 'number',
    sortable: true,
    width: 'col-block',
    textAlign: 'right',
  },
  {
    key: 'index',
    header: 'Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
    textAlign: 'right',
  },
  {
    key: 'validatorIndex',
    header: 'Validator Index',
    type: 'number',
    sortable: true,
    width: 'col-validator-index',
    textAlign: 'right',
  },
  {
    key: 'address',
    header: 'Address',
    type: 'address',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'timestamp',
    header: 'Timestamp',
    type: 'timestamp',
    sortable: true,
    width: 'col-timestamp',
  },
  {
    key: 'amount',
    header: 'Amount',
    type: 'text',
    sortable: true,
    width: 'col-amount',
    textAlign: 'right',
  },
  // EXISTING_CODE
];

// EXISTING_CODE
export function renderGasCost(row: Record<string, unknown>) {
  const gas = BigInt((row.gas as string) || '0');
  const gasPrice = BigInt((row.gasPrice as string) || '0');
  const gasCost = gas * gasPrice;
  return formatWeiToGigawei(gasCost.toString());
}

export function renderStatementNetAmount(row: Record<string, unknown>) {
  const amountIn = BigInt((row.amountIn as string) || '0');
  const amountOut = BigInt((row.amountOut as string) || '0');
  const netAmount = amountIn + amountOut;
  return formatWeiToEther(netAmount.toString());
}

export function renderPriorBalance(row: Record<string, unknown>) {
  const balance = (row.priorBalance as string) || '0';
  return formatWeiToEther(balance);
}

// EXISTING_CODE
