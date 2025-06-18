import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetChunksPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { TableKey, useFiltering, useSorting } from '@contexts';
import { useActiveProject, useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { chunks, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumns } from './columns';

const CHUNKS_ROUTE = '/chunks';
const CHUNKS_DEFAULT_LIST = types.ListKind.STATS;

export const Chunks = () => {
  const { lastTab } = useActiveProject();
  const [pageData, setPageData] = useState<chunks.ChunksPage | null>(null);
  const [_state, setState] = useState<types.LoadState>();

  const [listKind, setListKind] = useState<types.ListKind>(
    lastTab[CHUNKS_ROUTE] || CHUNKS_DEFAULT_LIST,
  );
  const tableKey = useMemo(
    (): TableKey => ({ viewName: CHUNKS_ROUTE, tabName: listKind }),
    [listKind],
  );

  const { handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter } = useFiltering(tableKey);

  const listKindRef = useRef(listKind);

  useEffect(() => {
    listKindRef.current = listKind;
  }, [listKind]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetChunksPage(
        listKindRef.current,
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

    switch (listKind) {
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
  }, [pageData, listKind]);

  useEffect(() => {
    const currentTab = lastTab[CHUNKS_ROUTE];
    if (currentTab && currentTab !== listKind) {
      setListKind(currentTab);
    }
  }, [lastTab, listKind]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleReload = useCallback(async () => {
    try {
      await Reload(listKind);
      await fetchData();
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${listKind}`);
    }
  }, [listKind, fetchData, handleError]);

  useEvent(msgs.EventType.DATA_LOADED, () => {
    fetchData();
  });

  useHotkeys([['mod+r', handleReload]]);

  const columns = useMemo(() => getColumns(listKind), [listKind]);

  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={columns}
        loading={!!pageData?.isFetching}
        error={null} // Chunks are read-only, no edit errors
        onSubmit={() => {}} // No-op since chunks are read-only
        tableKey={tableKey}
      />
    ),
    [currentData, columns, pageData?.isFetching, tableKey],
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
      <TabView tabs={tabs} route={CHUNKS_ROUTE} />
    </div>
  );
};
