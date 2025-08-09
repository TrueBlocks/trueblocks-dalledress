import { types } from '@models';

import {
  type MetaOverlay,
  buildDetailPanelFromColumns,
} from '../utils/detailPanel';
import { ROUTE, getColumns } from './columns';

const statsMeta: Record<string, MetaOverlay> = {
  range: { section: 'Range', detailLabel: 'Range', detailOrder: 1 },
  nAddrs: { section: 'Counts', detailLabel: 'Addresses', detailOrder: 1 },
  nApps: { section: 'Counts', detailLabel: 'Apps', detailOrder: 2 },
  nBlocks: { section: 'Counts', detailLabel: 'Blocks', detailOrder: 3 },
  nBlooms: { section: 'Counts', detailLabel: 'Blooms', detailOrder: 4 },
  recWid: { section: 'Sizes', detailLabel: 'Record Width', detailOrder: 1 },
  bloomSz: { section: 'Sizes', detailLabel: 'Bloom Size', detailOrder: 2 },
  chunkSz: { section: 'Sizes', detailLabel: 'Chunk Size', detailOrder: 3 },
  addrsPerBlock: {
    section: 'Efficiency',
    detailLabel: 'Addrs/Block',
    detailOrder: 1,
  },
  appsPerBlock: {
    section: 'Efficiency',
    detailLabel: 'Apps/Block',
    detailOrder: 2,
  },
  appsPerAddr: {
    section: 'Efficiency',
    detailLabel: 'Apps/Addr',
    detailOrder: 3,
  },
  ratio: { section: 'Efficiency', detailLabel: 'Ratio', detailOrder: 4 },
};

const indexMeta: Record<string, MetaOverlay> = {
  range: { section: 'Range', detailLabel: 'Range', detailOrder: 1 },
  magic: { section: 'Identity', detailLabel: 'Magic', detailOrder: 1 },
  hash: {
    section: 'Identity',
    detailLabel: 'Hash',
    detailOrder: 2,
    detailFormat: 'hash',
  },
  nAddresses: { section: 'Counts', detailLabel: 'Addresses', detailOrder: 1 },
  nAppearances: {
    section: 'Counts',
    detailLabel: 'Appearances',
    detailOrder: 2,
  },
  size: { section: 'Sizes', detailLabel: 'Size', detailOrder: 1 },
};

const bloomsMeta: Record<string, MetaOverlay> = {
  range: { section: 'Range', detailLabel: 'Range', detailOrder: 1 },
  magic: { section: 'Identity', detailLabel: 'Magic', detailOrder: 1 },
  hash: {
    section: 'Identity',
    detailLabel: 'Hash',
    detailOrder: 2,
    detailFormat: 'hash',
  },
  nBlooms: { section: 'Counts', detailLabel: 'Blooms', detailOrder: 1 },
  nInserted: { section: 'Counts', detailLabel: 'Inserted', detailOrder: 2 },
  size: { section: 'Sizes', detailLabel: 'Size', detailOrder: 1 },
  byteWidth: { section: 'Sizes', detailLabel: 'Byte Width', detailOrder: 2 },
};

const manifestMeta: Record<string, MetaOverlay> = {
  version: { section: 'Manifest', detailLabel: 'Version', detailOrder: 1 },
  chain: { section: 'Manifest', detailLabel: 'Chain', detailOrder: 2 },
  specification: {
    section: 'Manifest',
    detailLabel: 'Specification',
    detailOrder: 3,
    detailFormat: 'hash',
  },
};

const StatsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.STATS),
  statsMeta,
  { collapsedSections: ['Sizes', 'Efficiency'] },
);
const IndexPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.INDEX),
  indexMeta,
  { collapsedSections: ['Identity'] },
);
const BloomsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.BLOOMS),
  bloomsMeta,
  { collapsedSections: ['Identity', 'Sizes'] },
);
const ManifestPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.MANIFEST),
  manifestMeta,
  { collapsedSections: [] },
);

export const chunksDetailPanels: Record<
  string,
  (rowData: Record<string, unknown> | null) => React.ReactNode
> = {
  [`${ROUTE}.${types.DataFacet.STATS}`]: StatsPanel,
  [`${ROUTE}.${types.DataFacet.INDEX}`]: IndexPanel,
  [`${ROUTE}.${types.DataFacet.BLOOMS}`]: BloomsPanel,
  [`${ROUTE}.${types.DataFacet.MANIFEST}`]: ManifestPanel,
};
