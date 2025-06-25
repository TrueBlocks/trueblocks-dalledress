// === SECTION 1: Imports & Dependencies ===
// EXISTING_CODE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetExportsPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { exports, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumns } from './columns';
import {
  EXPORTS_DEFAULT_FACET,
  EXPORTS_ROUTE as ROUTE,
  exportsFacets,
} from './facets';

// EXISTING_CODE
// === END SECTION 1 ===

export const Exports = () => {
  // === SECTION 2: Hook Initialization ===
  const createPayload = usePayload();
  // EXISTING_CODE
  const [pageData, setPageData] = useState<exports.ExportsPage | null>(null);

  const activeFacetHook = useActiveFacet({
    facets: exportsFacets,
    defaultFacet: EXPORTS_DEFAULT_FACET,
    viewRoute: ROUTE,
  });

  const { getCurrentDataFacet } = activeFacetHook;

  const viewStateKey = useMemo(
    (): ViewStateKey => ({
      viewName: ROUTE,
      tabName: getCurrentDataFacet(),
    }),
    [getCurrentDataFacet],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);
  // EXISTING_CODE
  // === END SECTION 2 ===

  // === SECTION 3: Refs & Effects Setup ===
  // EXISTING_CODE
  // EXISTING_CODE
  const dataFacetRef = useRef(getCurrentDataFacet() as types.DataFacet);
  useEffect(() => {
    dataFacetRef.current = getCurrentDataFacet() as types.DataFacet;
  }, [getCurrentDataFacet]);
  // === END SECTION 3 ===

  // === SECTION 4: Data Fetching Logic ===
  // EXISTING_CODE
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetExportsPage(
        createPayload(dataFacetRef.current),
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setPageData(result);
      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      handleError(err, `Failed to fetch ${dataFacetRef.current}`);
    }
  }, [
    clearError,
    createPayload,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setTotalItems,
    handleError,
  ]);

  const currentData = useMemo(() => {
    if (!pageData) return [];

    const currentDataFacet = getCurrentDataFacet() as types.DataFacet;
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
      case types.DataFacet.ASSETS:
        return pageData.assets || [];
      case types.DataFacet.LOGS:
        return pageData.logs || [];
      case types.DataFacet.TRACES:
        return pageData.traces || [];
      case types.DataFacet.RECEIPTS:
        return pageData.receipts || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);
  // EXISTING_CODE
  // === END SECTION 4 ===

  // === SECTION 5: Event Handling ===
  // EXISTING_CODE
  // EXISTING_CODE
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'exports') {
        const eventDataFacet = payload.dataFacet;
        if (eventDataFacet === dataFacetRef.current) {
          fetchData();
        }
      }
    },
  );

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleReload = useCallback(async () => {
    try {
      Reload(createPayload(dataFacetRef.current)).then(() => {
        fetchData();
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, createPayload, fetchData, handleError]);

  useHotkeys([['mod+r', handleReload]]);
  // === END SECTION 5 ===

  // === SECTION 6: CRUD Operations ===
  // EXISTING_CODE
  // EXISTING_CODE
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE

  const handleSubmit = useCallback((_formData: Record<string, unknown>) => {
    // Exports are read-only, no submit action needed
  }, []);

  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(
      pageData?.facet || (getCurrentDataFacet() as types.DataFacet),
    );

    // Exports are read-only, so we filter out any actions column
    return baseColumns.filter((col) => col.key !== 'actions');
  }, [pageData?.facet, getCurrentDataFacet]);
  // EXISTING_CODE
  // === END SECTION 7 ===

  // === SECTION 8: Tab Configuration ===
  // EXISTING_CODE
  // EXISTING_CODE
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
  // === END SECTION 8 ===

  // === SECTION 9: Render/JSX ===
  // EXISTING_CODE
  // EXISTING_CODE
  const renderCnt = useRef(0);
  // renderCnt.current++;
  return (
    <div className="mainView">
      <TabView tabs={tabs} route={ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${getCurrentDataFacet()}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      {renderCnt.current > 0 && <div>{`renderCnt: ${renderCnt.current}`}</div>}
    </div>
  );
  // === END SECTION 9 ===
};
