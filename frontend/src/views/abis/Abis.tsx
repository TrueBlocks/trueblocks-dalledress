/*
TableConfigProps
getTableConfig
*/
// ADD_ROUTE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetAbisPage, Reload, RemoveAbi } from '@app';
import { usePagination } from '@components';
import { TableKey, useAppContext, useFiltering, useSorting } from '@contexts';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { msgs, sorting, types } from '@models';
import { EventsOn } from '@runtime';
import { Log, useEmitters } from '@utils';

import './Abis.css';
import { ABIS_ROUTE, ACTION_MESSAGES, DEFAULT_LIST_KIND } from './constants';
import { DownloadedTab, EventsTab, FunctionsTab, KnownTab } from './index';
import { IndexedAbi, IndexedFunction } from './types';

//--------------------------------------------------------------------
export const Abis = () => {
  const { emitStatus, emitError } = useEmitters();
  // const { lastTab, setSelectedAddress } = useAppContext();
  // const [, setLocation] = useLocation();
  const { lastTab } = useAppContext();

  const [listKind, setListKind] = useState<types.ListKind>(
    (lastTab[ABIS_ROUTE] as types.ListKind) || DEFAULT_LIST_KIND,
  );
  const listKindRef = useRef(listKind);
  const [downloaded, setDownloaded] = useState<IndexedAbi[]>([]);
  const [known, setKnown] = useState<IndexedAbi[]>([]);
  const [functions, setFunctions] = useState<IndexedFunction[]>([]);
  const [events, setEvents] = useState<IndexedFunction[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  // const [processingAddresses, setProcessingAddresses] = useState<Set<string>>(
  //   new Set(),
  // );
  // const [isDownloadedLoaded, setIsDownloadedLoaded] = useState<boolean>(false);
  // const [isKnownLoaded, setIsKnownLoaded] = useState<boolean>(false);
  // const [isFuncsLoaded, setIsFuncsLoaded] = useState<boolean>(false);
  // const [isEventsLoaded, setIsEventsLoaded] = useState<boolean>(false);
  const [, setProcessingAddresses] = useState<Set<string>>(new Set());
  const [, setIsDownloadedLoaded] = useState<boolean>(false);
  const [, setIsKnownLoaded] = useState<boolean>(false);
  const [, setIsFuncsLoaded] = useState<boolean>(false);
  const [, setIsEventsLoaded] = useState<boolean>(false);

  const tableKey = useMemo((): TableKey => {
    return { viewName: ABIS_ROUTE, tabName: listKind };
  }, [listKind]);

  const { pagination, setTotalItems } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter } = useFiltering(tableKey);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const sortArgument = sort === null ? undefined : sort;
      const result = await GetAbisPage(
        listKind,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sortArgument as sorting.SortDef,
        filter,
      );

      switch (listKind) {
        case types.ListKind.DOWNLOADED:
          setIsDownloadedLoaded(result.isLoaded);
          setDownloaded((result.abis as IndexedAbi[]) || []);
          break;
        case types.ListKind.KNOWN:
          setIsKnownLoaded(result.isLoaded);
          setKnown((result.abis as IndexedAbi[]) || []);
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
    sort,
    listKind,
    pagination.currentPage,
    pagination.pageSize,
    filter,
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
    const unlisten = EventsOn(eventName, (payload: types.DataLoadedPayload) => {
      if (payload) {
        switch (listKindRef.current) {
          case types.ListKind.DOWNLOADED:
            setIsDownloadedLoaded(payload.isLoaded);
            fetchData();
            break;
          case types.ListKind.KNOWN:
            setIsKnownLoaded(payload.isLoaded);
            fetchData();
            break;
          case types.ListKind.FUNCTIONS:
            setIsFuncsLoaded(payload.isLoaded);
            fetchData();
            break;
          case types.ListKind.EVENTS:
            setIsEventsLoaded(payload.isLoaded);
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
      const original = [...downloaded];
      const optimisticNames = original.filter((abi) => {
        const nameAddress =
          typeof abi.address === 'string' ? abi.address : String(abi.address);
        return nameAddress !== address;
      });
      setDownloaded(optimisticNames as IndexedAbi[]);

      RemoveAbi(address)
        .then(async () => {
          // If API call is successful, refresh the data to get the definitive state
          const sortArgument = sort === null ? undefined : sort;
          const result = await GetAbisPage(
            listKind,
            pagination.currentPage * pagination.pageSize,
            pagination.pageSize,
            sortArgument as sorting.SortDef,
            filter,
          );

          setDownloaded((result.abis || []) as IndexedAbi[]);
          setTotalItems(result.totalItems || 0);

          emitStatus(ACTION_MESSAGES.DELETE_SUCCESS(address));
        })
        .catch((err) => {
          // If there's an error, revert the optimistic update
          setDownloaded(optimisticNames as IndexedAbi[]);
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
    switch (listKind) {
      case types.ListKind.DOWNLOADED:
        return (
          <DownloadedTab
            data={downloaded}
            loading={loading && downloaded.length === 0}
            error={error}
            onSubmit={handleTableSubmit}
            onDelete={handleAction}
            tableKey={tableKey}
          />
        );
      case types.ListKind.KNOWN:
        return (
          <KnownTab
            data={known}
            loading={loading && known.length === 0}
            error={error}
            onSubmit={handleTableSubmit}
            tableKey={tableKey}
          />
        );
      case types.ListKind.FUNCTIONS:
        return (
          <FunctionsTab
            data={functions}
            loading={loading && functions.length === 0}
            error={error}
            onSubmit={handleTableSubmit}
            tableKey={tableKey}
          />
        );
      case types.ListKind.EVENTS:
        return (
          <EventsTab
            data={events}
            loading={loading && events.length === 0}
            error={error}
            onSubmit={handleTableSubmit}
            tableKey={tableKey}
          />
        );
      default:
        return <div>Unknown list kind: {listKind}</div>;
    }
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
