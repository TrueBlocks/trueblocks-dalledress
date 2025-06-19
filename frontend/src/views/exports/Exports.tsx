import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetExportsPage, LoadExportsData, ResetExportsData } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { useActiveFacet, useActiveProject, useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { exports, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';
import { getColumnsForExports } from './columns';
import {
  EXPORTS_DEFAULT_FACET,
  EXPORTS_ROUTE as ROUTE,
  exportsFacets,
} from './exportsFacets';

export const Exports = () => {
  const { effectiveAddress, effectiveChain } = useActiveProject();
  const [pageData, setPageData] = useState<exports.ExportsPage | null>(null);
  const [state, setState] = useState<types.LoadState>();

  const activeFacetHook = useActiveFacet({
    facets: exportsFacets,
    defaultFacet: EXPORTS_DEFAULT_FACET,
    viewRoute: ROUTE,
  });

  const viewStateKey = useMemo(
    (): ViewStateKey => ({
      viewName: ROUTE,
      tabName: activeFacetHook.getCurrentListKind(),
    }),
    [activeFacetHook],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);

  const listKindRef = useRef(activeFacetHook.getCurrentListKind());
  const renderCnt = useRef(0);

  useEffect(() => {
    listKindRef.current = activeFacetHook.getCurrentListKind();
  }, [activeFacetHook]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetExportsPage(
        listKindRef.current as types.ListKind,
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

    const currentListKind =
      activeFacetHook.getCurrentListKind() as types.ListKind;
    switch (currentListKind) {
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
  }, [pageData, activeFacetHook]);

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
  }, [fetchData, activeFacetHook.activeFacet]);

  // Load data when component mounts or facet changes
  useEffect(() => {
    const currentListKind =
      activeFacetHook.getCurrentListKind() as types.ListKind;
    LoadExportsData(currentListKind, effectiveChain, effectiveAddress);
  }, [activeFacetHook, effectiveChain, effectiveAddress]);

  useHotkeys([
    [
      'mod+r',
      () => {
        const currentListKind =
          activeFacetHook.getCurrentListKind() as types.ListKind;
        ResetExportsData(
          currentListKind,
          effectiveChain,
          effectiveAddress,
        ).then(() => {
          fetchData();
        });
      },
    ],
  ]);

  const handleSubmit = useCallback((_formData: Record<string, unknown>) => {
    // Exports are read-only, no submit action needed
  }, []);

  const currentColumns = useMemo(() => {
    const baseColumns = getColumnsForExports(
      pageData?.kind || activeFacetHook.getCurrentListKind(),
    );

    // Exports are read-only, so we filter out any actions column
    return baseColumns.filter((col) => col.key !== 'actions');
  }, [pageData?.kind, activeFacetHook]);

  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        onSubmit={handleSubmit}
        viewStateKey={viewStateKey}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      handleSubmit,
      viewStateKey,
    ],
  );

  const tabs = useMemo(
    () =>
      activeFacetHook.availableFacets.map((facetConfig: DataFacetConfig) => ({
        label: facetConfig.label,
        value: facetConfig.id,
        content: perTabTable,
      })),
    [activeFacetHook.availableFacets, perTabTable],
  );

  return (
    <div className="mainView">
      {(state as string) === '' && <div>{`state: ${state}`}</div>}
      <TabView tabs={tabs} route={ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${activeFacetHook.getCurrentListKind()}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      {renderCnt.current > 0 && <div>{`renderCnt: ${renderCnt.current}`}</div>}
    </div>
  );
};
