import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetExportsPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import {
  DataFacetConfig,
  useActiveFacet,
  useActiveProject,
  useEvent,
} from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { exports, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

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
      tabName: activeFacetHook.getCurrentDataFacet(),
    }),
    [activeFacetHook],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);

  const dataFacetRef = useRef(activeFacetHook.getCurrentDataFacet());
  const renderCnt = useRef(0);

  useEffect(() => {
    dataFacetRef.current = activeFacetHook.getCurrentDataFacet();
  }, [activeFacetHook]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetExportsPage(
        types.Payload.createFrom({
          dataFacet: dataFacetRef.current,
          chain: effectiveChain,
          address: effectiveAddress,
        }),
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setState(result.state);
      setPageData(result);
      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      handleError(err, `Failed to fetch ${dataFacetRef.current}`);
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

    const currentDataFacet =
      activeFacetHook.getCurrentDataFacet() as types.DataFacet;
    switch (currentDataFacet) {
      case types.DataFacet.STATEMENTS:
        return pageData.statements || [];
      case types.DataFacet.TRANSFERS:
        return pageData.transfers || [];
      case types.DataFacet.BALANCES:
        return pageData.balances || [];
      case types.DataFacet.TRANSACTIONS:
        return pageData.transactions || [];
      case types.DataFacet.WITHDRAWALS:
        return pageData.withdrawals || [];
      default:
        return pageData.transactions || [];
    }
  }, [pageData, activeFacetHook]);

  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'exports') {
        const eventDataFacet = payload.dataFacet as types.DataFacet | undefined;
        if (eventDataFacet === dataFacetRef.current) {
          fetchData();
        }
      }
    },
  );

  useEffect(() => {
    fetchData();
  }, [fetchData, activeFacetHook.activeFacet]);

  useHotkeys([
    [
      'mod+r',
      () => {
        const currentDataFacet =
          activeFacetHook.getCurrentDataFacet() as types.DataFacet;
        Reload(currentDataFacet, effectiveChain, effectiveAddress).then(() => {
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
      pageData?.facet || activeFacetHook.getCurrentDataFacet(),
    );

    // Exports are read-only, so we filter out any actions column
    return baseColumns.filter((col) => col.key !== 'actions');
  }, [pageData?.facet, activeFacetHook]);

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
          <h3>{`Error fetching ${activeFacetHook.getCurrentDataFacet()}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      {renderCnt.current > 0 && <div>{`renderCnt: ${renderCnt.current}`}</div>}
    </div>
  );
};
