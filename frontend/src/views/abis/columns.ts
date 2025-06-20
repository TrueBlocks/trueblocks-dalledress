import { FormField } from '@components';
import { types } from '@models';

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.FUNCTIONS:
    case types.DataFacet.EVENTS:
      return getColumnsForFunction();
    case types.DataFacet.KNOWN:
      return getColumnsForAbi().filter((col) => {
        const skip = col.key !== 'address';
        return skip;
      });
    case types.DataFacet.DOWNLOADED:
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
    textAlign: 'right',
  },
  {
    key: 'nFunctions',
    header: 'Functions',
    sortable: true,
    type: 'number',
    textAlign: 'right',
  },
  {
    key: 'nEvents',
    header: 'Events',
    sortable: true,
    type: 'number',
    textAlign: 'right',
  },
  {
    key: 'actions',
    header: 'Actions',
    sortable: false,
    type: 'text',
    width: '120px',
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
