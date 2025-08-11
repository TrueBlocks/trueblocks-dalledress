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
  toPageDataProp,
  useActiveFacet,
  useColumns,
  useEvent,
  usePayload,
  useViewConfig,
} from '@hooks';
import {
  ActionType,
  useActionMsgs,
  useActions,
  useSilencedDialog,
} from '@hooks';
import { TabView } from '@layout';
import { Group } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';
import { names } from '@models';
import { msgs, project, types } from '@models';
import { Debugger, useErrorHandler } from '@utils';
import { createDetailPanelFromViewConfig } from '@views';

const ROUTE = 'names' as const;

export const Names = () => {
  // === SECTION 2: Hook Initialization ===
  const renderCnt = useRef(0);
  const createPayload = usePayload();

  // ViewConfig hook - guaranteed to be available in Wails
  const { config: viewConfig } = useViewConfig({
    viewName: 'names',
  });

  // Generate facets from ViewConfig - no fallbacks needed in Wails
  const facetsFromConfig = useMemo((): DataFacetConfig[] => {
    // Define the correct order for Names facets
    const facetOrder = ['all', 'custom', 'prefund', 'regular', 'baddress'];

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
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

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
  const { emitSuccess } = useActionMsgs('names');
  const [confirmModal, setConfirmModal] = useState<{
    opened: boolean;
    address: string;
    title: string;
    message: string;
    onConfirm: () => void;
  }>({
    opened: false,
    address: '',
    title: '',
    message: '',
    onConfirm: () => {},
  });
  const { isSilenced } = useSilencedDialog('createCustomName');

  // === SECTION 3: Data Fetching ===
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetNamesPage(
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
      case types.DataFacet.ALL:
        return pageData.names || [];
      case types.DataFacet.CUSTOM:
        return pageData.names || [];
      case types.DataFacet.PREFUND:
        return pageData.names || [];
      case types.DataFacet.REGULAR:
        return pageData.names || [];
      case types.DataFacet.BADDRESS:
        return pageData.names || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);

  // === SECTION 4: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'names') {
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
  const enabledActions = useMemo(() => {
    const currentFacet = getCurrentDataFacet();
    let actions: ActionType[] = [];
    if (currentFacet === types.DataFacet.CUSTOM) {
      actions = [
        'add',
        'publish',
        'pin',
        'delete',
        'remove',
        'autoname',
        'update',
      ];
    } else if (currentFacet === types.DataFacet.BADDRESS) {
      actions = ['add'];
    } else {
      actions = ['add', 'autoname', 'update'];
    }
    if (!actions.includes('export')) {
      actions.push('export');
    }
    return actions;
  }, [getCurrentDataFacet]);
  const { handlers, config, exportFormatModal } = useActions({
    collection: 'names',
    viewStateKey,
    pagination,
    goToPage,
    sort,
    filter,
    enabledActions,
    pageData,
    setPageData,
    setTotalItems,
    crudFunc: NamesCrud,
    pageFunc: GetNamesPage,
    pageClass: names.NamesPage,
    updateItem: types.Name.createFrom({}),
    postFunc: useCallback((item: types.Name): types.Name => {
      item = types.Name.createFrom({
        ...item,
        source: item.source || 'TrueBlocks',
      });
      return item;
    }, []),
    createPayload,
    getCurrentDataFacet,
  });

  const {
    handleAutoname: originalHandleAutoname,
    handleRemove,
    handleToggle,
    handleUpdate,
  } = handlers;

  const handleAutoname = useCallback(
    (address: string) => {
      const currentFacet = getCurrentDataFacet();
      if (currentFacet === types.DataFacet.CUSTOM || isSilenced) {
        originalHandleAutoname(address);
        emitSuccess(
          'autoname',
          'Successfully created custom name for ${address}',
        );
        if (currentFacet === types.DataFacet.CUSTOM) {
          fetchData();
        } else {
          activeFacetHook.setActiveFacet(types.DataFacet.CUSTOM);
        }
        return;
      }
      setConfirmModal({
        opened: true,
        address,
        title: 'Create Custom Name',
        message:
          'This will create a custom name for this address. The new custom name will be available in the Custom tab.',
        onConfirm: () => {
          originalHandleAutoname(address);
          emitSuccess(
            'autoname',
            'Successfully created custom name for ${address}',
          );
          activeFacetHook.setActiveFacet(types.DataFacet.CUSTOM);
          setConfirmModal((prev) => ({ ...prev, opened: false }));
        },
      });
    },
    [
      getCurrentDataFacet,
      isSilenced,
      originalHandleAutoname,
      emitSuccess,
      activeFacetHook,
      fetchData,
    ],
  );

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
  const currentFacetConfig = viewConfig.facets[getCurrentDataFacet()];

  const currentColumns = useColumns(
    currentFacetConfig?.columns || [],
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
    },
    toPageDataProp(pageData),
    config,
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
        detailPanel={createDetailPanelFromViewConfig(
          viewConfig,
          getCurrentDataFacet,
          'Names Details',
        )}
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
    getCurrentDataFacet,
    handleToggle,
    handleRemove,
    handleAutoname,
    handleUpdate,
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

  const handleConfirmModalClose = useCallback(
    () => setConfirmModal((prev) => ({ ...prev, opened: false })),
    [],
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
        onClose={handleConfirmModalClose}
        onConfirm={confirmModal.onConfirm}
        title={confirmModal.title}
        message={confirmModal.message}
        dialogKey="confirmNamesModal"
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
