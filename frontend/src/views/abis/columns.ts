// ABIS_ROUTE
import { FormField } from '@components';
import { types } from '@models';

export const getColumns = (listKind: types.ListKind): FormField[] => {
  switch (listKind) {
    case 'Functions':
    case 'Events':
      return getColumnsForFunction();
    case 'Known':
      return getColumnsForAbi().filter((col) => {
        const skip = col.key !== 'address';
        return skip;
      });
    case 'Downloaded':
    // fallthrough intended
    default:
      return getColumnsForAbi();
  }
};

const getColumnsForAbi = (): FormField[] => [
  {
    key: 'address',
    header: 'Address',
    sortable: true,
    type: 'text',
  },
  {
    key: 'name',
    header: 'Name',
    sortable: true,
    type: 'text',
  },
  {
    key: 'fileSize',
    header: 'File Size',
    sortable: true,
    type: 'number',
  },
  {
    key: 'nFunctions',
    header: 'Functions',
    sortable: true,
    type: 'number',
  },
  {
    key: 'nEvents',
    header: 'Events',
    sortable: true,
    type: 'number',
  },
  {
    key: 'nErrors',
    header: 'Errors',
    sortable: true,
    type: 'number',
  },
];

const getColumnsForFunction = (): FormField[] => [
  {
    key: 'encoding',
    header: 'Encoding',
    type: 'text',
    sortable: true,
    width: 'col-encoding',
  },
  {
    key: 'name',
    header: 'Name',
    type: 'text',
    sortable: true,
    width: 'col-name',
  },
  {
    key: 'type',
    header: 'Type',
    type: 'text',
    sortable: true,
    width: 'col-type',
  },
  {
    key: 'signature',
    header: 'Signature',
    type: 'text',
    sortable: true,
    width: 'col-signature',
  },
];

// ABIS_ROUTE
