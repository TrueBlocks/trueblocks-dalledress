import { types } from '@models';

import {
  type MetaOverlay,
  buildDetailPanelFromColumns,
} from '../utils/detailPanel';
import { ROUTE, getColumns } from './columns';

const downloadedKnownMeta: Record<string, MetaOverlay> = {
  name: {
    section: 'Overview',
    detailLabel: 'Name',
    detailOrder: 1,
  },
  address: {
    section: 'Overview',
    detailLabel: 'Address',
    detailFormat: 'address' as const,
    detailOrder: 2,
  },
  fileSize: {
    section: 'Statistics',
    detailLabel: 'File Size',
    detailOrder: 1,
  },
  nFunctions: {
    section: 'Statistics',
    detailLabel: 'Functions',
    detailOrder: 2,
  },
  nEvents: {
    section: 'Statistics',
    detailLabel: 'Events',
    detailOrder: 3,
  },
};

const functionsEventsMeta: Record<string, MetaOverlay> = {
  name: {
    section: 'Definition',
    detailLabel: 'Name',
    detailOrder: 1,
  },
  type: {
    section: 'Definition',
    detailLabel: 'Type',
    detailOrder: 2,
  },
  encoding: {
    section: 'Encoding',
    detailLabel: 'Encoding',
    detailOrder: 1,
  },
  signature: {
    section: 'Encoding',
    detailLabel: 'Signature',
    detailOrder: 2,
  },
  paramCount: {
    section: 'Inputs',
    detailLabel: 'Param Count',
    detailOnly: true,
    detailOrder: 1,
  },
  params: {
    section: 'Inputs',
    detailLabel: 'Parameters',
    detailOnly: true,
    detailOrder: 2,
  },
};

const fileSizeFormatter = (v: unknown) => {
  if (v === null || v === undefined) return '-';
  const num = Number(v);
  if (!Number.isFinite(num)) return String(v);
  if (num < 1024) return `${num} B`;
  const units = ['KB', 'MB', 'GB', 'TB'];
  let val = num;
  let u = -1;
  while (val >= 1024 && u < units.length - 1) {
    val /= 1024;
    u++;
  }
  return `${val.toFixed(1)} ${units[u]}`;
};

const functionsEventsExtras = (row: Record<string, unknown>) => {
  const signature = String(row.signature || '');
  const open = signature.indexOf('(');
  const close = signature.lastIndexOf(')');
  let inner = '';
  if (open !== -1 && close !== -1 && close > open)
    inner = signature.slice(open + 1, close);
  const params: string[] = [];
  if (inner.trim().length) {
    let depth = 0;
    let token = '';
    for (let i = 0; i < inner.length; i++) {
      const ch = inner[i];
      if (ch === '(') {
        depth++;
        token += ch;
      } else if (ch === ')') {
        depth--;
        token += ch;
      } else if (ch === ',' && depth === 0) {
        if (token.trim()) params.push(token.trim());
        token = '';
      } else token += ch;
    }
    if (token.trim()) params.push(token.trim());
  }
  const paramCount = params.length;
  const paramsJoined = params.join('\n');
  return { paramCount, params: paramsJoined };
};

const DownloadedPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.DOWNLOADED),
  downloadedKnownMeta,
  { collapsedSections: ['Stats'], formatters: { fileSize: fileSizeFormatter } },
);
const KnownPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.KNOWN),
  downloadedKnownMeta,
  { collapsedSections: ['Stats'], formatters: { fileSize: fileSizeFormatter } },
);
const FunctionsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.FUNCTIONS),
  functionsEventsMeta,
  { extras: functionsEventsExtras },
);
const EventsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.EVENTS),
  functionsEventsMeta,
  { extras: functionsEventsExtras },
);

export const abisDetailPanels: Record<
  string,
  (rowData: Record<string, unknown> | null) => React.ReactNode
> = {
  [`${ROUTE}.${types.DataFacet.DOWNLOADED}`]: DownloadedPanel,
  [`${ROUTE}.${types.DataFacet.KNOWN}`]: KnownPanel,
  [`${ROUTE}.${types.DataFacet.FUNCTIONS}`]: FunctionsPanel,
  [`${ROUTE}.${types.DataFacet.EVENTS}`]: EventsPanel,
};
