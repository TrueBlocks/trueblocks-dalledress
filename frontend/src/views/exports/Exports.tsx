import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetExportsPage, LoadExportsData, ResetExportsData } from '@app';
import { BaseTab, usePagination } from '@components';
import { TableKey, useFiltering, useSorting } from '@contexts';
import { useActiveProject, useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { exports, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumnsForExports } from './columns';

export const Exports = () => {
  const { lastTab, effectiveAddress, effectiveChain } = useActiveProject();
  const [pageData, setPageData] = useState<exports.ExportsPage | null>(null);
  const [state, setState] = useState<types.LoadState>();

  const [listKind, setListKind] = useState<types.ListKind>(
    lastTab[EXPORTS_ROUTE] || EXPORTS_DEFAULT_LIST_KIND,
  );
  const tableKey = useMemo(
    (): TableKey => ({ viewName: EXPORTS_ROUTE, tabName: listKind }),
    [listKind],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter } = useFiltering(tableKey);

  const listKindRef = useRef(listKind);
  const renderCnt = useRef(0);

  useEffect(() => {
    listKindRef.current = listKind;
  }, [listKind]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetExportsPage(
        listKindRef.current,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
        effectiveChain,
        effectiveAddress,
      );
      setState(result.state);
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
    effectiveChain,
    effectiveAddress,
  ]);

  const currentData = useMemo(() => {
    if (!pageData) return [];

    switch (listKind) {
      case types.ListKind.STATEMENTS:
        return pageData.statements || [];
      case types.ListKind.TRANSFERS:
        return pageData.transfers || [];
      case types.ListKind.BALANCES:
        return pageData.balances || [];
      case types.ListKind.TRANSACTIONS:
        return pageData.transactions || [];
      default:
        return pageData.transactions || [];
    }
  }, [pageData, listKind]);

  useEffect(() => {
    const currentTab = lastTab[EXPORTS_ROUTE];
    if (currentTab && currentTab !== listKindRef.current) {
      setListKind(currentTab);
    }
  }, [lastTab]);

  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'exports') {
        const eventListKind = payload.listKind as types.ListKind | undefined;
        if (eventListKind === listKindRef.current) {
          fetchData();
        }
      }
    },
  );

  useEffect(() => {
    fetchData();
  }, [fetchData, listKind]);

  // Load data when component mounts or listKind changes
  useEffect(() => {
    LoadExportsData(listKind, effectiveChain, effectiveAddress);
  }, [listKind, effectiveChain, effectiveAddress]);

  useHotkeys([
    [
      'mod+r',
      () => {
        ResetExportsData(listKind, effectiveChain, effectiveAddress).then(
          () => {
            fetchData();
          },
        );
      },
    ],
  ]);

  const handleSubmit = useCallback((_formData: Record<string, unknown>) => {
    // Exports are read-only, no submit action needed
  }, []);

  const currentColumns = useMemo(() => {
    const baseColumns = getColumnsForExports(
      pageData?.kind || EXPORTS_DEFAULT_LIST_KIND,
    );

    // Exports are read-only, so we filter out any actions column
    return baseColumns.filter((col) => col.key !== 'actions');
  }, [pageData?.kind]);

  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        onSubmit={handleSubmit}
        tableKey={tableKey}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      handleSubmit,
      tableKey,
    ],
  );

  const tabs = useMemo(
    () => [
      {
        label: types.ListKind.STATEMENTS,
        value: types.ListKind.STATEMENTS,
        content: perTabTable,
      },
      {
        label: types.ListKind.TRANSFERS,
        value: types.ListKind.TRANSFERS,
        content: perTabTable,
      },
      {
        label: types.ListKind.BALANCES,
        value: types.ListKind.BALANCES,
        content: perTabTable,
      },
      {
        label: types.ListKind.TRANSACTIONS,
        value: types.ListKind.TRANSACTIONS,
        content: perTabTable,
      },
    ],
    [perTabTable],
  );

  return (
    <div className="mainView">
      {(state as string) === '' && <div>{`state: ${state}`}</div>}
      <TabView tabs={tabs} route={EXPORTS_ROUTE} />
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

const EXPORTS_DEFAULT_LIST_KIND = types.ListKind.TRANSACTIONS;
const EXPORTS_ROUTE = '/exports';
