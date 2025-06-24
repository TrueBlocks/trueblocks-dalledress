import { FormField } from '@components';
import { types } from '@models';

// EXISTING_CODE
// EXISTING_CODE

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.DOWNLOADED:
      return getAbisColumns();
    case types.DataFacet.KNOWN:
      return getAbisColumns();
    case types.DataFacet.FUNCTIONS:
      return getFunctionsColumns();
    case types.DataFacet.EVENTS:
      return getFunctionsColumns();
    default:
      return [];
  }
};

const getAbisColumns = (): FormField[] => [
  // EXISTING_CODE
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
  // EXISTING_CODE
];

const getFunctionsColumns = (): FormField[] => [
  // EXISTING_CODE
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
  // EXISTING_CODE
];

// EXISTING_CODE
// EXISTING_CODE
