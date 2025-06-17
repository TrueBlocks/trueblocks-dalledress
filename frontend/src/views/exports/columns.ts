import { FormField } from '@components';
import { types } from '@models';

// import { formatWeiToEther, formatWeiToGigawei } from '../../utils/ether';

export const getColumnsForExports = (
  listKind: types.ListKind | string,
): FormField[] => {
  switch (listKind) {
    case types.ListKind.TRANSACTIONS:
    case 'Transactions':
      return getTransactionColumns();
    case types.ListKind.STATEMENTS:
    case 'Statements':
      return getStatementColumns();
    case types.ListKind.TRANSFERS:
    case 'Transfers':
      return getTransferColumns();
    case types.ListKind.BALANCES:
    case 'Balances':
      return getBalanceColumns();
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
  },
  {
    key: 'transactionIndex',
    header: 'Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
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
    type: 'text',
    sortable: true,
    width: 'col-value',
  },
  {
    key: 'gas',
    header: 'Gas',
    type: 'text',
    sortable: true,
    width: 'col-gas',
  },
  {
    key: 'gasPrice',
    header: 'Gas Price',
    type: 'text',
    sortable: true,
    width: 'col-gas-price',
  },
  {
    key: 'gasCost',
    header: 'Gas Cost',
    type: 'text',
    sortable: true,
    width: 'col-gas-price',
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
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
  },
  {
    key: 'logIndex',
    header: 'Log Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
  },
  {
    key: 'accountedFor',
    header: 'Account',
    type: 'text',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'asset',
    header: 'Asset',
    type: 'text',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'prevBal',
    header: 'Prev Balance',
    type: 'text',
    sortable: true,
    width: 'col-balance',
  },
  {
    key: 'amountNet',
    header: 'Net Amount',
    type: 'text',
    sortable: true,
    width: 'col-amount',
  },
  {
    key: 'endBal',
    header: 'End Balance',
    type: 'text',
    sortable: true,
    width: 'col-balance',
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
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
  },
  {
    key: 'logIndex',
    header: 'Log Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
  },
  {
    key: 'sender',
    header: 'Sender',
    type: 'text',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'recipient',
    header: 'Recipient',
    type: 'text',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'asset',
    header: 'Asset',
    type: 'text',
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
    type: 'text',
    sortable: true,
    width: 'col-amount',
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
  },
  {
    key: 'transactionIndex',
    header: 'Tx Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
  },
  {
    key: 'logIndex',
    header: 'Log Index',
    type: 'number',
    sortable: true,
    width: 'col-index',
  },
  {
    key: 'address',
    header: 'Address',
    type: 'text',
    sortable: true,
    width: 'col-address',
  },
  {
    key: 'assetAddr',
    header: 'Asset',
    type: 'text',
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
    key: 'balance',
    header: 'Balance',
    type: 'text',
    sortable: true,
    width: 'col-balance',
  },
  {
    key: 'diff',
    header: 'Difference',
    type: 'text',
    sortable: true,
    width: 'col-amount',
  },
  {
    key: 'actions',
    header: 'Actions',
    type: 'text',
    sortable: false,
    width: '120px',
  },
];
