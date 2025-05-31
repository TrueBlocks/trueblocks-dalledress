// ADD_ROUTE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { DeleteAbi, GetAbisPage, Reload } from '@app';
import { Table, TableProvider, usePagination } from '@components';
import { TableKey, useAppContext } from '@contexts';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { msgs, sorting, types } from '@models';
import { EventsOn } from '@runtime';
import { Log, useEmitters } from '@utils';
import { getTableConfig } from 'src/views/abis/utilities';
import { useLocation } from 'wouter';

import './Abis.css';
import { ABIS_ROUTE, ACTION_MESSAGES, DEFAULT_LIST_KIND } from './constants';
import { AbiRow, IndexedAbi, IndexedFunction, TableConfigProps } from './types';

//--------------------------------------------------------------------
export const Abis = () => {
  const { emitStatus, emitError } = useEmitters();
  const { lastTab, setSelectedAddress } = useAppContext();
  const [, setLocation] = useLocation();

  const [listKind, setListKind] = useState<types.ListKind>(
    (lastTab[ABIS_ROUTE] as types.ListKind) || DEFAULT_LIST_KIND,
  );
  const listKindRef = useRef(listKind);
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
    return { viewName: ABIS_ROUTE, tabName: listKind };
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
          setDownloadedAbis((result.abis as IndexedAbi[]) || []);
          break;
        case types.ListKind.KNOWN:
          setIsKnownLoaded(result.isLoaded);
          setKnownAbis((result.abis as IndexedAbi[]) || []);
          break;
        case types.ListKind.FUNCTIONS:
          setIsFuncsLoaded(result.isLoaded);
          setFunctions((result.functions as IndexedFunction[]) || []);
          break;
        case types.ListKind.EVENTS:
          setIsEventsLoaded(result.isLoaded);
          setEvents((result.functions as IndexedFunction[]) || []);
          break;
        default:
          throw new Error(`Unknown list kind: ${listKind}`);
      }

      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      const e = err instanceof Error ? err : new Error(String(err));
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
    listKindRef.current = listKind;
  }, [listKind]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  useEffect(() => {
    const eventName = msgs.EventType.DATA_LOADED;
    const unlisten = EventsOn(eventName, (payload: unknown) => {
      const eventPayload = payload as types.DataLoadedPayload;
      if (eventPayload) {
        switch (listKindRef.current) {
          case types.ListKind.DOWNLOADED:
            setIsDownloadedLoaded(eventPayload.isLoaded);
            fetchData();
            break;
          case types.ListKind.KNOWN:
            setIsKnownLoaded(eventPayload.isLoaded);
            fetchData();
            break;
          case types.ListKind.FUNCTIONS:
            setIsFuncsLoaded(eventPayload.isLoaded);
            fetchData();
            break;
          case types.ListKind.EVENTS:
            setIsEventsLoaded(eventPayload.isLoaded);
            fetchData();
            break;
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
    const currentTabLabel = lastTab[ABIS_ROUTE] as types.ListKind | undefined;
    if (currentTabLabel && currentTabLabel !== listKindRef.current) {
      setListKind(currentTabLabel);
    }
  }, [lastTab]);

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
          emitStatus(ACTION_MESSAGES.RELOAD_STATUS);
        });
      },
    ],
  ]);

  const handleAction = (address: string) => {
    setProcessingAddresses((prev) => new Set(prev).add(address));
    try {
      const original = [...downloadedAbis];
      const optimisticNames = original.filter((abi) => {
        const nameAddress =
          typeof abi.address === 'string' ? abi.address : String(abi.address);
        return nameAddress !== address;
      });
      setDownloadedAbis(optimisticNames as IndexedAbi[]);

      DeleteAbi(address)
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
          setTotalItems(result.totalItems || 0);

          emitStatus(ACTION_MESSAGES.DELETE_SUCCESS(address));
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
        newSet.delete(address);
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

    const shouldShowLoading = (() => {
      switch (listKind) {
        case types.ListKind.DOWNLOADED:
          return loading && downloadedAbis.length === 0;
        case types.ListKind.KNOWN:
          return loading && knownAbis.length === 0;
        case types.ListKind.FUNCTIONS:
          return loading && functions.length === 0;
        case types.ListKind.EVENTS:
          return loading && events.length === 0;
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
            loading={shouldShowLoading}
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
      <TabView tabs={tabs} route={ABIS_ROUTE} />
      {error && (
        <div className="error-message-placeholder">
          <h3>{`Error fetching ${listKind}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
    </div>
  );
};

// ADD_ROUTE
