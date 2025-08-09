import { types } from '@models';

import {
  type MetaOverlay,
  buildDetailPanelFromColumns,
} from '../utils/detailPanel';
import { ROUTE, getColumns } from './columns';

const monitorsMeta: Record<string, MetaOverlay> = {
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
  nRecords: {
    section: 'Statistics',
    detailLabel: 'Records',
    detailOrder: 1,
  },
  fileSize: {
    section: 'Statistics',
    detailLabel: 'File Size',
    detailOrder: 2,
  },
  isEmpty: {
    section: 'Statistics',
    detailLabel: 'Empty',
    detailOrder: 3,
  },
  lastScanned: {
    section: 'Scanning',
    detailLabel: 'Last Scanned',
    detailOrder: 1,
  },
  lastScanAge: {
    section: 'Scanning',
    detailLabel: 'Age',
    detailOnly: true,
    detailOrder: 2,
  },
};

const extras = (row: Record<string, unknown>) => {
  const ts = Number(row.lastScanned || 0);
  if (!Number.isFinite(ts) || ts === 0) return { lastScanAge: '-' };
  const ageSec = Math.max(0, Math.floor(Date.now() / 1000 - ts));
  const parts: string[] = [];
  const d = Math.floor(ageSec / 86400);
  const h = Math.floor((ageSec % 86400) / 3600);
  const m = Math.floor((ageSec % 3600) / 60);
  if (d) parts.push(`${d}d`);
  if (h) parts.push(`${h}h`);
  if (m && !d) parts.push(`${m}m`);
  if (!parts.length) parts.push(`${ageSec}s`);
  return { lastScanAge: parts.join(' ') };
};

const formatters = {
  lastScanned: (v: unknown) => {
    if (v === null || v === undefined) return '-';
    const num = Number(v);
    if (!Number.isFinite(num) || num === 0) return '-';
    return new Date(num * 1000)
      .toISOString()
      .replace('T', ' ')
      .substring(0, 19);
  },
  fileSize: (v: unknown) => {
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
  },
};

const MonitorsPanel = buildDetailPanelFromColumns<Record<string, unknown>>(
  getColumns(types.DataFacet.MONITORS),
  monitorsMeta,
  {
    promptWidthPx: 220,
    formatters,
    extras,
    collapsedSections: ['Stats'],
  },
);

export const monitorsDetailPanels = {
  [`${ROUTE}.${types.DataFacet.MONITORS}`]: MonitorsPanel,
} as const;
