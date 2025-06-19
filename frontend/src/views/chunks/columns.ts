import { FormField } from '@components';
import { types } from '@models';

export const getColumns = (listKind: types.ListKind): FormField[] => {
  switch (listKind) {
    case types.ListKind.STATS:
      return getColumnsForChunkStats();
    case types.ListKind.INDEX:
      return getColumnsForChunkIndex();
    case types.ListKind.BLOOMS:
      return getColumnsForChunkBloom();
    case types.ListKind.MANIFEST:
      return getColumnsForChunkManifest();
    default:
      return getColumnsForChunkStats(); // Default to stats
  }
};

// Column definitions for ChunkStats
const getColumnsForChunkStats = (): FormField[] => [
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'text',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'nAddrs',
    name: 'nAddrs',
    header: 'Addresses',
    label: 'Addresses',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nApps',
    name: 'nApps',
    header: 'Apps',
    label: 'Apps',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nBlocks',
    name: 'nBlocks',
    header: 'Blocks',
    label: 'Blocks',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nBloomsHit',
    name: 'nBloomsHit',
    header: 'Blooms Hit',
    label: 'Blooms Hit',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nBloomsMiss',
    name: 'nBloomsMiss',
    header: 'Blooms Miss',
    label: 'Blooms Miss',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'ratio',
    name: 'ratio',
    header: 'Hit Ratio',
    label: 'Hit Ratio',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
];

// Column definitions for ChunkIndex
const getColumnsForChunkIndex = (): FormField[] => [
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'text',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'magic',
    name: 'magic',
    header: 'Magic',
    label: 'Magic',
    sortable: true,
    type: 'text',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'hash',
    name: 'hash',
    header: 'Hash',
    label: 'Hash',
    sortable: true,
    type: 'text',
    width: '300px',
    readOnly: true,
  },
  {
    key: 'nAddresses',
    name: 'nAddresses',
    header: 'Addresses',
    label: 'Addresses',
    sortable: true,
    type: 'number',
    width: '120px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nAppearances',
    name: 'nAppearances',
    header: 'Appearances',
    label: 'Appearances',
    sortable: true,
    type: 'number',
    width: '130px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'size',
    name: 'size',
    header: 'Size',
    label: 'Size',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
];

// Column definitions for ChunkBloom
const getColumnsForChunkBloom = (): FormField[] => [
  {
    key: 'range',
    name: 'range',
    header: 'Range',
    label: 'Range',
    sortable: true,
    type: 'text',
    width: '150px',
    readOnly: true,
  },
  {
    key: 'magic',
    name: 'magic',
    header: 'Magic',
    label: 'Magic',
    sortable: true,
    type: 'text',
    width: '120px',
    readOnly: true,
  },
  {
    key: 'hash',
    name: 'hash',
    header: 'Hash',
    label: 'Hash',
    sortable: true,
    type: 'text',
    width: '300px',
    readOnly: true,
  },
  {
    key: 'nBlooms',
    name: 'nBlooms',
    header: 'Blooms',
    label: 'Blooms',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'nInserted',
    name: 'nInserted',
    header: 'Inserted',
    label: 'Inserted',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'size',
    name: 'size',
    header: 'Size',
    label: 'Size',
    sortable: true,
    type: 'number',
    width: '100px',
    readOnly: true,
    textAlign: 'right',
  },
  {
    key: 'byteWidth',
    name: 'byteWidth',
    header: 'Byte Width',
    label: 'Byte Width',
    sortable: true,
    type: 'number',
    width: '110px',
    readOnly: true,
    textAlign: 'right',
  },
];

// Column definitions for ChunkManifest (types.Manifest)
const getColumnsForChunkManifest = (): FormField[] => [
  {
    key: 'version',
    name: 'version',
    header: 'Version',
    label: 'Version',
    sortable: true,
    type: 'text',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'chain',
    name: 'chain',
    header: 'Chain',
    label: 'Chain',
    sortable: true,
    type: 'text',
    width: '100px',
    readOnly: true,
  },
  {
    key: 'specification',
    name: 'specification',
    header: 'Specification',
    label: 'Specification',
    sortable: true,
    type: 'text',
    width: '200px',
    readOnly: true,
  },
];
