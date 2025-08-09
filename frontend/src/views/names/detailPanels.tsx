import { types } from '@models';

import {
  type MetaOverlay,
  buildDetailPanelFromColumns,
} from '../utils/detailPanel';
import { ROUTE, getColumns } from './columns';

const baseColumns = getColumns(types.DataFacet.ALL);

const meta: Record<string, MetaOverlay> = {
  address: { section: 'Identity', detailOrder: 1, detailFormat: 'address' },
  name: { section: 'Identity', detailOrder: 2 },
  symbol: { section: 'Token', detailOrder: 1 },
  decimals: { section: 'Token', detailOrder: 2 },
  tags: { section: 'Metadata', detailOrder: 1 },
  source: { section: 'Metadata', detailOrder: 2 },
};

const NamesDetailPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  baseColumns,
  meta,
  {
    collapsedSections: ['Metadata'],
  },
);

export const namesDetailPanels = {
  [`${ROUTE}.${types.DataFacet.ALL}`]: NamesDetailPanel,
} as const;
