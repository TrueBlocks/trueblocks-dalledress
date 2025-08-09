import { types } from '@models';

import {
  type MetaOverlay,
  buildDetailPanelFromColumns,
} from '../utils/detailPanel';
import { ROUTE, getColumns } from './columns';

const eventsMeta: Record<string, MetaOverlay> = {
  date: { section: 'Block/Tx', detailOrder: 1 },
  blockNumber: { section: 'Block/Tx', detailOrder: 2 },
  transactionIndex: { section: 'Block/Tx', detailOrder: 3 },
  transactionHash: {
    section: 'Block/Tx',
    detailOrder: 4,
    detailFormat: 'hash',
  },
  blockHash: { section: 'Block/Tx', detailOrder: 5, detailFormat: 'hash' },
  address: { section: 'Event', detailOrder: 1, detailFormat: 'address' },
  name: { section: 'Event', detailOrder: 2 },
  articulatedLog: { section: 'Event', detailOrder: 3, detailFormat: 'json' },
  topics: { section: 'Event', detailOrder: 4, detailFormat: 'json' },
  data: { section: 'Event', detailOrder: 5, detailFormat: 'json' },
};

const EventsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.EVENTS),
  eventsMeta,
  { collapsedSections: ['Event'] },
);

export const contractsDetailPanels: Record<
  string,
  (rowData: Record<string, unknown> | null) => React.ReactNode
> = {
  [`${ROUTE}.${types.DataFacet.EVENTS}`]: EventsPanel,
};
