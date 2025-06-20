import { FormField } from '@components';
import { types } from '@models';

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.MONITORS:
    default:
      return getColumnsForMonitor();
  }
};

const getColumnsForMonitor = (): FormField[] => [
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
  },
  {
    key: 'nRecords',
    name: 'nRecords',
    header: 'Records',
    label: 'Records',
    sortable: true,
    type: 'number',
    width: '100px',
    textAlign: 'right',
  },
  {
    key: 'fileSize',
    name: 'fileSize',
    header: 'File Size',
    label: 'File Size',
    sortable: true,
    type: 'number',
    width: '120px',
    textAlign: 'right',
  },
  {
    key: 'isEmpty',
    name: 'isEmpty',
    header: 'Empty',
    label: 'Empty',
    sortable: true,
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'lastScanned',
    name: 'lastScanned',
    header: 'Last Scanned',
    label: 'Last Scanned',
    sortable: true,
    type: 'timestamp',
    width: '140px',
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
    width: '120px',
  },
];
