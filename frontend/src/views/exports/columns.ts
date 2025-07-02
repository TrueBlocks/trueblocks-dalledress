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
    key: 'symbol',
    name: 'symbol',
    header: 'Symbol',
    label: 'Symbol',
    sortable: true,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'decimals',
    name: 'decimals',
    header: 'Decimals',
    label: 'Decimals',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
  },
  {
    key: 'source',
    name: 'source',
    header: 'Source',
    label: 'Source',
    sortable: true,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'tags',
    name: 'tags',
    header: 'Tags',
    label: 'Tags',
    sortable: true,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
];

const getBalancesColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'blockNumber',
    name: 'blockNumber',
    header: 'Block Number',
    label: 'Block Number',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
  },
  {
    key: 'transactionIndex',
    name: 'transactionIndex',
    header: 'Transaction Index',
    label: 'Transaction Index',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
  },
  {
    key: 'holder',
    name: 'holder',
    header: 'Holder',
    label: 'Holder',
    sortable: true,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
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
    key: 'symbol',
    name: 'symbol',
    header: 'Symbol',
    label: 'Symbol',
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
    key: 'balance',
    name: 'balance',
    header: 'Balance',
    label: 'Balance',
    sortable: true,
    type: 'ether',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'priorBalance',
    name: 'priorBalance',
    header: 'Prior Balance',
    label: 'Prior Balance',
    sortable: true,
    type: 'ether',
    width: '120px',
    textAlign: 'left',
    render: renderPriorBalance,
  },
  {
    key: 'decimals',
    name: 'decimals',
    header: 'Decimals',
    label: 'Decimals',
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
    type: 'text',
    width: '120px',
    textAlign: 'left',
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
    sortable: false,
    type: 'datetime',
    width: '120px',
    textAlign: 'left',
    render: renderDate,
  },
  {
    key: 'address',
    name: 'address',
    header: 'Address',
    label: 'Address',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'topics',
    name: 'topics',
    header: 'Topics',
    label: 'Topics',
    sortable: false,
    type: 'topic',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'data',
    name: 'data',
    header: 'Data',
    label: 'Data',
    sortable: false,
    type: 'bytes',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'compressedLog',
    name: 'compressedLog',
    header: 'Compressed Log',
    label: 'Compressed Log',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
    render: renderCompressedLog,
  },
  {
    key: 'isNFT',
    name: 'isNFT',
    header: 'N F T',
    label: 'N F T',
    sortable: false,
    type: 'checkbox',
    width: '80px',
    textAlign: 'center',
  },
];

const getReceiptsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'cumulativeGasUsed',
    name: 'cumulativeGasUsed',
    header: 'Cumulative Gas Used',
    label: 'Cumulative Gas Used',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'effectiveGasPrice',
    name: 'effectiveGasPrice',
    header: 'Effective Gas Price',
    label: 'Effective Gas Price',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'from',
    name: 'from',
    header: 'From',
    label: 'From',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'to',
    name: 'to',
    header: 'To',
    label: 'To',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'blockHash',
    name: 'blockHash',
    header: 'Block Hash',
    label: 'Block Hash',
    sortable: false,
    type: 'hash',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'blockNumber',
    name: 'blockNumber',
    header: 'Block Number',
    label: 'Block Number',
    sortable: false,
    type: 'blknum',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'contractAddress',
    name: 'contractAddress',
    header: 'Contract Address',
    label: 'Contract Address',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'gasUsed',
    name: 'gasUsed',
    header: 'Gas Used',
    label: 'Gas Used',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'isError',
    name: 'isError',
    header: 'Error',
    label: 'Error',
    sortable: false,
    type: 'checkbox',
    width: '80px',
    textAlign: 'center',
  },
  {
    key: 'status',
    name: 'status',
    header: 'Status',
    label: 'Status',
    sortable: false,
    type: 'value',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'transactionHash',
    name: 'transactionHash',
    header: 'Transaction Hash',
    label: 'Transaction Hash',
    sortable: false,
    type: 'hash',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'transactionIndex',
    name: 'transactionIndex',
    header: 'Transaction Index',
    label: 'Transaction Index',
    sortable: false,
    type: 'txnum',
    width: '120px',
    textAlign: 'left',
  },
];

const getStatementsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'date',
    name: 'date',
    header: 'Date',
    label: 'Date',
    sortable: false,
    type: 'datetime',
    width: '120px',
    textAlign: 'left',
    render: renderDate,
  },
  {
    key: 'asset',
    name: 'asset',
    header: 'Asset',
    label: 'Asset',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'symbol',
    name: 'symbol',
    header: 'Symbol',
    label: 'Symbol',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'decimals',
    name: 'decimals',
    header: 'Decimals',
    label: 'Decimals',
    sortable: false,
    type: 'value',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'spotPrice',
    name: 'spotPrice',
    header: 'Spot Price',
    label: 'Spot Price',
    sortable: false,
    type: 'float',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'priceSource',
    name: 'priceSource',
    header: 'Price Source',
    label: 'Price Source',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'accountedFor',
    name: 'accountedFor',
    header: 'Accounted For',
    label: 'Accounted For',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'sender',
    name: 'sender',
    header: 'Sender',
    label: 'Sender',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'recipient',
    name: 'recipient',
    header: 'Recipient',
    label: 'Recipient',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'begBal',
    name: 'begBal',
    header: 'Beg Bal',
    label: 'Beg Bal',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'amountNet',
    name: 'amountNet',
    header: 'Amount Net',
    label: 'Amount Net',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
    render: renderAmountNet,
  },
  {
    key: 'endBal',
    name: 'endBal',
    header: 'End Bal',
    label: 'End Bal',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'reconciled',
    name: 'reconciled',
    header: 'Reconciled',
    label: 'Reconciled',
    sortable: false,
    type: 'checkbox',
    width: '80px',
    textAlign: 'center',
  },
  {
    key: 'totalIn',
    name: 'totalIn',
    header: 'Total In',
    label: 'Total In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'amountIn',
    name: 'amountIn',
    header: 'Amount In',
    label: 'Amount In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'internalIn',
    name: 'internalIn',
    header: 'Internal In',
    label: 'Internal In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'selfDestructIn',
    name: 'selfDestructIn',
    header: 'Self Destruct In',
    label: 'Self Destruct In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'minerBaseRewardIn',
    name: 'minerBaseRewardIn',
    header: 'Miner Base Reward In',
    label: 'Miner Base Reward In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'minerNephewRewardIn',
    name: 'minerNephewRewardIn',
    header: 'Miner Nephew Reward In',
    label: 'Miner Nephew Reward In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'minerTxFeeIn',
    name: 'minerTxFeeIn',
    header: 'Miner Tx Fee In',
    label: 'Miner Tx Fee In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'minerUncleRewardIn',
    name: 'minerUncleRewardIn',
    header: 'Miner Uncle Reward In',
    label: 'Miner Uncle Reward In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'correctBegBalIn',
    name: 'correctBegBalIn',
    header: 'Correct Beg Bal In',
    label: 'Correct Beg Bal In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'correctAmountIn',
    name: 'correctAmountIn',
    header: 'Correct Amount In',
    label: 'Correct Amount In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'correctEndBalIn',
    name: 'correctEndBalIn',
    header: 'Correct End Bal In',
    label: 'Correct End Bal In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'prefundIn',
    name: 'prefundIn',
    header: 'Prefund In',
    label: 'Prefund In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'totalOut',
    name: 'totalOut',
    header: 'Total Out',
    label: 'Total Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'amountOut',
    name: 'amountOut',
    header: 'Amount Out',
    label: 'Amount Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'internalOut',
    name: 'internalOut',
    header: 'Internal Out',
    label: 'Internal Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'correctBegBalOut',
    name: 'correctBegBalOut',
    header: 'Correct Beg Bal Out',
    label: 'Correct Beg Bal Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'correctAmountOut',
    name: 'correctAmountOut',
    header: 'Correct Amount Out',
    label: 'Correct Amount Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'correctEndBalOut',
    name: 'correctEndBalOut',
    header: 'Correct End Bal Out',
    label: 'Correct End Bal Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'selfDestructOut',
    name: 'selfDestructOut',
    header: 'Self Destruct Out',
    label: 'Self Destruct Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'gasOut',
    name: 'gasOut',
    header: 'Gas Out',
    label: 'Gas Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'prevBal',
    name: 'prevBal',
    header: 'Prev Bal',
    label: 'Prev Bal',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'begBalDiff',
    name: 'begBalDiff',
    header: 'Beg Bal Diff',
    label: 'Beg Bal Diff',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'endBalDiff',
    name: 'endBalDiff',
    header: 'End Bal Diff',
    label: 'End Bal Diff',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'endBalCalc',
    name: 'endBalCalc',
    header: 'End Bal Calc',
    label: 'End Bal Calc',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'correctingReasons',
    name: 'correctingReasons',
    header: 'Correcting Reasons',
    label: 'Correcting Reasons',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
];

const getTracesColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'date',
    name: 'date',
    header: 'Date',
    label: 'Date',
    sortable: false,
    type: 'datetime',
    width: '120px',
    textAlign: 'left',
    render: renderDate,
  },
  {
    key: 'type',
    name: 'type',
    header: 'Type',
    label: 'Type',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'compressedTrace',
    name: 'compressedTrace',
    header: 'Compressed Trace',
    label: 'Compressed Trace',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
    render: renderCompressedTrace,
  },
  {
    key: 'error',
    name: 'error',
    header: 'Error',
    label: 'Error',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
];

const getTransactionsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'gasUsed',
    name: 'gasUsed',
    header: 'Gas Used',
    label: 'Gas Used',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'maxPriorityFeePerGas',
    name: 'maxPriorityFeePerGas',
    header: 'Max Priority Fee Per Gas',
    label: 'Max Priority Fee Per Gas',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'type',
    name: 'type',
    header: 'Type',
    label: 'Type',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
  },
  {
    key: 'maxFeePerGas',
    name: 'maxFeePerGas',
    header: 'Max Fee Per Gas',
    label: 'Max Fee Per Gas',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'nonce',
    name: 'nonce',
    header: 'Nonce',
    label: 'Nonce',
    sortable: false,
    type: 'value',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'date',
    name: 'date',
    header: 'Date',
    label: 'Date',
    sortable: false,
    type: 'datetime',
    width: '120px',
    textAlign: 'left',
    render: renderDate,
  },
  {
    key: 'from',
    name: 'from',
    header: 'From',
    label: 'From',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'to',
    name: 'to',
    header: 'To',
    label: 'To',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'value',
    name: 'value',
    header: 'Value',
    label: 'Value',
    sortable: false,
    type: 'wei',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'ether',
    name: 'ether',
    header: 'Ether',
    label: 'Ether',
    sortable: false,
    type: 'ether',
    width: '120px',
    textAlign: 'left',
    render: renderEther,
  },
  {
    key: 'gas',
    name: 'gas',
    header: 'Gas',
    label: 'Gas',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'gasPrice',
    name: 'gasPrice',
    header: 'Gas Price',
    label: 'Gas Price',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'gasCost',
    name: 'gasCost',
    header: 'Gas Cost',
    label: 'Gas Cost',
    sortable: false,
    type: 'gas',
    width: '120px',
    textAlign: 'left',
    render: renderGasCost,
  },
  {
    key: 'input',
    name: 'input',
    header: 'Input',
    label: 'Input',
    sortable: false,
    type: 'bytes',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'hasToken',
    name: 'hasToken',
    header: 'Has Token',
    label: 'Has Token',
    sortable: false,
    type: 'checkbox',
    width: '80px',
    textAlign: 'center',
  },
  {
    key: 'isError',
    name: 'isError',
    header: 'Error',
    label: 'Error',
    sortable: false,
    type: 'checkbox',
    width: '80px',
    textAlign: 'center',
  },
  {
    key: 'compressedTx',
    name: 'compressedTx',
    header: 'Compressed Tx',
    label: 'Compressed Tx',
    sortable: false,
    type: 'text',
    width: '200px',
    textAlign: 'left',
    render: renderCompressedTx,
  },
];

const getTransfersColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'blockNumber',
    name: 'blockNumber',
    header: 'Block Number',
    label: 'Block Number',
    sortable: false,
    type: 'blknum',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'transactionIndex',
    name: 'transactionIndex',
    header: 'Transaction Index',
    label: 'Transaction Index',
    sortable: false,
    type: 'txnum',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'logIndex',
    name: 'logIndex',
    header: 'Log Index',
    label: 'Log Index',
    sortable: false,
    type: 'lognum',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'date',
    name: 'date',
    header: 'Date',
    label: 'Date',
    sortable: false,
    type: 'datetime',
    width: '120px',
    textAlign: 'left',
    render: renderDate,
  },
  {
    key: 'holder',
    name: 'holder',
    header: 'Holder',
    label: 'Holder',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'asset',
    name: 'asset',
    header: 'Asset',
    label: 'Asset',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'decimals',
    name: 'decimals',
    header: 'Decimals',
    label: 'Decimals',
    sortable: false,
    type: 'number',
    width: '120px',
    textAlign: 'right',
  },
  {
    key: 'sender',
    name: 'sender',
    header: 'Sender',
    label: 'Sender',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'recipient',
    name: 'recipient',
    header: 'Recipient',
    label: 'Recipient',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'amountIn',
    name: 'amountIn',
    header: 'Amount In',
    label: 'Amount In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'internalIn',
    name: 'internalIn',
    header: 'Internal In',
    label: 'Internal In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'minerBaseRewardIn',
    name: 'minerBaseRewardIn',
    header: 'Miner Base Reward In',
    label: 'Miner Base Reward In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'minerNephewRewardIn',
    name: 'minerNephewRewardIn',
    header: 'Miner Nephew Reward In',
    label: 'Miner Nephew Reward In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'minerTxFeeIn',
    name: 'minerTxFeeIn',
    header: 'Miner Tx Fee In',
    label: 'Miner Tx Fee In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'minerUncleRewardIn',
    name: 'minerUncleRewardIn',
    header: 'Miner Uncle Reward In',
    label: 'Miner Uncle Reward In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'prefundIn',
    name: 'prefundIn',
    header: 'Prefund In',
    label: 'Prefund In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'selfDestructIn',
    name: 'selfDestructIn',
    header: 'Self Destruct In',
    label: 'Self Destruct In',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'amountOut',
    name: 'amountOut',
    header: 'Amount Out',
    label: 'Amount Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'internalOut',
    name: 'internalOut',
    header: 'Internal Out',
    label: 'Internal Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'gasOut',
    name: 'gasOut',
    header: 'Gas Out',
    label: 'Gas Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'selfDestructOut',
    name: 'selfDestructOut',
    header: 'Self Destruct Out',
    label: 'Self Destruct Out',
    sortable: false,
    type: 'int256',
    width: '120px',
    textAlign: 'left',
  },
];

const getWithdrawalsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'address',
    name: 'address',
    header: 'Address',
    label: 'Address',
    sortable: false,
    type: 'address',
    width: '340px',
    textAlign: 'left',
    readOnly: true,
  },
  {
    key: 'amount',
    name: 'amount',
    header: 'Amount',
    label: 'Amount',
    sortable: false,
    type: 'wei',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'ether',
    name: 'ether',
    header: 'Ether',
    label: 'Ether',
    sortable: false,
    type: 'ether',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'blockNumber',
    name: 'blockNumber',
    header: 'Block Number',
    label: 'Block Number',
    sortable: false,
    type: 'blknum',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'index',
    name: 'index',
    header: 'Index',
    label: 'Index',
    sortable: false,
    type: 'value',
    width: '120px',
    textAlign: 'left',
  },
  {
    key: 'date',
    name: 'date',
    header: 'Date',
    label: 'Date',
    sortable: false,
    type: 'datetime',
    width: '120px',
    textAlign: 'left',
    render: renderDate,
  },
  {
    key: 'validatorIndex',
    name: 'validatorIndex',
    header: 'Validator Index',
    label: 'Validator Index',
    sortable: false,
    type: 'value',
    width: '120px',
    textAlign: 'left',
  },
];

export function renderAmountNet(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    const amountIn = BigInt((row.amountIn as string) || '0');
    const amountOut = BigInt((row.amountOut as string) || '0');
    const netAmount = amountIn + amountOut;
    return formatWeiToEther(netAmount.toString());
    // EXISTING_CODE
  }
  return '';
}

export function renderCompressedLog(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    const log = row['articulatedLog'] as unknown as types.Function;
    return log?.name;
    // EXISTING_CODE
  }
  return '';
}

export function renderCompressedTrace(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    return 'renderCompressedTrace';
    // EXISTING_CODE
  }
  return '';
}

export function renderCompressedTx(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    return 'renderCompressedTx';
    // EXISTING_CODE
  }
  return '';
}

export function renderDate(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    var timestamp = row.timestamp as string | number | undefined;
    if (timestamp === undefined) {
      if (row.transaction) {
        const tx = row.transaction as types.Transaction | undefined;
        if (tx != null) {
          timestamp = tx.timestamp as string | number | undefined;
        }
      }
    }
    const blockNumber = row.blockNumber as string | number | undefined;
    const transactionIndex = row.transactionIndex as
      | string
      | number
      | undefined;
    const transactionHash = row.transactionHash as string | undefined;
    const blockHash = row.blockHash as string | undefined;
    const node = row.node as string | undefined;

    // Format date
    let dateStr = '';
    if (timestamp) {
      const date = new Date(Number(timestamp) * 1000);
      dateStr = date.toISOString().replace('T', ' ').substring(0, 19);
    }

    // Compose extra info
    const parts: string[] = [];
    if (blockNumber !== undefined) parts.push(`Block: ${blockNumber}`);
    if (transactionIndex !== undefined)
      parts.push(`TxIdx: ${transactionIndex}`);
    if (transactionHash) parts.push(`Tx: ${transactionHash.slice(0, 10)}…`);
    if (blockHash) parts.push(`BlkHash: ${blockHash.slice(0, 10)}…`);
    if (node) parts.push(`Node: ${node}`);

    return [dateStr, ...parts].join(' | ');
    // EXISTING_CODE
  }
  return '';
}

export function renderEther(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    return 'renderEther';
    // EXISTING_CODE
  }
  return '';
}

export function renderGasCost(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    const gas = BigInt((row.gas as string) || '0');
    const gasPrice = BigInt((row.gasPrice as string) || '0');
    const gasCost = gas * gasPrice;
    return formatWeiToGigawei(gasCost.toString());
    // EXISTING_CODE
  }
  return '';
}

export function renderPriorBalance(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    const balance = (row.priorBalance as string) || '0';
    return formatWeiToEther(balance);
    // EXISTING_CODE
  }
  return '';
}

export function renderStatements(row: Record<string, unknown>) {
  if (row != null) {
    // EXISTING_CODE
    return 'renderStatements';
    // EXISTING_CODE
  }
  return '';
}

// EXISTING_CODE
// EXISTING_CODE
