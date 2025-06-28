// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetExportsPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
// prettier-ignore
import { useActionConfig } from '@hooks';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { exports, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, exportsFacets } from './facets';

// === END SECTION 1 ===

export const Exports = () => {
  // === SECTION 2: Hook Initialization ===
  const createPayload = usePayload();

  const activeFacetHook = useActiveFacet({
    facets: exportsFacets,
    defaultFacet: DEFAULT_FACET,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<exports.ExportsPage | null>(null);
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
  // === END SECTION 2 ===

  // === SECTION 3: Refs & Effects Setup ===
  const dataFacetRef = useRef(getCurrentDataFacet());
  useEffect(() => {
    dataFacetRef.current = getCurrentDataFacet();
  }, [getCurrentDataFacet]);
  // === END SECTION 3 ===

  // === SECTION 4: Data Fetching Logic ===
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
      handleError(err, `Failed to fetch ${getCurrentDataFacet()}`);
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
    getCurrentDataFacet,
  ]);

  const currentData = useMemo(() => {
    if (!pageData) return [];

    const facet = getCurrentDataFacet();
    switch (facet) {
      case types.DataFacet.STATEMENTS:
        return pageData.statements || [];
      case types.DataFacet.BALANCES:
        return pageData.balances || [];
      case types.DataFacet.TRANSFERS:
        return pageData.transfers || [];
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
  // === END SECTION 4 ===

  // === SECTION 5: Event Handling ===
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
        // The data will reload when the DataLoaded event is fired.
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, createPayload, handleError]);

  useHotkeys([['mod+r', handleReload]]);
  // === END SECTION 5 ===

  // === SECTION 6: CRUD Operations ===
  const actionConfig = useActionConfig({
    operations: [],
  });
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE
  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(pageData?.facet || getCurrentDataFacet());
    return actionConfig.injectActionColumn(baseColumns, () => null);
  }, [pageData?.facet, getCurrentDataFacet, actionConfig]);
  // EXISTING_CODE
  // === END SECTION 7 ===

  // === SECTION 8: Tab Configuration ===
  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
      />
    ),
    // prettier-ignore
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      viewStateKey,
    ],
  );

  const tabs = useMemo(
    () =>
      availableFacets.map((facetConfig: DataFacetConfig) => ({
        label: facetConfig.label,
        value: facetConfig.id,
        content: perTabTable,
      })),
    [availableFacets, perTabTable],
  );
  // === END SECTION 8 ===

  // === SECTION 9: Render/JSX ===
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

// EXISTING_CODE
// EXISTING_CODE
