// ADD_ROUTE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { usePagination } from '@components';
import { TableKey, useAppContext, useFiltering, useSorting } from '@contexts';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { msgs, types } from '@models';
import { EventsOn } from '@runtime';
import { Log, useEmitters } from '@utils';

import {
  ABIS_ROUTE,
  ACTION_MESSAGES,
  DEFAULT_LIST_KIND,
  DownloadedTab,
  EventsTab,
  FunctionsTab,
  KnownTab,
  getAbisPage,
  reload,
  removeAbi,
} from './';
import './Abis.css';

//--------------------------------------------------------------------
export const Abis = () => {
  const { emitStatus, emitError } = useEmitters();
  // const { lastTab, setSelectedAddress } = useAppContext();
  // const [, setLocation] = useLocation();
  const { lastTab } = useAppContext();

  const [listKind, setListKind] = useState<types.ListKind>(
    lastTab[ABIS_ROUTE] || DEFAULT_LIST_KIND,
  );
  const listKindRef = useRef(listKind);
  const tableKey = useMemo((): TableKey => {
    return { viewName: ABIS_ROUTE, tabName: listKind };
  }, [listKind]);

  const [downloaded, setDownloaded] = useState<types.Abi[]>([]);
  const [known, setKnown] = useState<types.Abi[]>([]);
  const [functions, setFunctions] = useState<types.Function[]>([]);
  const [events, setEvents] = useState<types.Function[]>([]);

  // const [processingAddresses, setProcessingAddresses] = useState<Set<string>>(
  //   new Set(),
  // );
  const [, setProcessingAddresses] = useState<Set<string>>(new Set());
  const [error, setError] = useState<Error | null>(null);
  const { pagination, setTotalItems } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter } = useFiltering(tableKey);

  const fetchData = useCallback(async () => {
    setError(null);
    try {
      const result = await getAbisPage(
        listKind,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );

      switch (listKind) {
        case types.ListKind.DOWNLOADED:
          setDownloaded(result.abis || []);
          break;
        case types.ListKind.KNOWN:
          setKnown(result.abis || []);
          break;
        case types.ListKind.FUNCTIONS:
          setFunctions(result.functions || []);
          break;
        case types.ListKind.EVENTS:
          setEvents(result.functions || []);
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
            fetchData();
            break;
          case types.ListKind.KNOWN:
            fetchData();
            break;
          case types.ListKind.FUNCTIONS:
            fetchData();
            break;
          case types.ListKind.EVENTS:
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
    const currentTabLabel = lastTab[ABIS_ROUTE];
    if (currentTabLabel && currentTabLabel !== listKindRef.current) {
      setListKind(currentTabLabel);
    }
  }, [lastTab]);

  useHotkeys([
    [
      'mod+r',
      () => {
        reload().then(() => {
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
      const optimisticValues = original.filter((abi) => {
        const nameAddress =
          typeof abi.address === 'string' ? abi.address : String(abi.address);
        return nameAddress !== address;
      });
      setDownloaded(optimisticValues as types.Abi[]);

      removeAbi(address)
        .then(async () => {
          const result = await getAbisPage(
            listKind,
            pagination.currentPage * pagination.pageSize,
            pagination.pageSize,
            sort,
            filter,
          );
          setDownloaded(result.abis || []);
          setTotalItems(result.totalItems || 0);
          emitStatus(ACTION_MESSAGES.DELETE_SUCCESS(address));
        })
        .catch((err) => {
          // If there's an error, revert the optimistic update
          setDownloaded(optimisticValues as types.Abi[]);
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
            loading={downloaded.length === 0}
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
            loading={known.length === 0}
            error={error}
            onSubmit={handleTableSubmit}
            tableKey={tableKey}
          />
        );
      case types.ListKind.FUNCTIONS:
        return (
          <FunctionsTab
            data={functions}
            loading={functions.length === 0}
            error={error}
            onSubmit={handleTableSubmit}
            tableKey={tableKey}
          />
        );
      case types.ListKind.EVENTS:
        return (
          <EventsTab
            data={events}
            loading={events.length === 0}
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
