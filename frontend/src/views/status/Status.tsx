// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetStatusPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import {
  DataFacetConfig,
  toPageDataProp,
  useActiveFacet,
  useColumns,
  useEvent,
  usePayload,
} from '@hooks';
import { FormView, TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { status } from '@models';
import { msgs, types } from '@models';
import { Debugger, useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, statusFacets } from './facets';

export const Status = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();
  const activeFacetHook = useActiveFacet({
    facets: statusFacets,
    defaultFacet: DEFAULT_FACET,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<status.StatusPage | null>(null);
  const viewStateKey = useMemo(
    (): ViewStateKey => ({
      viewName: ROUTE,
      facetName: getCurrentDataFacet(),
    }),
    [getCurrentDataFacet],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);

  // === SECTION 3: Data Fetching ===
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetStatusPage(
        createPayload(getCurrentDataFacet()),
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
    getCurrentDataFacet,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setTotalItems,
    handleError,
  ]);

  const currentData = useMemo(() => {
    if (!pageData) return [];
    const facet = getCurrentDataFacet();
    switch (facet) {
      case types.DataFacet.STATUS:
        return pageData.status || [];
      case types.DataFacet.CACHES:
        return pageData.caches || [];
      case types.DataFacet.CHAINS:
        return pageData.chains || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);

  // === SECTION 4: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'status') {
        const eventDataFacet = payload.dataFacet;
        if (eventDataFacet === getCurrentDataFacet()) {
          fetchData();
        }
      }
    },
  );

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleReload = useCallback(async () => {
    clearError();
    try {
      Reload(createPayload(getCurrentDataFacet())).then(() => {
        // The data will reload when the DataLoaded event is fired.
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [clearError, getCurrentDataFacet, createPayload, handleError]);

  useHotkeys([['mod+r', handleReload]]);

  // === SECTION 6: UI Configuration ===
  const currentColumns = useColumns(
    getColumns(getCurrentDataFacet()),
    {
      showActions: false,
      actions: [],
      getCanRemove: useCallback((_row: unknown) => false, []),
    },
    {},
    toPageDataProp(pageData),
    { rowActions: [] },
  );

  const isForm = useCallback((facet: types.DataFacet) => {
    switch (facet) {
      case types.DataFacet.STATUS:
        return true;
      default:
        return false;
    }
  }, []);

  const perTabContent = useMemo(() => {
    const facet = getCurrentDataFacet();
    if (isForm(facet)) {
      const statusData = currentData[0] as unknown as Record<string, unknown>;
      if (!statusData) {
        return <div>No status data available</div>;
      }
      return (
        <FormView
          title="Status Information"
          formFields={getColumns(getCurrentDataFacet()).map((field) => {
            let value = statusData?.[field.name as string] ?? field.value;
            if (
              value !== undefined &&
              typeof value !== 'string' &&
              typeof value !== 'number' &&
              typeof value !== 'boolean'
            ) {
              value = JSON.stringify(value);
            }
            if (
              value !== undefined &&
              typeof value !== 'string' &&
              typeof value !== 'number' &&
              typeof value !== 'boolean'
            ) {
              value = undefined;
            }
            return {
              ...field,
              value,
              readOnly: true,
            };
          })}
          onSubmit={() => {}}
        />
      );
    }

    return (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        headerActions={[]}
      />
    );
  }, [
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    viewStateKey,
    isForm,
    getCurrentDataFacet,
  ]);

  const tabs = useMemo(
    () =>
      availableFacets.map((facetConfig: DataFacetConfig) => ({
        key: facetConfig.id,
        label: facetConfig.label,
        value: facetConfig.id,
        content: perTabContent,
        dividerBefore: facetConfig.dividerBefore,
      })),
    [availableFacets, perTabContent],
  );

  // === SECTION 7: Render ===
  return (
    <div className="mainView">
      <TabView tabs={tabs} route={ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${getCurrentDataFacet()}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      <Debugger
        rowActions={[]}
        headerActions={[]}
        count={++renderCnt.current}
      />
    </div>
  );
};

// EXISTING_CODE
