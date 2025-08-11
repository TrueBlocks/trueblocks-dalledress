// Copyright 2016, 2026 The Authors. All rights reserved.
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
import { useFiltering, useSorting } from '@contexts';
import {
  DataFacetConfig,
  toPageDataProp,
  useActiveFacet,
  useColumns,
  useEvent,
  usePayload,
  useViewConfig,
} from '@hooks';
import { TabView } from '@layout';
import { FormView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { status } from '@models';
import { msgs, project, types } from '@models';
import { Debugger, useErrorHandler } from '@utils';
import { toProperCase } from 'src/utils/toProper';

import { createDetailPanelFromViewConfig } from '../utils/detailPanel';

const ROUTE = 'status';

export const Status = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();

  // === SECTION 2.5: Initial ViewConfig Load ===
  const { config: viewConfig } = useViewConfig({
    viewName: 'status',
  });

  // Generate facets from ViewConfig
  const statusFacets: DataFacetConfig[] = useMemo(() => {
    if (!viewConfig?.facets) {
      // Fallback to default facets if ViewConfig not loaded yet
      return [
        {
          id: types.DataFacet.STATUS,
          label: toProperCase(types.DataFacet.STATUS),
        },
        {
          id: types.DataFacet.CACHES,
          label: toProperCase(types.DataFacet.CACHES),
        },
        {
          id: types.DataFacet.CHAINS,
          label: toProperCase(types.DataFacet.CHAINS),
        },
      ];
    }
    return Object.keys(viewConfig.facets).map((facetKey) => ({
      id: facetKey as types.DataFacet,
      label: toProperCase(facetKey),
    }));
  }, [viewConfig]);

  const activeFacetHook = useActiveFacet({
    facets: statusFacets,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<status.StatusPage | null>(null);
  const viewStateKey = useMemo(
    (): project.ViewStateKey => ({
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

  // Listen for active address/chain/period changes to refresh data
  useEvent(msgs.EventType.ADDRESS_CHANGED, fetchData);
  useEvent(msgs.EventType.CHAIN_CHANGED, fetchData);
  useEvent(msgs.EventType.PERIOD_CHANGED, fetchData);

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

  // === SECTION 5: CRUD Operations ===
  // No CRUD operations for Status view

  // === SECTION 6: UI Configuration ===
  const currentColumns = useColumns(
    viewConfig?.facets?.[getCurrentDataFacet()]?.columns || [],
    {
      showActions: false,
      actions: [],
      getCanRemove: useCallback((_row: unknown) => false, []),
    },
    {},
    toPageDataProp(pageData),
    { rowActions: [] },
  );

  // Create detail panel from ViewConfig
  const detailPanel = useMemo(() => {
    return createDetailPanelFromViewConfig(
      viewConfig,
      getCurrentDataFacet,
      'Status Details',
    );
  }, [viewConfig, getCurrentDataFacet]);

  const perTabContent = useMemo(() => {
    const currentFacet = getCurrentDataFacet();
    const facetConfig = viewConfig?.facets?.[currentFacet];

    // Check if this facet should be displayed as a form
    if (facetConfig?.isForm && currentData.length > 0) {
      const statusData = currentData[0] as types.Status;

      // Convert columns to form fields with values
      const fieldsWithValues = currentColumns.map((column) => {
        const rawValue = statusData[
          column.name as keyof types.Status
        ] as unknown;
        let displayValue = '';

        // Convert complex types to strings for display
        if (typeof rawValue === 'string') {
          displayValue = rawValue;
        } else if (typeof rawValue === 'boolean') {
          displayValue = rawValue ? 'Yes' : 'No';
        } else if (typeof rawValue === 'number') {
          displayValue = rawValue.toString();
        } else if (rawValue && typeof rawValue === 'object') {
          displayValue = JSON.stringify(rawValue);
        } else if (rawValue !== null && rawValue !== undefined) {
          displayValue = String(rawValue);
        }

        return {
          ...column,
          value: displayValue,
        };
      });

      return (
        <FormView
          formFields={fieldsWithValues}
          title="System Status"
          onSubmit={() => {}} // No submit functionality for status
        />
      );
    }

    // Default table view
    return (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        headerActions={[]}
        detailPanel={detailPanel}
      />
    );
  }, [
    getCurrentDataFacet,
    viewConfig,
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    viewStateKey,
    detailPanel,
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

// EXISTING_CODE
