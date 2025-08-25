// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetDalleDressPage, Reload } from '@app';
import { SeriesCrud } from '@app';
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
import { dalle, dalledress } from '@models';
import { crud, msgs, project, types } from '@models';
import { Debugger, LogError, useErrorHandler } from '@utils';

import { ViewRoute, assertRouteConsistency } from '../routes';
import { createDetailPanel } from '../utils/detailPanel';
import { SeriesModal } from './components';
import { renderers } from './components';
import { useSeriesModal } from './hooks/seriesModal';

const ROUTE: ViewRoute = 'dalledress';
export const DalleDress = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();
  // === SECTION 2.5: Initial ViewConfig Load ===
  const { config: viewConfig } = useViewConfig({ viewName: ROUTE });
  assertRouteConsistency(ROUTE, viewConfig);

  const facetsFromConfig: DataFacetConfig[] = useMemo(
    () => buildFacetConfigs(viewConfig),
    [viewConfig],
  );

  const activeFacetHook = useActiveFacet({
    facets: facetsFromConfig,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<dalledress.DalleDressPage | null>(
    null,
  );
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
      const result = await GetDalleDressPage(
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
      case types.DataFacet.GENERATOR:
        return [];
      case types.DataFacet.SERIES:
        return pageData.series || [];
      case types.DataFacet.DATABASES:
        return pageData.databases || [];
      case types.DataFacet.EVENTS:
        return pageData.logs || [];
      case types.DataFacet.GALLERY:
        return pageData.logs || [];
      default:
        LogError('[DalleDress] unexpected facet=' + String(facet));
        return [];
    }
  }, [pageData, getCurrentDataFacet]);

  // === SECTION 4: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === ROUTE) {
        if (payload.dataFacet === getCurrentDataFacet()) {
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
      Reload(createPayload(getCurrentDataFacet())).then(() => {});
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [clearError, getCurrentDataFacet, createPayload, handleError]);

  useHotkeys([['mod+r', handleReload]]);

  // === SECTION 5: Actions ===
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
      payload: types.Payload,
      op: crud.Operation,
      item: unknown,
    ) => {
      if (getCurrentDataFacet() === types.DataFacet.SERIES) {
        await SeriesCrud(payload, op, item as unknown as dalle.Series);
      }
    },
    pageFunc: GetDalleDressPage,
    pageClass: dalledress.DalleDressPage,
    updateItem: types.Log.createFrom({}),
    createPayload,
    getCurrentDataFacet,
  });
  const { seriesModal, closeSeriesModal, submitSeriesModal } = useSeriesModal({
    getCurrentDataFacet,
    createPayload,
    collection: ROUTE,
  });

  const headerActions = useMemo(() => {
    if (!config.headerActions.length) return null;
    return (
      <Group gap="xs" style={{ flexShrink: 0 }}>
        {config.headerActions.map((action) => {
          const handlerKey =
            `handle${action.type.charAt(0).toUpperCase() + action.type.slice(1)}` as keyof typeof handlers;
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
    );
  }, [config.headerActions, config.isWalletConnected, handlers]);

  // === SECTION 6: UI Configuration ===
  const currentColumns = useFacetColumns(
    viewConfig,
    getCurrentDataFacet,
    {
      showActions: useCallback(
        () => getCurrentDataFacet() === types.DataFacet.SERIES,
        [getCurrentDataFacet],
      ),
      actions: config.rowActions.map(
        (a) => a.type as unknown as import('@hooks').ActionType,
      ),
      getCanRemove: useCallback((_row: unknown) => false, []),
    },
    handlers,
    pageData,
    config,
  );

  const detailPanel = useMemo(
    () =>
      createDetailPanel(viewConfig, getCurrentDataFacet, 'Dalledress Details'),
    [viewConfig, getCurrentDataFacet],
  );

  const rendererMap = useMemo(
    () => renderers(pageData, viewStateKey),
    [pageData, viewStateKey],
  );
  const { isForm, node: formNode } = useFacetForm<Record<string, unknown>>({
    viewConfig,
    getCurrentDataFacet,
    currentData: currentData as unknown as Record<string, unknown>[],
    currentColumns:
      currentColumns as unknown as import('@components').FormField<
        Record<string, unknown>
      >[],
    renderers: rendererMap,
  });

  const perTabContent = useMemo(() => {
    const facet = getCurrentDataFacet();
    // Prefer custom renderer for certain facets even if no data yet
    if (
      rendererMap[facet] &&
      (facet === types.DataFacet.GENERATOR || facet === types.DataFacet.GALLERY)
    ) {
      return rendererMap[facet]();
    }
    if (isForm && formNode) return formNode;
    return (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        headerActions={headerActions}
        detailPanel={detailPanel}
        onDelete={(_rowData) => {}}
      />
    );
  }, [
    getCurrentDataFacet,
    rendererMap,
    currentColumns,
    isForm,
    formNode,
    currentData,
    pageData?.isFetching,
    error,
    viewStateKey,
    headerActions,
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
      <SeriesModal
        opened={
          seriesModal.opened && getCurrentDataFacet() === types.DataFacet.SERIES
        }
        mode={seriesModal.mode}
        existing={(pageData?.series || []) as dalle.Series[]}
        initial={seriesModal.initial as dalle.Series | undefined}
        onSubmit={submitSeriesModal}
        onClose={closeSeriesModal}
      />
    </div>
  );
};

// EXISTING_CODE
