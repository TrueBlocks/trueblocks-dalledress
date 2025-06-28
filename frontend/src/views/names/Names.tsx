// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetNamesPage, NamesCrud, Reload } from '@app';
import { Action, Chips, mapNameToChips } from '@components';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
// prettier-ignore
import { ActionData, useActionConfig, useCrudOperations } from '@hooks';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { msgs, names, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, namesFacets } from './facets';

// === END SECTION 1 ===

export const Names = () => {
  // === SECTION 2: Hook Initialization ===
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
      const result = await GetNamesPage(
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

  // === SECTION 5: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'names') {
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
        // The data will reload when the DataLoaded event is fired.
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, createPayload, handleError]);

  useHotkeys([['mod+r', handleReload]]);
  // === END SECTION 5 ===

  // === SECTION 6: CRUD Operations ===
  const actionConfig = useActionConfig({
    operations: ['delete', 'undelete', 'remove', 'autoname'],
  });

  // EXISTING_CODE
  const postFunc = useCallback((item: types.Name): types.Name => {
    return types.Name.createFrom({
      ...item,
      source: item.source || 'TrueBlocks',
    });
  }, []);
  // EXISTING_CODE

  // prettier-ignore
  const { handleRemove, handleToggle, handleAutoname, handleUpdate } = useCrudOperations({
    collectionName: 'names',
    crudFunc: NamesCrud,
    pageFunc: GetNamesPage,
    postFunc: postFunc,
    pageClass: names.NamesPage,
    updateItem: types.Name.createFrom({}),
    getCurrentDataFacet,
    pageData,
    setPageData,
    setTotalItems,
    dataFacetRef,
    actionConfig,
  });
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE
  const { setFiltering } = useFiltering(viewStateKey);

  const currentColumns = useMemo(() => {
    const handleChipClick = (chip: string) => {
      setFiltering(chip);
      Reload(createPayload(dataFacetRef.current)).then(() => {
        fetchData();
      });
    };

    const baseColumns = getColumns(getCurrentDataFacet()).map((col) =>
      col.key === 'chips'
        ? {
            ...col,
            render: (row: Record<string, unknown>) => {
              const nameObject = row as unknown as types.Name;
              const chipItems = mapNameToChips(nameObject);
              return <Chips items={chipItems} onChipClick={handleChipClick} />;
            },
          }
        : col,
    );

    const renderActions = (actionData: ActionData) => {
      const isDeleted = actionData.isDeleted;
      const effectiveDeletedState = actionData.isProcessing
        ? !isDeleted
        : isDeleted;

      return (
        <div className="action-buttons-container">
          <Action
            icon="Delete"
            iconOff="Undelete"
            isOn={!effectiveDeletedState}
            onClick={() => handleToggle(actionData.addressStr)}
            disabled={actionData.isProcessing}
            title={effectiveDeletedState ? 'Undelete' : 'Delete'}
            size="sm"
          />
          <Action
            icon="Remove"
            onClick={() => handleRemove(actionData.addressStr)}
            disabled={actionData.isProcessing || !effectiveDeletedState}
            title="Remove"
            size="sm"
          />
          <Action
            icon="Autoname"
            onClick={() => handleAutoname(actionData.addressStr)}
            disabled={actionData.isProcessing}
            title="Auto-generate name"
            size="sm"
          />
        </div>
      );
    };

    const getCanRemove = (row: Record<string, unknown>) => {
      const name = row as unknown as types.Name;
      return Boolean(name.deleted);
    };

    return actionConfig.injectActionColumn(
      baseColumns,
      renderActions,
      getCanRemove,
    );
  }, [
    setFiltering,
    fetchData,
    getCurrentDataFacet,
    actionConfig,
    handleToggle,
    handleRemove,
    handleAutoname,
    createPayload,
  ]);
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
        onSubmit={handleUpdate}
        viewStateKey={viewStateKey}
      />
    ),
    // prettier-ignore
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      handleUpdate,
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

// EXISTING_CODE
// EXISTING_CODE
