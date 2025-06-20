import { FormField } from '@components';
import { types } from '@models';

export const getColumns = (backendDataFacet: types.DataFacet): FormField[] => {
  switch (backendDataFacet) {
    case types.DataFacet.ALL:
    // fallthrough
    case types.DataFacet.CUSTOM:
    // fallthrough
    case types.DataFacet.PREFUND:
    // fallthrough
    case types.DataFacet.REGULAR:
    // fallthrough
    case types.DataFacet.BADDRESS:
    // fallthrough
    default:
      return getColumnsForName();
  }
};

const getColumnsForName = (): FormField[] => [
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
];
