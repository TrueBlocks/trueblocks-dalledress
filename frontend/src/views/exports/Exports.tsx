// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetExportsPage, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { Action, ExportFormatModal } from '@components';
import { useFiltering, useSorting } from '@contexts';
import {
  ActionType,
  DataFacetConfig,
  toPageDataProp,
  useActions,
  useActiveFacet,
  useColumns,
  useEvent,
  usePayload,
  useViewConfig,
} from '@hooks';
import { TabView } from '@layout';
import { Group } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';
import { exports } from '@models';
import { msgs, project, types } from '@models';
import { Debugger, useErrorHandler } from '@utils';
import { createDetailPanelFromViewConfig } from '@views';

export const Exports = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();

  // ViewConfig hook - guaranteed to be available in Wails
  const { config: viewConfig } = useViewConfig({
    viewName: 'exports',
  });

  // Generate facets from ViewConfig - no fallbacks needed in Wails
  const facetsFromConfig = useMemo((): DataFacetConfig[] => {
    // Define the correct order for Exports facets
    const facetOrder = [
      'statements',
      'balances',
      'transfers',
      'transactions',
      'withdrawals',
      'assets',
      'logs',
      'traces',
      'receipts',
    ];

    return facetOrder
      .filter((facetId) => viewConfig.facets[facetId]) // Only include facets that exist
      .map((facetId) => {
        const facetConfig = viewConfig.facets[facetId];
        return {
          id: facetId as types.DataFacet,
          label: facetConfig?.name || facetId,
          dividerBefore: facetConfig?.dividerBefore,
        };
      });
  }, [viewConfig.facets]);

  const activeFacetHook = useActiveFacet({
    facets: facetsFromConfig,
    viewRoute: 'exports',
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<exports.ExportsPage | null>(null);
  const viewStateKey = useMemo(
    (): project.ViewStateKey => ({
      viewName: 'exports',
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
    const currentFacet = getCurrentDataFacet();

    try {
      const result = await GetExportsPage(
        createPayload(currentFacet),
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

  // === SECTION 4: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'exports') {
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

  // === SECTION 5: Actions Configuration ===
  const enabledActions = useMemo(() => {
    // Auto-enable Export for all DataTable views that are not forms
    return ['export'] as ActionType[];
  }, []);

  const { handlers, config, exportFormatModal } = useActions({
    collection: 'exports',
    viewStateKey,
    pagination,
    goToPage: () => {}, // Exports typically don't have pagination navigation
    sort,
    filter,
    enabledActions,
    pageData,
    setPageData,
    setTotalItems,
    crudFunc: () => Promise.resolve(), // Exports don't have CRUD operations
    pageFunc: GetExportsPage,
    pageClass: exports.ExportsPage,
    updateItem: undefined,
    createPayload,
    getCurrentDataFacet,
  });

  const headerActions = useMemo(() => {
    if (!config.headerActions?.length) return null;
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
  const currentData = useMemo(() => {
    if (!pageData) return [];
    const facet = getCurrentDataFacet();
    switch (facet) {
      case types.DataFacet.STATEMENTS:
        return pageData.statements || [];
      case types.DataFacet.BALANCES:
        return pageData.balances || [];
      case types.DataFacet.TRANSFERS:
        return pageData.transfers || [];
      case types.DataFacet.TRANSACTIONS:
        return pageData.transactions || [];
      case types.DataFacet.WITHDRAWALS:
        return pageData.withdrawals || [];
      case types.DataFacet.ASSETS:
        return pageData.assets || [];
      case types.DataFacet.LOGS:
        return pageData.logs || [];
      case types.DataFacet.TRACES:
        return pageData.traces || [];
      case types.DataFacet.RECEIPTS:
        return pageData.receipts || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);

  const currentColumns = useColumns(
    viewConfig.facets[getCurrentDataFacet()]?.columns || [],
    {
      showActions: false,
      actions: [],
      getCanRemove: useCallback((_row: unknown) => false, []),
    },
    {},
    toPageDataProp(pageData),
    { rowActions: [] },
  );

  const perTabContent = useMemo(
    () => (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        headerActions={headerActions}
        detailPanel={createDetailPanelFromViewConfig(
          viewConfig,
          getCurrentDataFacet,
          'Exports Details',
        )}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      viewStateKey,
      headerActions,
      getCurrentDataFacet,
      viewConfig,
    ],
  );

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
      <TabView tabs={tabs} route="exports" />
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
      <ExportFormatModal
        opened={exportFormatModal.opened}
        onClose={exportFormatModal.onClose}
        onFormatSelected={exportFormatModal.onFormatSelected}
      />
    </div>
  );
};
