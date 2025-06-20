import { FormField } from '@components';
import { types } from '@models';

import { formatWeiToEther, formatWeiToGigawei } from '../../utils/ether';

export const getColumnsForExports = (
  dataFacet: types.DataFacet | string,
): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.STATEMENTS:
      return getStatementColumns();
    case types.DataFacet.TRANSFERS:
      return getTransferColumns();
    case types.DataFacet.BALANCES:
      return getBalanceColumns();
    case types.DataFacet.TRANSACTIONS:
      return getTransactionColumns();
    default:
      return getTransactionColumns();
  }
};

const getTransactionColumns = (): FormField[] => [
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
    render: (row) => {
      const gas = BigInt((row.gas as string) || '0');
      const gasPrice = BigInt((row.gasPrice as string) || '0');
      const gasCost = gas * gasPrice;
      return formatWeiToGigawei(gasCost.toString());
    },
  },
  {
    key: 'actions',
    header: 'Actions',
    type: 'text',
    sortable: false,
    width: '120px',
  },
];

const getStatementColumns = (): FormField[] => [
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
    render: (row) => {
      const amountIn = BigInt((row.amountIn as string) || '0');
      const amountOut = BigInt((row.amountOut as string) || '0');
      const netAmount = amountIn + amountOut;
      return formatWeiToEther(netAmount.toString());
    },
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
];

const getTransferColumns = (): FormField[] => [
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
];

const getBalanceColumns = (): FormField[] => [
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
    render: (row) => {
      const balance = (row.priorBalance as string) || '0';
      return formatWeiToEther(balance);
    },
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
];
