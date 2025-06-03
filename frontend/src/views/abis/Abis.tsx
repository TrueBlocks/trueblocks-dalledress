// ADD_ROUTE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { BaseTab, usePagination } from '@components';
import type { FormField } from '@components';
import { TableKey, useAppContext, useFiltering, useSorting } from '@contexts';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { abis, msgs, types } from '@models';
import { EventsOn } from '@runtime';
import { Log, useEmitters } from '@utils';

import {
  ABIS_ROUTE,
  ACTION_MESSAGES,
  DEFAULT_LIST_KIND,
  getAbisPage,
  reload,
} from './';
import './Abis.css';
import {
  ABI_COLUMNS,
  FUNCTION_COLUMNS,
  KNOWN_ABI_COLUMNS,
} from './columnDefinitions';

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

  const [pageData, setPageData] = useState<abis.AbisPage | null>(null);
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
      setPageData(result);
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

  // const handleAction = (address: string) => {
  //   // setProcessingAddresses((prev) => new Set(prev).add(address));
  //   try {
  //     const original = [...(pageData?.abis || [])];
  //     const optimisticValues = original.filter((abi) => {
  //       const nameAddress =
  //         typeof abi.address === 'string' ? abi.address : String(abi.address);
  //       return nameAddress !== address;
  //     });
  //     setPageData((prev) => {
  //       if (!prev) return null;
  //       // Create a new AbisPage instance to preserve methods like convertValues
  //       return new abis.AbisPage({
  //         ...prev,
  //         abis: optimisticValues as types.Abi[],
  //       });
  //     });
  //
  //     removeAbi(address)
  //       .then(async () => {
  //         const result = await getAbisPage(
  //           listKind,
  //           pagination.currentPage * pagination.pageSize,
  //           pagination.pageSize,
  //           sort,
  //           filter,
  //         );
  //         setPageData(result);
  //         setTotalItems(result.totalItems || 0);
  //         emitStatus(ACTION_MESSAGES.DELETE_SUCCESS(address));
  //       })
  //       .catch((err) => {
  //         // If there's an error, revert the optimistic update
  //         setPageData((prev) => {
  //           if (!prev) return null;
  //           return new abis.AbisPage({
  //             ...prev,
  //             abis: optimisticValues as types.Abi[],
  //           });
  //         });
  //         emitError(err);
  //         Log(`Error in handleAction: ${err}`);
  //       });
  //   } finally {
  //     // Always clean up the processing state
  //     // setProcessingAddresses((prev) => {
  //     //   const newSet = new Set(prev);
  //     //   newSet.delete(address);
  //     //   return newSet;
  //     // });
  //   }
  // };

  const handleTableSubmit = (formData: Record<string, unknown>) => {
    Log(`Table submitted: ${formData}`);
  };

  const getCurrentData = useCallback(() => {
    return pageData?.abis || pageData?.functions || [];
  }, [pageData]);

  const getColumnsForKind = (
    kind: string,
  ): FormField<Record<string, unknown>>[] => {
    switch (kind) {
      case 'Downloaded':
        return ABI_COLUMNS as unknown as FormField<Record<string, unknown>>[];
      case 'Known':
        return KNOWN_ABI_COLUMNS as unknown as FormField<
          Record<string, unknown>
        >[];
      case 'Functions':
      case 'Events':
        return FUNCTION_COLUMNS as unknown as FormField<
          Record<string, unknown>
        >[];
      default:
        return ABI_COLUMNS as unknown as FormField<Record<string, unknown>>[];
    }
  };

  const renderTable = () => {
    return (
      <BaseTab
        data={getCurrentData() as unknown as Record<string, unknown>[]}
        columns={getColumnsForKind(pageData?.kind || 'Downloaded')}
        loading={!!pageData?.isLoading}
        error={error}
        onSubmit={handleTableSubmit}
        tableKey={tableKey}
      />
    );
  };

  const createOneTab = (listKind: types.ListKind) => {
    return {
      label: listKind,
      value: listKind,
      content: renderTable(),
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
