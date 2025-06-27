// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetChunksPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
// prettier-ignore
import { useActionConfig } from '@hooks';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { chunks, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, chunksFacets } from './facets';

// === END SECTION 1 ===

export const Chunks = () => {
  // === SECTION 2: Hook Initialization ===
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
      const result = await GetChunksPage(
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

  // === SECTION 5: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'chunks') {
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
  const actionConfig = useActionConfig({
    operations: [],
  });
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE
  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(getCurrentDataFacet());
    return actionConfig.injectActionColumn(baseColumns, () => null);
  }, [getCurrentDataFacet, actionConfig]);
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
