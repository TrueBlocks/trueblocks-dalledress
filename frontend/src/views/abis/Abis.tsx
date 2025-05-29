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

export type AbisListType = 'Downloaded' | 'Known' | 'Functions' | 'Events';

type IndexedAbi = types.Abi & {
  [key: string]: unknown;
};

type IndexedFunction = types.Function & {
  [key: string]: unknown;
};

export const Abis = () => {
  // Example usage of composable CSS column sizing:
  // Basic: width: 'col-address'           -> 340px fixed width
  // Basic: width: 'col-min'              -> min-content
  // Composed: width: 'col-base-address col-min'    -> Start with 340px, but shrink to min-content
  // Composed: width: 'col-base-sm col-fit'         -> Start with 100px, but fit to content
  // Composed: width: 'col-base-lg col-expand'      -> Start with 160px, but expand to fill

  const { emitStatus, emitError } = useEmitters();
  const { lastTab, setSelectedAddress } = useAppContext();
  const [, setLocation] = useLocation();

  const [listType, setListType] = useState<AbisListType>(
    (lastTab['/abis'] as AbisListType) || 'Downloaded',
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

  const [collectionIsLoading, setCollectionIsLoading] =
    useState<boolean>(false);
  const [collectionIsLoaded, setCollectionIsLoaded] = useState<boolean>(false);

  const tableKey = useMemo((): TableKey => {
    return { viewName: '/abis', tabName: listType };
  }, [listType]);

  const { pagination, setTotalItems } = usePagination(tableKey);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const sortArgument = currentSort === null ? undefined : currentSort;
      const result = await GetAbisPage(
        listType,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sortArgument as sorting.SortDef,
        currentFilter,
      );

      setCollectionIsLoading(result.IsLoading);
      setCollectionIsLoaded(result.IsLoaded);

      if (listType === 'Downloaded') {
        setDownloadedAbis((result.Abis as IndexedAbi[]) || []);
      } else if (listType === 'Known') {
        setKnownAbis((result.Abis as IndexedAbi[]) || []);
      } else if (listType === 'Functions') {
        setFunctions((result.Functions as IndexedFunction[]) || []);
      } else if (listType === 'Events') {
        setEvents((result.Functions as IndexedFunction[]) || []);
      }
      setTotalItems(result.TotalItems || 0);
      // emitStatus(
      //   `Fetched ${result.Abis?.length || result.Functions?.length} ${listType} items for current page.`,
      // );
    } catch (err: unknown) {
      const e = err as Error;
      setError(e);
      emitError(`${e.message} Failed to fetch ${listType}`);
      Log(`Error fetching ${listType}: ${e}`);
    } finally {
      setLoading(false);
    }
  }, [
    currentSort,
    listType,
    pagination.currentPage,
    pagination.pageSize,
    currentFilter,
    // emitStatus,
    setTotalItems,
    emitError,
  ]);

  useEffect(() => {
    fetchData();
  }, [fetchData, listType]);

  useEffect(() => {
    const eventName = msgs.EventType.DATA_LOADED;
    const unlisten = EventsOn(eventName, (payload: unknown) => {
      const eventPayload = payload as types.DataLoadedPayload;

      if (eventPayload) {
        setCollectionIsLoaded(eventPayload.isLoaded);
        setCollectionIsLoading(!eventPayload.isLoaded);

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
  }, [fetchData]);

  useEffect(() => {
    const currentTabLabel = lastTab['/abis'] as AbisListType | undefined;
    if (currentTabLabel && currentTabLabel !== listType) {
      setListType(currentTabLabel);
    }
  }, [lastTab, listType]);

  const handleRefresh = () => {
    setCollectionIsLoading(true);
    setCollectionIsLoaded(false);
    Reload().then(() => {
      fetchData();
      emitStatus('Reloaded ABI data. Fetching fresh data...');
    });
  };

  useHotkeys([['mod+r', handleRefresh]]);

  type IndexableAbi = types.Abi & { [key: string]: unknown };
  const handleAction = (address: string, _isDeleted: boolean) => {
    const addressStr = address;
    setProcessingAddresses((prev) => new Set(prev).add(addressStr));
    try {
      const original = [...downloadedAbis];

      // Optimistic UI Update - remove the deleted item
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
            listType,
            pagination.currentPage * pagination.pageSize,
            pagination.pageSize,
            sortArgument as sorting.SortDef,
            currentFilter,
          );

          setDownloadedAbis((result.Abis || []) as IndexableAbi[]);
          setTotalItems(result.ExpectedTotal || 0);

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

  const renderActions = (item: types.Abi | types.Function) => {
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

  let abiColumns: FormField<types.Abi>[] = [
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
  if (listType === 'Known') {
    abiColumns = abiColumns.slice(1);
  } else {
    abiColumns.push({
      key: 'actions',
      header: 'Actions',
      width: 'col-actions',
      render: (row: types.Abi) => renderActions(row),
    });
  }

  const functionColumns: FormField<types.Function>[] = [
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

  type TableRow = (types.Abi | types.Function) & {
    [key: string]: unknown;
  };

  const handleTableSubmit = (formData: Record<string, unknown>) => {
    Log(`Table submitted: ${formData}`);
  };

  const renderTable = (currentListType: AbisListType) => {
    let data: TableRow[] = [];
    let columns: FormField<TableRow>[] = [];
    if (currentListType === 'Downloaded') {
      data = downloadedAbis as TableRow[];
      columns = abiColumns as FormField<TableRow>[];
    } else if (currentListType === 'Known') {
      data = knownAbis as TableRow[];
      columns = abiColumns as FormField<TableRow>[];
    } else if (currentListType === 'Functions') {
      data = functions as TableRow[];
      columns = functionColumns as FormField<TableRow>[];
    } else if (currentListType === 'Events') {
      data = events as TableRow[];
      columns = functionColumns as FormField<TableRow>[];
    }

    return (
      <TableProvider>
        <div className="tableContainer">
          <Table<TableRow>
            columns={columns}
            data={data}
            loading={loading}
            collectionIsLoading={collectionIsLoading}
            collectionIsLoaded={collectionIsLoaded}
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

  const tabsConfig = [
    {
      label: 'Downloaded',
      value: 'Downloaded',
      content: renderTable('Downloaded'),
    },
    {
      label: 'Known',
      value: 'Known',
      content: renderTable('Known'),
    },
    {
      label: 'Functions',
      value: 'Functions',
      content: renderTable('Functions'),
    },
    {
      label: 'Events',
      value: 'Events',
      content: renderTable('Events'),
    },
  ];

  return (
    <div className="abisView">
      <TabView tabs={tabsConfig} route="/abis" />
      {error && (
        <div className="error-message-placeholder">
          <h3>{`Error fetching ${listType}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
    </div>
  );
};
// ADD_ROUTE
