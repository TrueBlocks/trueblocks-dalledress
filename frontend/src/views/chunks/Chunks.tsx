// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetChunksPage, Reload } from '@app';
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
import { FormView, TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { chunks } from '@models';
import { msgs, project, types } from '@models';
import { Debugger, useErrorHandler } from '@utils';

import { createDetailPanelFromViewConfig } from '../utils/detailPanel';

export const Chunks = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();

  // === SECTION 2.5: Initial ViewConfig Load ===
  const { config: viewConfig } = useViewConfig({
    viewName: 'chunks',
  });

  // Convert ViewConfig to DataFacetConfig format for useActiveFacet
  const facetsFromConfig = useMemo((): DataFacetConfig[] => {
    if (!viewConfig?.facets) {
      // Fallback to default facets if ViewConfig not loaded yet
      return [
        {
          id: types.DataFacet.STATS,
          label: 'Stats',
        },
        {
          id: types.DataFacet.INDEX,
          label: 'Index',
        },
        {
          id: types.DataFacet.BLOOMS,
          label: 'Blooms',
        },
        {
          id: types.DataFacet.MANIFEST,
          label: 'Manifest',
        },
      ];
    }
    // Maintain specific order for chunks facets: Stats, Index, Blooms, Manifest
    const orderedKeys = ['stats', 'index', 'blooms', 'manifest'];
    const availableKeys = Object.keys(viewConfig.facets);
    const sortedKeys = orderedKeys.filter((key) => availableKeys.includes(key));

    return sortedKeys.map((facetKey) => ({
      id: facetKey as types.DataFacet,
      label: viewConfig.facets[facetKey]?.name || facetKey,
    }));
  }, [viewConfig?.facets]);

  const activeFacetHook = useActiveFacet({
    facets: facetsFromConfig,
    viewRoute: 'chunks',
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<chunks.ChunksPage | null>(null);
  const viewStateKey = useMemo(
    (): project.ViewStateKey => ({
      viewName: 'chunks',
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
      const result = await GetChunksPage(
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
      case types.DataFacet.STATS:
        return pageData.stats || [];
      case types.DataFacet.INDEX:
        return pageData.index || [];
      case types.DataFacet.BLOOMS:
        return pageData.blooms || [];
      case types.DataFacet.MANIFEST:
        return pageData.manifest || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);

  // === SECTION 4: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'chunks') {
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
  // No CRUD operations for Chunks view

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

  const perTabContent = useMemo(() => {
    const facet = getCurrentDataFacet();
    const currentFacetConfig = viewConfig?.facets?.[facet];

    if (currentFacetConfig?.isForm) {
      // Form view for MANIFEST facet
      const chunksData = currentData[0] as unknown as Record<string, unknown>;
      if (!chunksData) {
        return <div>No chunks data available</div>;
      }
      return (
        <FormView
          title="Chunks Information"
          formFields={[
            {
              key: 'version',
              name: 'version',
              header: 'Version',
              label: 'Version',
              type: 'text',
              width: '100px',
              readOnly: true,
              value: (chunksData?.version as string) ?? '',
            },
            {
              key: 'chain',
              name: 'chain',
              header: 'Chain',
              label: 'Chain',
              type: 'text',
              width: '100px',
              readOnly: true,
              value: (chunksData?.chain as string) ?? '',
            },
            {
              key: 'specification',
              name: 'specification',
              header: 'Specification',
              label: 'Specification',
              type: 'ipfshash',
              width: '100px',
              readOnly: true,
              value: (chunksData?.specification as string) ?? '',
            },
          ]}
          onSubmit={() => {}}
        />
      );
    }

    // Table view for other facets (STATS, INDEX, BLOOMS)
    const detailPanel = createDetailPanelFromViewConfig(
      viewConfig,
      getCurrentDataFacet,
      'Chunks Details',
    );

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
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    viewStateKey,
    getCurrentDataFacet,
    viewConfig,
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
      <TabView tabs={tabs} route="chunks" />
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
