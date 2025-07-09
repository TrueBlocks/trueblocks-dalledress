// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
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
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import {
  ActionConfig,
  DataFacetConfig,
  toPageDataProp,
  useColumns,
} from '@hooks';
// prettier-ignore
import { useActiveFacet, useEvent, usePayload } from '@hooks';
import { FormView, TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { chunks, msgs, types } from '@models';
import { ActionDebugger, useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, chunksFacets } from './facets';

// === END SECTION 1 ===

export const Chunks = () => {
  // === SECTION 2.2: Hook Initialization ===
  const createPayload = usePayload();

  const activeFacetHook = useActiveFacet({
    facets: chunksFacets,
    defaultFacet: DEFAULT_FACET,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<chunks.ChunksPage | null>(null);
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

  // === END SECTION 2.2 ===

  // === SECTION 3: Data Fetching Logic ===
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
  // === END SECTION 4 ===

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

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleReload = useCallback(async () => {
    try {
      Reload(createPayload(getCurrentDataFacet())).then(() => {
        // The data will reload when the DataLoaded event is fired.
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, createPayload, handleError]);

  useHotkeys([['mod+r', handleReload]]);
  // === END SECTION 4 ===

  // === SECTION 6: Actions ===
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  const showActions = false;
  const getCanRemove = useCallback((_row: unknown): boolean => {
    return false;
  }, []);

  const currentColumns = useColumns(
    getColumns(getCurrentDataFacet()),
    {
      showActions,
      actions: [],
      getCanRemove,
    },
    {},
    toPageDataProp(pageData),
    {} as ActionConfig,
    true /* perRowCrud */,
  );
  // === END SECTION 7 ===

  // === SECTION 8: Tab Configuration ===
  const isForm = useCallback((facet: types.DataFacet) => {
    switch (facet) {
      case types.DataFacet.MANIFEST:
        return true;
      default:
        return false;
    }
  }, []);

  const perTabContent = useMemo(() => {
    const actionDebugger = (
      <ActionDebugger
        enabledActions={[]}
        setActiveFacet={activeFacetHook.setActiveFacet}
      />
    );
    const facet = getCurrentDataFacet();
    if (isForm(facet)) {
      const chunksData = currentData[0] as unknown as Record<string, unknown>;
      if (!chunksData) {
        return <div>No chunks data available</div>;
      }
      const fieldsWithValues = getColumns(getCurrentDataFacet()).map(
        (field) => ({
          ...field,
          value:
            (chunksData?.[field.name as string] as
              | string
              | number
              | boolean
              | undefined) || field.value,
          readOnly: true,
        }),
      );
      return (
        <FormView
          title="Chunks Information"
          formFields={fieldsWithValues}
          onSubmit={() => {}}
        />
      );
    } else {
      return (
        <BaseTab<Record<string, unknown>>
          data={currentData as unknown as Record<string, unknown>[]}
          columns={currentColumns}
          loading={!!pageData?.isFetching}
          error={error}
          debugComponent={actionDebugger}
          headerActions={[]}
          viewStateKey={viewStateKey}
        />
      );
    }
  }, [
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    viewStateKey,
    isForm,
    getCurrentDataFacet,
    activeFacetHook.setActiveFacet,
  ]);

  const tabs = useMemo(
    () =>
      availableFacets.map((facetConfig: DataFacetConfig) => ({
        label: facetConfig.label,
        value: facetConfig.id,
        content: perTabContent,
        dividerBefore: facetConfig.dividerBefore,
      })),
    [availableFacets, perTabContent],
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
