// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetMonitorsPage, MonitorsCrud, Reload } from '@app';
import { Action } from '@components';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
// prettier-ignore
import { ActionData, useActionConfig, useCrudOperations } from '@hooks';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { monitors, msgs, types } from '@models';
import { useErrorHandler } from '@utils';

import { getColumns } from './columns';
import { DEFAULT_FACET, ROUTE, monitorsFacets } from './facets';

// === END SECTION 1 ===

export const Monitors = () => {
  // === SECTION 2: Hook Initialization ===
  const createPayload = usePayload();

  const activeFacetHook = useActiveFacet({
    facets: monitorsFacets,
    defaultFacet: DEFAULT_FACET,
    viewRoute: ROUTE,
  });
  const { availableFacets, getCurrentDataFacet } = activeFacetHook;

  const [pageData, setPageData] = useState<monitors.MonitorsPage | null>(null);
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
      const result = await GetMonitorsPage(
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
      case types.DataFacet.MONITORS:
        return pageData.monitors || [];
      default:
        return [];
    }
  }, [pageData, getCurrentDataFacet]);
  // === END SECTION 4 ===

  // === SECTION 5: Event Handling ===
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'monitors') {
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
        fetchData();
      });
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, createPayload, fetchData, handleError]);

  useHotkeys([['mod+r', handleReload]]);
  // === END SECTION 5 ===

  // === SECTION 6: CRUD Operations ===
  const actionConfig = useActionConfig({
    operations: ['delete', 'undelete', 'remove'],
  });

  // prettier-ignore
  const { handleDelete, handleUndelete, handleRemove } = useCrudOperations({
    collectionName: 'monitors',
    crudFunc: MonitorsCrud,
    pageFunc: GetMonitorsPage,
    pageClass: monitors.MonitorsPage,
    emptyItem: types.Monitor.createFrom({}),
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
  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(getCurrentDataFacet());

    const renderActions = (actionData: ActionData) => {
      const isDeleted = actionData.isDeleted;

      return (
        <div className="action-buttons-container">
          <Action
            icon={isDeleted ? 'Undelete' : 'Delete'}
            onClick={() => {
              if (isDeleted) {
                handleUndelete(actionData.addressStr);
              } else {
                handleDelete(actionData.addressStr);
              }
            }}
            disabled={actionData.isProcessing}
            title={isDeleted ? 'Undelete' : 'Delete'}
            size="sm"
          />
          <Action
            icon="Remove"
            onClick={() => handleRemove(actionData.addressStr)}
            disabled={actionData.isProcessing || !isDeleted}
            title="Remove"
            size="sm"
          />
        </div>
      );
    };

    const getCanRemove = (row: Record<string, unknown>) => {
      const monitor = row as unknown as types.Monitor;
      return Boolean(monitor.deleted);
    };

    return actionConfig.injectActionColumn(
      baseColumns,
      renderActions,
      getCanRemove,
    );
  }, [
    getCurrentDataFacet,
    actionConfig,
    handleDelete,
    handleUndelete,
    handleRemove,
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
        viewStateKey={viewStateKey}
      />
    ),
    // prettier-ignore
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
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
