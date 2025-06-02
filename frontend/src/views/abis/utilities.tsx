import { ReactNode } from 'react';

import { FormField } from '@components';
import { Action } from '@components';
import { Group, Text } from '@mantine/core';
import { types } from '@models';

import { AbiRow, TableConfigProps } from './types';

//--------------------------------------------------------------------
export const getTableConfig = (
  listKind: types.ListKind,
  config: TableConfigProps,
): {
  data: AbiRow[];
  columns: FormField<AbiRow>[];
} => {
  const abiColumns = createAbiColumns(
    listKind,
    listKind === types.ListKind.DOWNLOADED
      ? config.isDownloadedLoaded
      : listKind === types.ListKind.KNOWN
        ? config.isKnownLoaded
        : listKind === types.ListKind.FUNCTIONS
          ? config.isFuncsLoaded
          : config.isEventsLoaded,
    config.processingAddresses,
    config.setSelectedAddress,
    config.setLocation,
    config.handleAction,
  ) as FormField<AbiRow>[];

  const funcColumns = createFunctionColumns() as FormField<AbiRow>[];

  switch (listKind) {
    case types.ListKind.FUNCTIONS:
      return {
        data: config.functions as AbiRow[],
        columns: funcColumns,
      };
    case types.ListKind.EVENTS:
      return {
        data: config.events as AbiRow[],
        columns: funcColumns,
      };
    case types.ListKind.KNOWN:
      return {
        data: config.known as AbiRow[],
        columns: abiColumns,
      };
    case types.ListKind.DOWNLOADED:
    // fall through
    default:
      return {
        data: config.downloaded as AbiRow[],
        columns: abiColumns,
      };
  }
};

//--------------------------------------------------------------------
const createAbiColumns = (
  listKind: types.ListKind,
  collectionIsLoaded: boolean,
  processingAddresses: Set<string>,
  setSelectedAddress: (address: string) => void,
  setLocation: (path: string) => void,
  handleAction: (address: string) => void,
): FormField<types.Abi>[] => {
  const baseColumns: FormField<types.Abi>[] = [
    {
      key: 'address',
      header: 'Address',
      sortable: true,
      width: 'col-address',
      render: (row: types.Abi) => (
        <Text size="sm" style={{ fontFamily: 'monospace' }}>
          {row.address.toString()}
        </Text>
      ),
    },
    {
      key: 'name',
      header: 'Name',
      sortable: true,
    },
    {
      key: 'nFunctions',
      header: 'Functions',
      sortable: true,
      width: 'col-base-md',
      render: (row: types.Abi) => <Text ta="right">{row.nFunctions}</Text>,
    },
    {
      key: 'nEvents',
      header: 'Events',
      sortable: true,
      width: 'col-base-md',
      render: (row: types.Abi) => <Text ta="right">{row.nEvents}</Text>,
    },
    {
      key: 'fileSize',
      header: 'Size (bytes)',
      sortable: true,
      width: 'col-base-md',
      render: (row: types.Abi) => <Text ta="right">{row.fileSize}</Text>,
    },
    {
      key: 'lastModDate',
      header: 'Last Modified',
      sortable: true,
      width: 'col-date',
    },
  ];

  if (listKind === types.ListKind.KNOWN) {
    return baseColumns.slice(1);
  } else {
    return [
      ...baseColumns,
      {
        key: 'actions',
        header: 'Actions',
        width: 'col-actions',
        render: (row: types.Abi) => {
          return renderActions(
            row as AbiRow,
            collectionIsLoaded,
            processingAddresses,
            setSelectedAddress,
            setLocation,
            handleAction,
          );
        },
      },
    ];
  }
};

//--------------------------------------------------------------------
const createFunctionColumns = (): FormField<types.Function>[] => {
  return [
    {
      key: 'encoding',
      header: 'Encoding',
      sortable: true,
      width: 'col-encoding',
    },
    {
      key: 'name',
      header: 'Name',
      sortable: true,
    },
    {
      key: 'signature',
      header: 'Signature',
      sortable: true,
    },
  ];
};

//--------------------------------------------------------------------
const renderActions = (
  item: types.Abi | types.Function,
  collectionIsLoaded: boolean,
  processingAddresses: Set<string>,
  setSelectedAddress: (address: string) => void,
  setLocation: (path: string) => void,
  handleAction: (address: string) => void,
): ReactNode => {
  const addressStr = 'address' in item ? String(item.address) : '';
  const isProcessing = processingAddresses.has(addressStr);
  return (
    <Group gap="xs">
      <Action
        icon="History"
        onClick={() => {
          setSelectedAddress(addressStr);
          setLocation(`/history/${addressStr}`);
        }}
        disabled={!collectionIsLoaded || isProcessing}
        title="View History"
        size="sm"
      />
      <Action
        icon={'Delete'}
        onClick={() => handleAction(addressStr)}
        disabled={!collectionIsLoaded || isProcessing}
        title={'Delete'}
        size="sm"
      />
    </Group>
  );
};
