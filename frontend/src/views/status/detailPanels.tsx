import { types } from '@models';

import {
  type MetaOverlay,
  buildDetailPanelFromColumns,
} from '../utils/detailPanel';
import { ROUTE, getColumns } from './columns';

const statusMeta: Record<string, MetaOverlay> = {
  cachePath: { section: 'Paths', detailLabel: 'Cache Path', detailOrder: 1 },
  indexPath: { section: 'Paths', detailLabel: 'Index Path', detailOrder: 2 },
  chain: { section: 'Chain', detailLabel: 'Chain', detailOrder: 1 },
  chainId: { section: 'Chain', detailLabel: 'Chain ID', detailOrder: 2 },
  networkId: { section: 'Chain', detailLabel: 'Network ID', detailOrder: 3 },
  chainConfig: {
    section: 'Chain',
    detailLabel: 'Chain Config',
    detailOrder: 4,
  },
  rootConfig: { section: 'Config', detailLabel: 'Root Config', detailOrder: 1 },
  clientVersion: {
    section: 'Config',
    detailLabel: 'Client Version',
    detailOrder: 2,
  },
  version: { section: 'Config', detailLabel: 'Version', detailOrder: 3 },
  progress: { section: 'Progress', detailLabel: 'Progress', detailOrder: 1 },
  rpcProvider: {
    section: 'Providers',
    detailLabel: 'RPC Provider',
    detailOrder: 1,
  },
  hasEsKey: { section: 'Flags', detailLabel: 'Has ES Key', detailOrder: 1 },
  hasPinKey: { section: 'Flags', detailLabel: 'Has Pin Key', detailOrder: 2 },
  isApi: { section: 'Flags', detailLabel: 'Is API', detailOrder: 3 },
  isArchive: { section: 'Flags', detailLabel: 'Is Archive', detailOrder: 4 },
  isScraping: { section: 'Flags', detailLabel: 'Is Scraping', detailOrder: 5 },
  isTesting: { section: 'Flags', detailLabel: 'Is Testing', detailOrder: 6 },
  isTracing: { section: 'Flags', detailLabel: 'Is Tracing', detailOrder: 7 },
};

const cacheMeta: Record<string, MetaOverlay> = {
  type: {
    section: 'General',
    detailLabel: 'Type',
    detailOrder: 1,
  },
  path: {
    section: 'General',
    detailLabel: 'Path',
    detailOrder: 2,
  },
  nFiles: {
    section: 'Statistics',
    detailLabel: 'Files',
    detailOrder: 1,
  },
  nFolders: {
    section: 'Statistics',
    detailLabel: 'Folders',
    detailOrder: 2,
  },
  sizeInBytes: {
    section: 'Statistics',
    detailLabel: 'Size (Bytes)',
    detailOrder: 3,
  },
  lastCached: {
    section: 'Timestamps',
    detailLabel: 'Last Cached',
    detailOrder: 1,
  },
};

const chainMeta: Record<string, MetaOverlay> = {
  chain: {
    section: 'General',
    detailLabel: 'Chain',
    detailOrder: 1,
  },
  chainId: {
    section: 'General',
    detailLabel: 'Chain ID',
    detailOrder: 2,
  },
  symbol: {
    section: 'General',
    detailLabel: 'Symbol',
    detailOrder: 3,
  },
  rpcProvider: {
    section: 'Providers',
    detailLabel: 'RPC Provider',
    detailOrder: 1,
  },
  ipfsGateway: {
    section: 'Providers',
    detailLabel: 'IPFS Gateway',
    detailOrder: 2,
  },
  localExplorer: {
    section: 'Explorers',
    detailLabel: 'Local Explorer',
    detailOrder: 1,
  },
  remoteExplorer: {
    section: 'Explorers',
    detailLabel: 'Remote Explorer',
    detailOrder: 2,
  },
};

const StatusPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.STATUS),
  statusMeta,
  { collapsedSections: ['Flags'] },
);

const CachePanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.CACHES),
  cacheMeta,
  { collapsedSections: ['Statistics', 'Timestamps'] },
);

const ChainPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.CHAINS),
  chainMeta,
  { collapsedSections: ['Providers', 'Explorers'] },
);

export const statusDetailPanels = {
  [`${ROUTE}.${types.DataFacet.STATUS}`]: StatusPanel,
  [`${ROUTE}.${types.DataFacet.CACHES}`]: CachePanel,
  [`${ROUTE}.${types.DataFacet.CHAINS}`]: ChainPanel,
} as const;
