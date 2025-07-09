// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetNamesPage, NamesCrud, Reload } from '@app';
import { Action, BaseTab, ConfirmModal, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { ActionType } from '@hooks';
import {
  DataFacetConfig,
  toPageDataProp,
  useActionMsgs,
  useActions,
  useColumns,
  useSilencedDialog,
} from '@hooks';
// prettier-ignore
import { useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { Group } from '@mantine/core';
import { useHotkeys } from '@mantine/hooks';
import { msgs, names, types } from '@models';
import { ActionDebugger, useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, namesFacets } from './facets';

// === END SECTION 1 ===

export const Names = () => {
  // === SECTION 2.2: Hook Initialization ===
  const createPayload = usePayload();

  const activeFacetHook = useActiveFacet({
    facets: namesFacets,
    defaultFacet: DEFAULT_FACET,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<names.NamesPage | null>(null);
  const viewStateKey = useMemo(
    (): ViewStateKey => ({
      viewName: ROUTE,
      tabName: getCurrentDataFacet(),
    }),
    [getCurrentDataFacet],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems, goToPage } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);
  // === SECTION 2.1: Modal State ===
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
  // === END SECTION 2.1 ===

  // === END SECTION 2.2 ===

  // === SECTION 3: Data Fetching Logic ===
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
  // === END SECTION 4 ===

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
  const postFunc = useCallback((item: types.Name): types.Name => {
    // EXISTING_CODE
    item = types.Name.createFrom({
      ...item,
      source: item.source || 'TrueBlocks',
    });
    // EXISTING_CODE
    return item;
  }, []);

  const enabledActions = useMemo(() => {
    // EXISTING_CODE
    const currentFacet = getCurrentDataFacet();
    if (currentFacet === types.DataFacet.CUSTOM) {
      return [
        'add',
        'publish',
        'pin',
        'delete',
        'remove',
        'autoname',
        'update',
      ] as ActionType[];
    }
    if (currentFacet === types.DataFacet.BADDRESS) {
      return ['add'] as ActionType[];
    }
    return ['add', 'autoname', 'update'] as ActionType[];
    // EXISTING_CODE
  }, [getCurrentDataFacet]);

  // prettier-ignore
  const { handlers, config } = useActions({
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
    postFunc,
    pageClass: names.NamesPage,
    updateItem: types.Name.createFrom({}),
    createPayload,
    getCurrentDataFacet,
  });

  const {
    handleAutoname: originalHandleAutoname,
    handleRemove,
    handleToggle,
    handleUpdate,
  } = handlers;

  // EXISTING_CODE
  const handleAutoname = useCallback(
    (address: string) => {
      const currentFacet = getCurrentDataFacet();
      if (currentFacet === types.DataFacet.CUSTOM || isSilenced) {
        originalHandleAutoname(address);
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
          emitSuccess('autoname', address);
          activeFacetHook.setActiveFacet(types.DataFacet.CUSTOM);
        },
      });
    },
    [
      getCurrentDataFacet,
      isSilenced,
      originalHandleAutoname,
      emitSuccess,
      activeFacetHook,
    ],
  );

  const handleCloseModal = useCallback(() => {
    setConfirmModal((prev) => ({ ...prev, opened: false }));
  }, []);
  // EXISTING_CODE

  const headerActions = useMemo(() => {
    if (config.headerActions.length === 0) return null;
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
              size="sm"
              isSubdued={action.requiresWallet && !config.isWalletConnected}
            />
          );
        })}
      </Group>
    );
  }, [config.headerActions, config.isWalletConnected, handlers]);
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  const showActions = useMemo(() => {
    return (
      enabledActions.includes('delete' as ActionType) ||
      enabledActions.includes('remove' as ActionType) ||
      enabledActions.includes('autoname' as ActionType)
    );
  }, [enabledActions]);

  const getCanRemove = useCallback(
    (row: unknown): boolean => {
      return (
        Boolean((row as unknown as types.Name)?.deleted) &&
        enabledActions.includes('remove' as ActionType)
      );
    },
    [enabledActions],
  );

  const currentColumns = useColumns(
    getColumns(getCurrentDataFacet()),
    {
      showActions,
      actions: ['delete', 'remove', 'autoname'],
      getCanRemove,
    },
    {
      handleAutoname,
      handleRemove,
      handleToggle,
    },
    toPageDataProp(pageData),
    config,
    true /* perRowCrud */,
  );
  // === END SECTION 7 ===

  // === SECTION 8: Tab Configuration ===
  const canUpdate =
    enabledActions.includes('update' as ActionType) ||
    enabledActions.includes('add' as ActionType);

  const perTabContent = useMemo(() => {
    const actionDebugger = (
      <ActionDebugger
        enabledActions={enabledActions}
        setActiveFacet={activeFacetHook.setActiveFacet}
      />
    );
    return (
      <BaseTab<Record<string, unknown>>
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        viewStateKey={viewStateKey}
        onSubmit={canUpdate ? handleUpdate : undefined}
        debugComponent={actionDebugger}
        headerActions={headerActions}
        onDelete={
          enabledActions.includes('delete' as ActionType)
            ? (rowData: Record<string, unknown>) => {
                const address = String(rowData.address || '');
                handleToggle(address);
              }
            : undefined
        }
        onRemove={
          enabledActions.includes('remove' as ActionType)
            ? (rowData: Record<string, unknown>) => {
                const address = String(rowData.address || '');
                handleRemove(address);
              }
            : undefined
        }
        onAutoname={
          enabledActions.includes('autoname' as ActionType)
            ? (rowData: Record<string, unknown>) => {
                const address = String(rowData.address || '');
                handleAutoname(address);
              }
            : undefined
        }
      />
    );
  }, [
    currentData,
    currentColumns,
    pageData?.isFetching,
    error,
    handleUpdate,
    canUpdate,
    viewStateKey,
    headerActions,
    handleToggle,
    handleRemove,
    handleAutoname,
    enabledActions,
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
      <ConfirmModal
        opened={confirmModal.opened}
        onClose={handleCloseModal}
        onConfirm={confirmModal.onConfirm}
        title={confirmModal.title}
        message={confirmModal.message}
        dialogKey="confirmNamesModal"
      />
    </div>
  );
  // === END SECTION 9 ===
};

// EXISTING_CODE
// EXISTING_CODE
