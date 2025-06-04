// ADD_ROUTE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { BaseTab, usePagination } from '@components';
import type { FormField } from '@components';
import { TableKey, useAppContext, useFiltering, useSorting } from '@contexts';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { abis, msgs, types } from '@models';
import { EventsOn } from '@runtime';
import { getAddressString, useEmitters, useErrorHandler } from '@utils';

import {
  ABIS_ROUTE,
  ACTION_MESSAGES,
  DEFAULT_LIST_KIND,
  getAbisPage,
  reload,
  removeAbi,
} from './';
import './Abis.css';

//--------------------------------------------------------------------
export const Abis = () => {
  const { emitStatus } = useEmitters();
  const { lastTab } = useAppContext();
  const { error, handleError, clearError } = useErrorHandler();

  const renderCount = useRef(0);
  renderCount.current++;

  const [listKind, setListKind] = useState<types.ListKind>(
    lastTab[ABIS_ROUTE] || DEFAULT_LIST_KIND,
  );
  const listKindRef = useRef(listKind);
  const tableKey = useMemo(
    (): TableKey => ({ viewName: ABIS_ROUTE, tabName: listKind }),
    [listKind],
  );

  const [pageData, setPageData] = useState<abis.AbisPage | null>(null);
  const { pagination, setTotalItems } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter } = useFiltering(tableKey);

  // Fetch data from backend and update state
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await getAbisPage(
        listKind,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setPageData(result);
      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      handleError(err, `Failed to fetch ${listKind}`);
    }
  }, [
    clearError,
    listKind,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setTotalItems,
    handleError,
  ]);

  useEffect(() => {
    listKindRef.current = listKind;
  }, [listKind]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Listen for backend data loaded events
  useEffect(() => {
    const eventName = msgs.EventType.DATA_LOADED;
    const unlisten = EventsOn(eventName, (payload: types.DataLoadedPayload) => {
      if (payload?.listKind == listKind) {
        // && payload?.reason === 'initial') {
        fetchData();
      }
    });
    return () => {
      if (typeof unlisten === 'function') unlisten();
    };
  }, [fetchData, listKind]);

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

  // Optimistic delete action with proper type safety
  const _handleAction = (address: string) => {
    clearError();
    try {
      const original = [...(pageData?.abis || [])];
      const optimisticValues = original.filter((abi) => {
        const abiAddress = getAddressString(abi.address);
        return abiAddress !== address;
      });
      setPageData((prev) => {
        if (!prev) return null;
        return new abis.AbisPage({
          ...prev,
          abis: optimisticValues,
        });
      });
      removeAbi(address)
        .then(async () => {
          const result = await getAbisPage(
            listKind,
            pagination.currentPage * pagination.pageSize,
            pagination.pageSize,
            sort,
            filter,
          );
          setPageData(result);
          setTotalItems(result.totalItems || 0);
          emitStatus(ACTION_MESSAGES.DELETE_SUCCESS(address));
        })
        .catch((err) => {
          setPageData((prev) => {
            if (!prev) return null;
            return new abis.AbisPage({
              ...prev,
              abis: original,
            });
          });
          handleError(err, 'handleAction');
        });
    } finally {
      // Always clean up the processing state if needed
    }
  };

  const handleTableSubmit = useCallback(
    (_formData: Record<string, unknown>) => {
      // Log(`Table submitted: ${formData}`);
    },
    [],
  );

  // Get current data for the active tab
  const currentData = useMemo(() => {
    return pageData?.abis || pageData?.functions || [];
  }, [pageData?.abis, pageData?.functions]);

  const getColumnsForKind = useCallback(
    (listKind: types.ListKind): FormField<Record<string, unknown>>[] => {
      switch (listKind) {
        case 'Functions':
        case 'Events':
          return getFunctionColumns();
        case 'Known':
          return getAbisColumns().filter((col) => {
            const skip = col.key !== 'address';
            return skip;
          });
        case 'Downloaded':
        // fallthrough intended
        default:
          return getAbisColumns();
      }
    },
    [],
  );

  const currentColumns = useMemo(
    () => getColumnsForKind(pageData?.kind || types.ListKind.DOWNLOADED),
    [pageData?.kind, getColumnsForKind],
  );

  const memoizedTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isLoading}
        error={error}
        onSubmit={handleTableSubmit}
        tableKey={tableKey}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isLoading,
      error,
      handleTableSubmit,
      tableKey,
    ],
  );

  const tabs = useMemo(
    () => [
      {
        label: types.ListKind.DOWNLOADED,
        value: types.ListKind.DOWNLOADED,
        content: memoizedTable,
      },
      {
        label: types.ListKind.KNOWN,
        value: types.ListKind.KNOWN,
        content: memoizedTable,
      },
      {
        label: types.ListKind.FUNCTIONS,
        value: types.ListKind.FUNCTIONS,
        content: memoizedTable,
      },
      {
        label: types.ListKind.EVENTS,
        value: types.ListKind.EVENTS,
        content: memoizedTable,
      },
    ],
    [memoizedTable],
  );

  return (
    <div className="abisView">
      <TabView tabs={tabs} route={ABIS_ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${listKind}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      <div>{`renderCnt: ${renderCount.current}`}</div>
    </div>
  );
};

const getAbisColumns = (): FormField[] => [
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
  },
  {
    key: 'nFunctions',
    header: 'Functions',
    sortable: true,
    type: 'number',
  },
  {
    key: 'nEvents',
    header: 'Events',
    sortable: true,
    type: 'number',
  },
  {
    key: 'nErrors',
    header: 'Errors',
    sortable: true,
    type: 'number',
  },
];

const getFunctionColumns = (): FormField[] => [
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
];
