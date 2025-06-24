import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetChunksPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { useActiveFacet, useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { chunks, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumns } from './columns';
import {
  CHUNKS_DEFAULT_FACET,
  CHUNKS_ROUTE as ROUTE,
  chunksFacets,
} from './facets';

export const Chunks = () => {
  const activeFacetHook = useActiveFacet({
    facets: chunksFacets,
    defaultFacet: CHUNKS_DEFAULT_FACET,
    viewRoute: ROUTE,
  });

  const { getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<chunks.ChunksPage | null>(null);
  const [_state, setState] = useState<types.LoadState>();

  const viewStateKey = useMemo(
    (): ViewStateKey => ({
      viewName: ROUTE,
      tabName: getCurrentDataFacet(),
    }),
    [getCurrentDataFacet],
  );

  const { handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);

  const dataFacetRef = useRef(getCurrentDataFacet());

  useEffect(() => {
    dataFacetRef.current = getCurrentDataFacet();
  }, [getCurrentDataFacet]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetChunksPage(
        types.Payload.createFrom({
          dataFacet: dataFacetRef.current,
        }),
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter ?? '',
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
  ]);

  const currentData = useMemo(() => {
    if (!pageData) return [];

    const currentDataFacet = getCurrentDataFacet();
    switch (currentDataFacet) {
      case types.DataFacet.STATS:
        return pageData.stats || [];
      case types.DataFacet.INDEX:
        return pageData.index || [];
      case types.DataFacet.BLOOMS:
        return pageData.blooms || [];
      case types.DataFacet.MANIFEST:
        return pageData.manifest || [];
      default:
        return pageData.stats || [];
    }
  }, [pageData, getCurrentDataFacet]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleReload = useCallback(async () => {
    try {
      await Reload(getCurrentDataFacet() as types.DataFacet, '', '');
      await fetchData();
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, fetchData, handleError]);

  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'chunks') {
        fetchData();
      }
    },
  );

  useHotkeys([['mod+r', handleReload]]);

  const columns = useMemo(
    () => getColumns(getCurrentDataFacet() as types.DataFacet),
    [getCurrentDataFacet],
  );

  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={columns}
        loading={!!pageData?.isFetching}
        error={null} // Chunks are read-only, no edit errors
        onSubmit={() => {}} // No-op since chunks are read-only
        viewStateKey={viewStateKey}
      />
    ),
    [currentData, columns, pageData?.isFetching, viewStateKey],
  );

  const tabs = useMemo(
    () => [
      {
        label: 'Stats',
        value: types.DataFacet.STATS,
        content: perTabTable,
      },
      {
        label: 'Index',
        value: types.DataFacet.INDEX,
        content: perTabTable,
      },
      {
        label: 'Blooms',
        value: types.DataFacet.BLOOMS,
        content: perTabTable,
      },
      {
        label: 'Manifest',
        value: types.DataFacet.MANIFEST,
        content: perTabTable,
      },
    ],
    [perTabTable],
  );

  return (
    <div className="mainView">
      <TabView tabs={tabs} route={ROUTE} />
    </div>
  );
};
