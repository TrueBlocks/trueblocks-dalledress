// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetNamesPage, NamesCrud, Reload } from '@app';
import { BaseTab, usePagination } from '@components';
import { Action, ConfirmModal, ExportFormatModal } from '@components';
import { useFiltering, useSorting } from '@contexts';
import {
  DataFacetConfig,
  useActiveFacet,
  useEvent,
  useFacetColumns,
  usePayload,
  useViewConfig,
} from '@hooks';
import { buildFacetConfigs } from '@hooks';
import { useActionMsgs, useActions, useSilencedDialog } from '@hooks';
import { TabView } from '@layout';
import { Group } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';
import { names } from '@models';
import { msgs, project, types } from '@models';
import { Debugger, LogError, useErrorHandler } from '@utils';
import { createDetailPanelFromViewConfig } from '@views';

import { ViewRoute, assertRouteConsistency } from '../routes';

const ROUTE: ViewRoute = 'names';

export const Names = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();

  // ViewConfig hook - guaranteed to be available in Wails
  const { config: viewConfig } = useViewConfig({ viewName: ROUTE });
  assertRouteConsistency(ROUTE, viewConfig);

  // Generate facets from ViewConfig - no fallbacks needed in Wails
  const facetsFromConfig = useMemo(
    () => buildFacetConfigs(viewConfig),
    [viewConfig],
  );

  const activeFacetHook = useActiveFacet({
    facets: facetsFromConfig,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet, setActiveFacet } =
    activeFacetHook;

  const [pageData, setPageData] = useState<names.NamesPage | null>(null);
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

  // Names-specific state
  const { emitSuccess: _emitSuccess } = useActionMsgs(ROUTE);
  const { isSilenced: _isSilenced } = useSilencedDialog('createCustomName');

  // === SECTION 3: Data Fetching ===
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const payload = createPayload(getCurrentDataFacet());
      const first = pagination.currentPage * pagination.pageSize;
      const result = await GetNamesPage(
        payload,
        first,
        pagination.pageSize,
        sort,
        filter,
      );
      setPageData(result);
      setTotalItems(result.totalItems || 0);
    } catch (err: unknown) {
      // LogError('[NAMES] fetchData error ' + String(err));
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
    if (!pageData) return [] as types.Name[];
    const facet = getCurrentDataFacet();
    switch (facet) {
      case types.DataFacet.ALL:
      case types.DataFacet.CUSTOM:
      case types.DataFacet.PREFUND:
      case types.DataFacet.REGULAR:
      case types.DataFacet.BADDRESS:
        return pageData.names || [];
      default:
        LogError('[NAMES] unexpected facet=' + String(facet));
        return [] as types.Name[];
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
      const payload = createPayload(getCurrentDataFacet());
      Reload(payload).then(() => {});
    } catch (err: unknown) {
      // LogError('[NAMES] handleReload error ' + String(err));
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [clearError, getCurrentDataFacet, createPayload, handleError]);

  useHotkeys([['mod+r', handleReload]]);

  // === SECTION 5: CRUD Operations ===
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
    crudFunc: NamesCrud,
    pageFunc: GetNamesPage,
    pageClass: names.NamesPage,
    updateItem: types.Name.createFrom({}),
    postFunc: useCallback((item: types.Name): types.Name => {
      const rawDecimals = (item as unknown as Record<string, unknown>).decimals;
      const rawParts = (item as unknown as Record<string, unknown>).parts;
      const decimals =
        typeof rawDecimals === 'string'
          ? parseInt(rawDecimals || '0', 10) || 0
          : (rawDecimals as number | undefined);
      const parts =
        typeof rawParts === 'string'
          ? parseInt(rawParts || '0', 10) || undefined
          : (rawParts as number | undefined);
      item = types.Name.createFrom({
        ...item,
        decimals,
        parts,
        source: item.source || 'TrueBlocks',
      });
      return item;
    }, []),
    createPayload,
    getCurrentDataFacet,
  });

  const { handleAutoname, handleRemove, handleToggle, handleUpdate } = handlers;

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
      showActions: true,
      actions: ['delete', 'remove', 'update', 'autoname'],
      getCanRemove: useCallback(
        (row: unknown) => Boolean((row as unknown as types.Name)?.deleted),
        [],
      ),
    },
    {
      handleAutoname,
      handleRemove,
      handleToggle,
      handleUpdate,
    },
    pageData,
    config,
  );

  const detailPanel = useMemo(
    () =>
      createDetailPanelFromViewConfig(
        viewConfig,
        getCurrentDataFacet,
        'Names Details',
      ),
    [viewConfig, getCurrentDataFacet],
  );

  const perTabContent = useMemo(() => {
    return (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        headerActions={headerActions}
        detailPanel={detailPanel}
        onDelete={(rowData) => handleToggle(String(rowData.address || ''))}
        onRemove={(rowData) => handleRemove(String(rowData.address || ''))}
        onAutoname={(rowData) => handleAutoname(String(rowData.address || ''))}
        onSubmit={handleUpdate}
      />
    );
  }, [
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    viewStateKey,
    headerActions,
    detailPanel,
    handleToggle,
    handleRemove,
    handleAutoname,
    handleUpdate,
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

  // When autoname confirmed and not already on CUSTOM facet, switch facets (post-confirm behavior)
  // This effect watches for a completed autoname action by inspecting modal state transitions.
  // For simplicity, rely on refresh triggered by DATA_LOADED event after backend processes autoname.

  const handleConfirmModalClose = confirmModal.onClose;

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
        onClose={handleConfirmModalClose}
        onConfirm={() => {
          const wasCustom = getCurrentDataFacet() === types.DataFacet.CUSTOM;
          confirmModal.onConfirm();
          if (!wasCustom) {
            setActiveFacet(types.DataFacet.CUSTOM);
          }
        }}
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
