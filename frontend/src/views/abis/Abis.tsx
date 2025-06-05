// ADD_ROUTE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { AbisCrud } from '@app';
import { BaseTab, usePagination } from '@components';
import { TableKey, useAppContext, useFiltering, useSorting } from '@contexts';
import { useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { abis, crud, msgs, types } from '@models';
import { getAddressString, useEmitters, useErrorHandler } from '@utils';

import { Address } from '../../types/address';
import {
  ABIS_DEFAULT_LIST,
  ABIS_ROUTE,
  ACTION_MESSAGES,
  getAbisPage,
  getColumns,
  reload,
} from './';

export const Abis = () => {
  const { lastTab } = useAppContext();
  const [pageData, setPageData] = useState<abis.AbisPage | null>(null);

  const [listKind, setListKind] = useState<types.ListKind>(
    lastTab[ABIS_ROUTE] || ABIS_DEFAULT_LIST,
  );
  const tableKey = useMemo(
    (): TableKey => ({ viewName: ABIS_ROUTE, tabName: listKind }),
    [listKind],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter } = useFiltering(tableKey);
  const { emitStatus } = useEmitters();

  const listKindRef = useRef(listKind);
  const renderCnt = useRef(0);
  // renderCnt.current++;

  useEffect(() => {
    listKindRef.current = listKind;
  }, [listKind]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await getAbisPage(
        listKindRef.current,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setPageData(result);
      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      handleError(err, `Failed to fetch ${listKindRef.current}`);
    }
  }, [
    clearError,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setTotalItems,
    handleError,
  ]);

  const currentData = useMemo(() => {
    return pageData?.abis || pageData?.functions || [];
  }, [pageData?.abis, pageData?.functions]);

  useEffect(() => {
    const currentTab = lastTab[ABIS_ROUTE];
    if (currentTab && currentTab !== listKindRef.current) {
      setListKind(currentTab);
    }
  }, [lastTab]);

  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: types.DataLoadedPayload) => {
      if (payload?.listKind === listKindRef.current) {
        fetchData();
      }
    },
  );

  useEffect(() => {
    fetchData();
  }, [fetchData, listKind]);

  useHotkeys([
    [
      'mod+r',
      () => {
        reload().then(() => {
          fetchData();
        });
      },
    ],
  ]);

  // Optimistic delete action with proper type safety
  const _handleAction = (address: Address) => {
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
      AbisCrud(crud.Operation.REMOVE, {} as types.Abi, address)
        .then(async () => {
          const result = await getAbisPage(
            listKindRef.current,
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

  const handleSubmit = useCallback((_formData: Record<string, unknown>) => {
    // Log(`Table submitted: ${formData}`);
  }, []);

  const currentColumns = useMemo(
    () => getColumns(pageData?.kind || ABIS_DEFAULT_LIST),
    [pageData?.kind],
  );

  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isLoading}
        error={error}
        onSubmit={handleSubmit}
        tableKey={tableKey}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isLoading,
      error,
      handleSubmit,
      tableKey,
    ],
  );

  const tabs = useMemo(
    () => [
      {
        label: types.ListKind.DOWNLOADED,
        value: types.ListKind.DOWNLOADED,
        content: perTabTable,
      },
      {
        label: types.ListKind.KNOWN,
        value: types.ListKind.KNOWN,
        content: perTabTable,
      },
      {
        label: types.ListKind.FUNCTIONS,
        value: types.ListKind.FUNCTIONS,
        content: perTabTable,
      },
      {
        label: types.ListKind.EVENTS,
        value: types.ListKind.EVENTS,
        content: perTabTable,
      },
    ],
    [perTabTable],
  );

  return (
    <div className="mainView">
      <TabView tabs={tabs} route={ABIS_ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${listKind}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      {renderCnt.current > 0 && <div>{`renderCnt: ${renderCnt.current}`}</div>}
    </div>
  );
};

// ADD_ROUTE
