import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetChunksPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { useActiveFacet, useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { chunks, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import {
  CHUNKS_DEFAULT_FACET,
  CHUNKS_ROUTE as ROUTE,
  chunksFacets,
} from './chunksFacets';
import { getColumns } from './columns';

export const Chunks = () => {
  const activeFacetHook = useActiveFacet({
    facets: chunksFacets,
    defaultFacet: CHUNKS_DEFAULT_FACET,
    viewRoute: ROUTE,
  });

  const { getCurrentListKind } = activeFacetHook;

  const [pageData, setPageData] = useState<chunks.ChunksPage | null>(null);
  const [_state, setState] = useState<types.LoadState>();

  const viewStateKey = useMemo(
    (): ViewStateKey => ({
      viewName: ROUTE,
      tabName: getCurrentListKind(),
    }),
    [getCurrentListKind],
  );

  const { handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);

  const listKindRef = useRef(getCurrentListKind());

  useEffect(() => {
    listKindRef.current = getCurrentListKind();
  }, [getCurrentListKind]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetChunksPage(
        listKindRef.current as types.ListKind,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter ?? '',
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
  ]);

  const currentData = useMemo(() => {
    if (!pageData) return [];

    const currentListKind = getCurrentListKind();
    switch (currentListKind) {
      case types.ListKind.STATS:
        return pageData.chunksStats || [];
      case types.ListKind.INDEX:
        return pageData.chunksIndex || [];
      case types.ListKind.BLOOMS:
        return pageData.chunksBlooms || [];
      case types.ListKind.MANIFEST:
        return pageData.chunksManifest || [];
      default:
        return pageData.chunksStats || [];
    }
  }, [pageData, getCurrentListKind]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleReload = useCallback(async () => {
    try {
      await Reload(getCurrentListKind() as types.ListKind);
      await fetchData();
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentListKind()}`);
    }
  }, [getCurrentListKind, fetchData, handleError]);

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
    () => getColumns(getCurrentListKind() as types.ListKind),
    [getCurrentListKind],
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
        value: types.ListKind.STATS,
        content: perTabTable,
      },
      {
        label: 'Index',
        value: types.ListKind.INDEX,
        content: perTabTable,
      },
      {
        label: 'Blooms',
        value: types.ListKind.BLOOMS,
        content: perTabTable,
      },
      {
        label: 'Manifest',
        value: types.ListKind.MANIFEST,
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
