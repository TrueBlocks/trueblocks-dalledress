// ADD_ROUTE
import { useCallback, useEffect, useMemo, useState } from 'react';

import { DeleteAbi, GetAbisPage, Reload } from '@app';
import {
  Action,
  FormField,
  Table,
  TableProvider,
  usePagination,
} from '@components';
import { TableKey, useAppContext } from '@contexts';
import { TabView } from '@layout';
import { Group, Text } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';
import { msgs, sorting, types } from '@models';
import { EventsOn } from '@runtime';
import { Log, useEmitters } from '@utils';
import { useLocation } from 'wouter';

import './Abis.css';
import { AbiRow, IndexedAbi, IndexedFunction, TableConfigProps } from './Types';

//--------------------------------------------------------------------
export const Abis = () => {
  const { emitStatus, emitError } = useEmitters();
  const { lastTab, setSelectedAddress } = useAppContext();
  const [, setLocation] = useLocation();

  const [listKind, setListKind] = useState<types.ListKind>(
    (lastTab['/abis'] as types.ListKind) || types.ListKind.DOWNLOADED,
  );
  const [downloadedAbis, setDownloadedAbis] = useState<IndexedAbi[]>([]);
  const [knownAbis, setKnownAbis] = useState<IndexedAbi[]>([]);
  const [functions, setFunctions] = useState<IndexedFunction[]>([]);
  const [events, setEvents] = useState<IndexedFunction[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  const [currentSort, setCurrentSort] = useState<sorting.SortDef | null>(null);
  const [currentFilter, setCurrentFilter] = useState<string>('');
  const [processingAddresses, setProcessingAddresses] = useState<Set<string>>(
    new Set(),
  );
  const [isDownloadedLoaded, setIsDownloadedLoaded] = useState<boolean>(false);
  const [isKnownLoaded, setIsKnownLoaded] = useState<boolean>(false);
  const [isFuncsLoaded, setIsFuncsLoaded] = useState<boolean>(false);
  const [isEventsLoaded, setIsEventsLoaded] = useState<boolean>(false);

  const tableKey = useMemo((): TableKey => {
    return { viewName: '/abis', tabName: listKind };
  }, [listKind]);

  const { pagination, setTotalItems } = usePagination(tableKey);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const sortArgument = currentSort === null ? undefined : currentSort;
      const result = await GetAbisPage(
        listKind,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sortArgument as sorting.SortDef,
        currentFilter,
      );

      switch (listKind) {
        case types.ListKind.DOWNLOADED:
          setIsDownloadedLoaded(result.isLoaded);
          break;
        case types.ListKind.KNOWN:
          setIsKnownLoaded(result.isLoaded);
          break;
        case types.ListKind.FUNCTIONS:
          setIsFuncsLoaded(result.isLoaded);
          break;
        case types.ListKind.EVENTS:
          setIsEventsLoaded(result.isLoaded);
          break;
      }

      switch (listKind) {
        case types.ListKind.DOWNLOADED:
          setDownloadedAbis((result.abis as IndexedAbi[]) || []);
          break;
        case types.ListKind.KNOWN:
          setKnownAbis((result.abis as IndexedAbi[]) || []);
          break;
        case types.ListKind.FUNCTIONS:
          setFunctions((result.functions as IndexedFunction[]) || []);
          break;
        case types.ListKind.EVENTS:
          setEvents((result.functions as IndexedFunction[]) || []);
          break;
        default:
          throw new Error(`Unknown list kind: ${listKind}`);
      }
      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      const e = err as Error;
      setError(e);
      emitError(`${e.message} Failed to fetch ${listKind}`);
      Log(`Error fetching ${listKind}: ${e}`);
    } finally {
      setLoading(false);
    }
  }, [
    currentSort,
    listKind,
    pagination.currentPage,
    pagination.pageSize,
    currentFilter,
    setTotalItems,
    emitError,
  ]);

  useEffect(() => {
    fetchData();
  }, [fetchData, listKind]);

  useEffect(() => {
    const eventName = msgs.EventType.DATA_LOADED;
    const unlisten = EventsOn(eventName, (payload: unknown) => {
      const eventPayload = payload as types.DataLoadedPayload;

      if (eventPayload) {
        switch (listKind) {
          case types.ListKind.DOWNLOADED:
            setIsDownloadedLoaded(eventPayload.isLoaded);
            break;
          case types.ListKind.KNOWN:
            setIsKnownLoaded(eventPayload.isLoaded);
            break;
          case types.ListKind.FUNCTIONS:
            setIsFuncsLoaded(eventPayload.isLoaded);
            break;
          case types.ListKind.EVENTS:
            setIsEventsLoaded(eventPayload.isLoaded);
            break;
        }

        if (
          eventPayload.dataType === 'functions-events' ||
          !eventPayload.isLoaded
        ) {
          fetchData();
        }
      }
    });

    return () => {
      if (typeof unlisten === 'function') {
        unlisten();
      }
    };
  }, [fetchData, listKind]);

  useEffect(() => {
    const currentTabLabel = lastTab['/abis'] as types.ListKind | undefined;
    if (currentTabLabel && currentTabLabel !== listKind) {
      setListKind(currentTabLabel);
    }
  }, [lastTab, listKind]);

  useHotkeys([
    [
      'mod+r',
      () => {
        setLoading(true);
        setIsDownloadedLoaded(false);
        setIsKnownLoaded(false);
        setIsFuncsLoaded(false);
        setIsEventsLoaded(false);
        Reload().then(() => {
          fetchData();
          emitStatus('Reloaded ABI data. Fetching fresh data...');
        });
      },
    ],
  ]);

  const handleAction = (address: string, _isDeleted: boolean) => {
    const addressStr = address;
    setProcessingAddresses((prev) => new Set(prev).add(addressStr));
    try {
      const original = [...downloadedAbis];
      const optimisticNames = original.filter((abi) => {
        const nameAddress =
          typeof abi.address === 'string' ? abi.address : String(abi.address);
        return nameAddress !== addressStr;
      });
      setDownloadedAbis(optimisticNames as IndexedAbi[]);

      DeleteAbi(addressStr)
        .then(async () => {
          // If API call is successful, refresh the data to get the definitive state
          const sortArgument = currentSort === null ? undefined : currentSort;
          const result = await GetAbisPage(
            listKind,
            pagination.currentPage * pagination.pageSize,
            pagination.pageSize,
            sortArgument as sorting.SortDef,
            currentFilter,
          );

          setDownloadedAbis((result.abis || []) as IndexedAbi[]);
          setTotalItems(result.expectedTotal || 0);

          const action = 'deleted';
          emitStatus(`Address ${addressStr} was ${action} successfully`);
        })
        .catch((err) => {
          // If there's an error, revert the optimistic update
          setDownloadedAbis(optimisticNames as IndexedAbi[]);
          emitError(err);
          Log(`Error in handleAction: ${err}`);
        });
    } finally {
      // Always clean up the processing state
      setProcessingAddresses((prev) => {
        const newSet = new Set(prev);
        newSet.delete(addressStr);
        return newSet;
      });
    }
  };

  const handleTableSubmit = (formData: Record<string, unknown>) => {
    Log(`Table submitted: ${formData}`);
  };

  const renderTable = (listKind: types.ListKind) => {
    const config: TableConfigProps = {
      downloadedAbis,
      knownAbis,
      functions,
      events,
      isDownloadedLoaded,
      isKnownLoaded,
      isFuncsLoaded,
      isEventsLoaded,
      processingAddresses,
      setSelectedAddress,
      setLocation,
      handleAction,
    };

    const { data, columns } = getTableConfig(listKind, config);

    // Determine if the current collection is still loading
    const isCurrentCollectionLoading = (() => {
      switch (listKind) {
        case types.ListKind.DOWNLOADED:
          return !isDownloadedLoaded;
        case types.ListKind.KNOWN:
          return !isKnownLoaded;
        case types.ListKind.FUNCTIONS:
          return !isFuncsLoaded;
        case types.ListKind.EVENTS:
          return !isEventsLoaded;
        default:
          return loading;
      }
    })();

    return (
      <TableProvider>
        <div className="tableContainer">
          <Table<AbiRow>
            columns={columns}
            data={data}
            loading={loading || isCurrentCollectionLoading}
            sort={currentSort}
            onSortChange={setCurrentSort}
            filter={currentFilter}
            onFilterChange={setCurrentFilter}
            tableKey={tableKey}
            onSubmit={handleTableSubmit}
          />
        </div>
      </TableProvider>
    );
  };

  const createOneTab = (listKind: types.ListKind) => {
    return {
      label: listKind,
      value: listKind,
      content: renderTable(listKind),
    };
  };

  const tabs = [
    createOneTab(types.ListKind.DOWNLOADED),
    createOneTab(types.ListKind.KNOWN),
    createOneTab(types.ListKind.FUNCTIONS),
    createOneTab(types.ListKind.EVENTS),
  ];

  return (
    <div className="abisView">
      <TabView tabs={tabs} route="/abis" />
      {error && (
        <div className="error-message-placeholder">
          <h3>{`Error fetching ${listKind}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
    </div>
  );
};

//--------------------------------------------------------------------
const getTableConfig = (
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
        data: config.knownAbis as AbiRow[],
        columns: abiColumns,
      };
    case types.ListKind.DOWNLOADED:
    // fall through
    default:
      return {
        data: config.downloadedAbis as AbiRow[],
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
  handleAction: (address: string, isDeleted: boolean) => void,
): FormField<types.Abi>[] => {
  const baseColumns: FormField<types.Abi>[] = [
    {
      key: 'address',
      header: 'Address',
      sortable: true,
      width: 'col-address',
      render: (row: types.Abi) => (
        <Text size="sm" ff="monospace">
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
            row,
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
  handleAction: (address: string, isDeleted: boolean) => void,
): React.ReactNode => {
  const addressStr = 'address' in item ? item.address.toString() : '';
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
        onClick={() => handleAction(addressStr, false)}
        disabled={!collectionIsLoaded || isProcessing}
        title={'Delete'}
        size="sm"
      />
    </Group>
  );
};

// ADD_ROUTE
