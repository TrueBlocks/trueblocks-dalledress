import { types } from '@models';

import {
  type MetaOverlay,
  buildDetailPanelFromColumns,
} from '../utils/detailPanel';
import { ROUTE, getColumns } from './columns';

// === Detail Panel Meta Definitions (moved from columns.ts) ===
const statementsMeta: Record<string, MetaOverlay> = {
  date: {
    section: 'Statement',
    detailLabel: 'Timestamp',
    detailOrder: 1,
  },
  begBal: {
    section: 'Statement',
    detailLabel: 'Begin Balance',
    detailOrder: 2,
  },
  amountIn: {
    section: 'Statement',
    detailLabel: 'In',
    detailOrder: 3,
  },
  amountOut: {
    section: 'Statement',
    detailLabel: 'Out',
    detailOrder: 4,
  },
  amountNet: {
    section: 'Statement',
    detailLabel: 'Net',
    detailOrder: 5,
  },
  endBal: {
    section: 'Statement',
    detailLabel: 'End Balance',
    detailOrder: 6,
  },
  symbol: {
    section: 'Statement',
    detailLabel: 'Symbol',
    detailOrder: 7,
  },
};

const balancesMeta: Record<string, MetaOverlay> = {
  address: {
    section: 'Balance',
    detailLabel: 'Address',
    detailFormat: 'address' as const,
    detailOrder: 1,
  },
  balance: {
    section: 'Balance',
    detailLabel: 'Balance',
    detailOrder: 2,
  },
  symbol: {
    section: 'Balance',
    detailLabel: 'Symbol',
    detailOrder: 3,
  },
};

const transfersMeta: Record<string, MetaOverlay> = {
  from: {
    section: 'Transfer',
    detailLabel: 'From',
    detailFormat: 'address' as const,
    detailOrder: 1,
  },
  to: {
    section: 'Transfer',
    detailLabel: 'To',
    detailFormat: 'address' as const,
    detailOrder: 2,
  },
  amount: {
    section: 'Transfer',
    detailLabel: 'Amount',
    detailOrder: 3,
  },
  symbol: {
    section: 'Transfer',
    detailLabel: 'Symbol',
    detailOrder: 4,
  },
};

const transactionsMeta: Record<string, MetaOverlay> = {
  hash: {
    section: 'Transaction',
    detailLabel: 'Hash',
    detailFormat: 'hash' as const,
    detailOrder: 1,
  },
  from: {
    section: 'Transaction',
    detailLabel: 'From',
    detailFormat: 'address' as const,
    detailOrder: 2,
  },
  to: {
    section: 'Transaction',
    detailLabel: 'To',
    detailFormat: 'address' as const,
    detailOrder: 3,
  },
  value: {
    section: 'Transaction',
    detailLabel: 'Value',
    detailOrder: 4,
  },
  gasUsed: {
    section: 'Transaction',
    detailLabel: 'Gas Used',
    detailOrder: 5,
  },
  effectiveGasPrice: {
    section: 'Transaction',
    detailLabel: 'Gas Price',
    detailOrder: 6,
  },
};

const withdrawalsMeta: Record<string, MetaOverlay> = {
  address: {
    section: 'Withdrawal',
    detailLabel: 'Address',
    detailFormat: 'address' as const,
  },
  amount: {
    section: 'Withdrawal',
    detailLabel: 'Amount',
    detailOrder: 1,
  },
  ether: {
    section: 'Withdrawal',
    detailLabel: 'Ether',
    detailOrder: 2,
  },
};

const assetsMeta: Record<string, MetaOverlay> = {
  address: {
    section: 'Asset',
    detailLabel: 'Address',
    detailFormat: 'address' as const,
    detailOrder: 1,
  },
  name: {
    section: 'Asset',
    detailLabel: 'Name',
    detailOrder: 2,
  },
  symbol: {
    section: 'Asset',
    detailLabel: 'Symbol',
    detailOrder: 3,
  },
  decimals: {
    section: 'Details',
    detailLabel: 'Decimals',
    detailOrder: 1,
  },
  source: {
    section: 'Details',
    detailLabel: 'Source',
    detailOrder: 2,
  },
  tags: {
    section: 'Details',
    detailLabel: 'Tags',
    detailOrder: 3,
  },
};

const logsMeta: Record<string, MetaOverlay> = {
  date: {
    section: 'Event',
    detailLabel: 'Date',
    detailOrder: 1,
  },
  address: {
    section: 'Event',
    detailLabel: 'Address',
    detailFormat: 'address' as const,
    detailOrder: 2,
  },
  name: {
    section: 'Event',
    detailLabel: 'Name',
    detailOrder: 3,
  },
  articulatedLog: {
    section: 'Event',
    detailLabel: 'Articulated Log',
    detailOrder: 4,
  },
};

const tracesMeta: Record<string, MetaOverlay> = {
  date: {
    section: 'Trace',
    detailLabel: 'Date',
    detailOrder: 1,
  },
  type: {
    section: 'Trace',
    detailLabel: 'Type',
    detailOrder: 2,
  },
  compressedTrace: {
    section: 'Trace',
    detailLabel: 'Compressed Trace',
    detailOrder: 3,
  },
  error: {
    section: 'Trace',
    detailLabel: 'Error',
    detailOrder: 4,
  },
};

const receiptsMeta: Record<string, MetaOverlay> = {
  cumulativeGasUsed: {
    section: 'Gas',
    detailLabel: 'Cumulative Gas Used',
    detailOrder: 1,
  },
  effectiveGasPrice: {
    section: 'Gas',
    detailLabel: 'Effective Gas Price',
    detailOrder: 2,
  },
  from: {
    section: 'Addresses',
    detailLabel: 'From',
    detailFormat: 'address' as const,
    detailOrder: 1,
  },
  to: {
    section: 'Addresses',
    detailLabel: 'To',
    detailFormat: 'address' as const,
    detailOrder: 2,
  },
  blockHash: {
    section: 'Block Info',
    detailLabel: 'Block Hash',
    detailOrder: 1,
  },
};

const StatementsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.STATEMENTS),
  statementsMeta,
  { collapsedSections: [] },
);
const BalancesPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.BALANCES),
  balancesMeta,
  { collapsedSections: [] },
);
const TransfersPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.TRANSFERS),
  transfersMeta,
  { collapsedSections: [] },
);
const TransactionsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.TRANSACTIONS),
  transactionsMeta,
  { collapsedSections: [] },
);
const WithdrawalsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.WITHDRAWALS),
  withdrawalsMeta,
  { collapsedSections: [] },
);
const AssetsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.ASSETS),
  assetsMeta,
  { collapsedSections: ['Details'] },
);
const LogsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.LOGS),
  logsMeta,
  { collapsedSections: [] },
);
const TracesPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.TRACES),
  tracesMeta,
  { collapsedSections: [] },
);
const ReceiptsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.RECEIPTS),
  receiptsMeta,
  { collapsedSections: ['Gas', 'Block Info'] },
);

export const exportsDetailPanels = {
  [`${ROUTE}.${types.DataFacet.STATEMENTS}`]: StatementsPanel,
  [`${ROUTE}.${types.DataFacet.BALANCES}`]: BalancesPanel,
  [`${ROUTE}.${types.DataFacet.TRANSFERS}`]: TransfersPanel,
  [`${ROUTE}.${types.DataFacet.TRANSACTIONS}`]: TransactionsPanel,
  [`${ROUTE}.${types.DataFacet.WITHDRAWALS}`]: WithdrawalsPanel,
  [`${ROUTE}.${types.DataFacet.ASSETS}`]: AssetsPanel,
  [`${ROUTE}.${types.DataFacet.LOGS}`]: LogsPanel,
  [`${ROUTE}.${types.DataFacet.TRACES}`]: TracesPanel,
  [`${ROUTE}.${types.DataFacet.RECEIPTS}`]: ReceiptsPanel,
} as const;
