// === SECTION 1: Imports & Dependencies ===
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetMonitorsPage, MonitorsClean, MonitorsCrud, Reload } from '@app';
import { Action } from '@components';
import { BaseTab, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import {
  ActionData,
  useActionConfig,
  useActionMsgs,
  useCrudOperations,
} from '@hooks';
import { DataFacetConfig, useActiveFacet, useEvent, usePayload } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { crud, monitors, msgs, types } from '@models';
import { getAddressString, useErrorHandler } from '@utils';

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
  const dataFacetRef = useRef(getCurrentDataFacet() as types.DataFacet);
  useEffect(() => {
    dataFacetRef.current = getCurrentDataFacet() as types.DataFacet;
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
      if (payload?.collection === ROUTE) {
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
  // EXISTING_CODE
  const actionConfig = useActionConfig({
    operations: ['delete', 'undelete', 'remove'],
  });

  const { emitSuccess, emitCleaningStatus, failure } =
    useActionMsgs('monitors');

  // Use the new CRUD operations hook for handleRemove
  const { handleRemove } = useCrudOperations({
    collectionName: ROUTE,
    getCurrentDataFacet,
    pageData,
    setPageData,
    setTotalItems,
    crudFunction: MonitorsCrud,
    getPageFunction: GetMonitorsPage,
    dataFacetRef,
    actionConfig,
    PageClass: monitors.MonitorsPage,
    emptyItem: types.Monitor.createFrom({}),
  });

  const handleDelete = useCallback(
    (address: string) => {
      clearError();
      actionConfig.startProcessing(address);

      try {
        const original = [...(pageData?.monitors || [])];
        const optimisticValues = original.map((monitor) => {
          const monitorAddress = getAddressString(monitor.address);
          if (monitorAddress === address) {
            return { ...monitor, deleted: true };
          }
          return monitor;
        });
        setPageData((prev) => {
          if (!prev) return null;
          return new monitors.MonitorsPage({
            ...prev,
            monitors: optimisticValues,
          });
        });
        MonitorsCrud(
          createPayload(dataFacetRef.current, address),
          crud.Operation.DELETE,
          {} as types.Monitor,
        )
          .then(async () => {
            const result = await GetMonitorsPage(
              createPayload(dataFacetRef.current),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('delete', address);
          })
          .catch((err) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new monitors.MonitorsPage({
                ...prev,
                monitors: original,
              });
            });
            handleError(err, failure('delete', address, err.message));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(err, `Failed to delete monitor ${address}`);
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
      pageData?.monitors,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      handleError,
      failure,
      createPayload,
    ],
  );

  const handleUndelete = useCallback(
    (address: string) => {
      clearError();
      actionConfig.startProcessing(address);

      try {
        const original = [...(pageData?.monitors || [])];
        const optimisticValues = original.map((monitor) => {
          const monitorAddress = getAddressString(monitor.address);
          if (monitorAddress === address) {
            return { ...monitor, deleted: false };
          }
          return monitor;
        });
        setPageData((prev) => {
          if (!prev) return null;
          return new monitors.MonitorsPage({
            ...prev,
            monitors: optimisticValues,
          });
        });
        MonitorsCrud(
          createPayload(dataFacetRef.current, address),
          crud.Operation.UNDELETE,
          {} as types.Monitor,
        )
          .then(async () => {
            const result = await GetMonitorsPage(
              createPayload(dataFacetRef.current),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('undelete', address);
          })
          .catch((err) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new monitors.MonitorsPage({
                ...prev,
                monitors: original,
              });
            });
            handleError(err, failure('undelete', address, err.message));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(err, `Failed to undelete monitor ${address}`);
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
      pageData?.monitors,
      handleError,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      failure,
      createPayload,
    ],
  );

  const _handleCleanAll = useCallback(async () => {
    clearError();
    try {
      emitCleaningStatus();
      await MonitorsClean(createPayload(dataFacetRef.current), []);
      await fetchData();
      emitSuccess('clean', 0);
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : String(err);
      handleError(err, failure('clean', undefined, errorMessage));
    }
  }, [
    clearError,
    fetchData,
    emitCleaningStatus,
    emitSuccess,
    failure,
    handleError,
    createPayload,
  ]);

  // Handle clean selected monitors
  const _handleCleanSelected = useCallback(
    async (addresses: string[]) => {
      clearError();
      try {
        emitCleaningStatus(addresses.length);
        await MonitorsClean(createPayload(dataFacetRef.current), addresses);
        await fetchData();
        emitSuccess('clean', addresses.length);
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : String(err);
        handleError(err, failure('clean', undefined, errorMessage));
      }
    },
    [
      clearError,
      fetchData,
      emitCleaningStatus,
      emitSuccess,
      failure,
      handleError,
      createPayload,
    ],
  );
  // EXISTING_CODE
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE
  const handleSubmit = useCallback(
    (data: Record<string, unknown>) => {
      const monitor = data as unknown as types.Monitor;
      const address = getAddressString(monitor.address);

      if (monitor.deleted) {
        handleUndelete(address);
      } else {
        handleDelete(address);
      }
    },
    [handleDelete, handleUndelete],
  );

  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(getCurrentDataFacet() as types.DataFacet);

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
        onSubmit={handleSubmit}
        viewStateKey={viewStateKey}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      handleSubmit,
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
