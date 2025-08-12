// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetChunksPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { Action, ConfirmModal, ExportFormatModal } from '@components';
import { useFiltering, useSorting } from '@contexts';
import {
  DataFacetConfig,
  buildFacetConfigs,
  useActions,
  useActiveFacet,
  useEvent,
  useFacetColumns,
  useFacetForm,
  usePayload,
  useViewConfig,
} from '@hooks';
import { TabView } from '@layout';
import { Group } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';
import { chunks } from '@models';
import { crud, msgs, project, types } from '@models';
import { Debugger, LogError, useErrorHandler } from '@utils';

import { ViewRoute, assertRouteConsistency } from '../routes';
import { createDetailPanelFromViewConfig } from '../utils/detailPanel';

const ROUTE: ViewRoute = 'chunks';

export const Chunks = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();

  // === SECTION 2.5: Initial ViewConfig Load ===
  const { config: viewConfig } = useViewConfig({ viewName: ROUTE });
  assertRouteConsistency(ROUTE, viewConfig);

  // Convert ViewConfig to DataFacetConfig format for useActiveFacet
  const facetsFromConfig = useMemo(
    () => buildFacetConfigs(viewConfig),
    [viewConfig],
  );

  const activeFacetHook = useActiveFacet({
    facets: facetsFromConfig,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<chunks.ChunksPage | null>(null);
  const viewStateKey = useMemo(
    (): project.ViewStateKey => ({
      viewName: ROUTE,
      facetName: getCurrentDataFacet(),
    }),
    [getCurrentDataFacet],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems, goToPage } = usePagination(viewStateKey);
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
        LogError('[CHUNKS] unexpected facet=' + String(facet));
        return [];
    }
  }, [pageData, getCurrentDataFacet]);

  // === SECTION 4: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === ROUTE) {
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

  // === SECTION 5: Actions (standardized placeholder) ===
  const { handlers, config, exportFormatModal, confirmModal } = useActions({
    collection: ROUTE,
    viewStateKey,
    pagination,
    goToPage,
    sort,
    filter,
    viewConfig,
    pageData,
    setPageData,
    setTotalItems,
    crudFunc: async (
      _payload: types.Payload,
      _op: crud.Operation,
      _item:
        | types.ChunkStats
        | types.ChunkIndex
        | types.ChunkBloom
        | types.Manifest,
    ) => {},
    pageFunc: GetChunksPage,
    pageClass: chunks.ChunksPage,
    updateItem: types.ChunkStats.createFrom({}) as unknown as types.ChunkStats,
    createPayload,
    getCurrentDataFacet,
  });

  // header actions are built lazily inside perTabContent

  // === SECTION 6: UI Configuration ===
  const currentColumns = useFacetColumns(
    viewConfig,
    getCurrentDataFacet,
    {
      showActions: false,
      actions: [],
      getCanRemove: useCallback((_row: unknown) => false, []),
    },
    {},
    pageData,
    { rowActions: [] },
  );

  const detailPanel = useMemo(
    () =>
      createDetailPanelFromViewConfig(
        viewConfig,
        getCurrentDataFacet,
        'Chunks Details',
      ),
    [viewConfig, getCurrentDataFacet],
  );

  const { isForm, node: formNode } = useFacetForm<Record<string, unknown>>({
    viewConfig,
    getCurrentDataFacet,
    currentData: currentData as unknown as Record<string, unknown>[],
    currentColumns:
      currentColumns as unknown as import('@components').FormField<
        Record<string, unknown>
      >[],
  });

  const perTabContent = useMemo(() => {
    if (isForm && formNode) return formNode;

    const headerActions = config.headerActions.length ? (
      <Group gap="xs" style={{ flexShrink: 0 }}>
        {config.headerActions.map((action) => {
          const handlerKey = `handle${
            action.type.charAt(0).toUpperCase() + action.type.slice(1)
          }` as keyof typeof handlers;
          const handler = handlers[handlerKey] as () => void;
          return (
            <Action
              key={action.type}
              icon={
                action.icon as keyof ReturnType<
                  typeof import('@hooks').useIconSets
                >
              }
              onClick={handler}
              title={
                action.requiresWallet && !config.isWalletConnected
                  ? `${action.title} (requires wallet connection)`
                  : action.title
              }
              hotkey={action.type === 'export' ? 'mod+x' : undefined}
              size="sm"
            />
          );
        })}
      </Group>
    ) : null;
    return (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        headerActions={headerActions}
        detailPanel={detailPanel}
      />
    );
  }, [
    isForm,
    formNode,
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    viewStateKey,
    detailPanel,
    config.headerActions,
    config.isWalletConnected,
    handlers,
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
        rowActions={config.rowActions}
        headerActions={config.headerActions}
        count={++renderCnt.current}
      />
      <ConfirmModal
        opened={confirmModal.opened}
        onClose={confirmModal.onClose}
        onConfirm={confirmModal.onConfirm}
        title={confirmModal.title}
        message={confirmModal.message}
        dialogKey={confirmModal.dialogKey}
      />
      <ExportFormatModal
        opened={exportFormatModal.opened}
        onClose={exportFormatModal.onClose}
        onFormatSelected={exportFormatModal.onFormatSelected}
      />
    </div>
  );
};

// EXISTING_CODE
// EXISTING_CODE
