import { FormField } from '@components';
import { types } from '@models';

// EXISTING_CODE
// EXISTING_CODE

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.ALL:
      return getNamesColumns();
    case types.DataFacet.CUSTOM:
      return getNamesColumns();
    case types.DataFacet.PREFUND:
      return getNamesColumns();
    case types.DataFacet.REGULAR:
      return getNamesColumns();
    case types.DataFacet.BADDRESS:
      return getNamesColumns();
    default:
      return [];
  }
};

const getNamesColumns = (): FormField[] => [
  // EXISTING_CODE
  {
    key: 'address',
    name: 'address',
    header: 'Address',
    label: 'Address',
    sortable: true,
    type: 'text',
    width: '340px',
    readOnly: true,
  },
  {
    key: 'name',
    name: 'name',
    header: 'Name',
    label: 'Name',
    sortable: true,
    type: 'text',
    width: '200px',
    required: true,
  },
  {
    key: 'tags',
    name: 'tags',
    header: 'Tags',
    label: 'Tags',
    sortable: true,
    type: 'text',
    width: '150px',
  },
  {
    key: 'source',
    name: 'source',
    header: 'Source',
    label: 'Source',
    sortable: true,
    type: 'text',
    width: '120px',
  },
  {
    key: 'symbol',
    name: 'symbol',
    header: 'Symbol',
    label: 'Symbol',
    sortable: true,
    type: 'text',
    width: '100px',
  },
  {
    key: 'decimals',
    name: 'decimals',
    header: 'Decimals',
    label: 'Decimals',
    sortable: true,
    type: 'number',
    width: '100px',
    textAlign: 'right',
  },
  {
    key: 'chips',
    name: 'chips',
    header: 'Chips',
    label: 'Chips',
    sortable: false,
    editable: false,
    width: '180px',
  },
  {
    key: 'actions',
    name: 'actions',
    header: 'Actions',
    label: 'Actions',
    sortable: false,
    editable: false,
    visible: true,
    type: 'button',
    width: '200px',
  },
  // EXISTING_CODE
];

// EXISTING_CODE
// EXISTING_CODE
