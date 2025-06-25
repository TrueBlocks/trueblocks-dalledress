// === SECTION 1: Imports & Dependencies ===
// EXISTING_CODE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetMonitorsPage, MonitorsClean, MonitorsCrud, Reload } from '@app';
import { Action, BaseTab, FormField, usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import {
  DataFacetConfig,
  useActionMsgs,
  useActiveFacet,
  useEvent,
} from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { crud, monitors, msgs, types } from '@models';
import { getAddressString, useEmitters, useErrorHandler } from '@utils';

import { Address } from '../../types/address';
import { getColumns } from './';
import {
  MONITORS_DEFAULT_FACET,
  MONITORS_ROUTE as ROUTE,
  monitorsFacets,
} from './facets';

// EXISTING_CODE
// === END SECTION 1 ===

export const Monitors = () => {
  // === SECTION 2: Hook Initialization ===
  // EXISTING_CODE
  const activeFacetHook = useActiveFacet({
    facets: monitorsFacets,
    defaultFacet: MONITORS_DEFAULT_FACET,
    viewRoute: ROUTE,
  });

  const { getCurrentDataFacet } = activeFacetHook;

  const { emitSuccess, emitCleaningStatus, failure } =
    useActionMsgs('monitors');
  const [pageData, setPageData] = useState<monitors.MonitorsPage | null>(null);
  const [processingAddresses, setProcessingAddresses] = useState<Set<string>>(
    new Set(),
  );

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
  const { emitStatus } = useEmitters();
  // EXISTING_CODE
  // === END SECTION 2 ===

  // === SECTION 3: Refs & Effects Setup ===
  // EXISTING_CODE
  const dataFacetRef = useRef(getCurrentDataFacet());
  useEffect(() => {
    dataFacetRef.current = getCurrentDataFacet();
  }, [getCurrentDataFacet]);
  // EXISTING_CODE
  // === END SECTION 3 ===

  // === SECTION 4: Data Fetching Logic ===
  // EXISTING_CODE
  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetMonitorsPage(
        types.Payload.createFrom({ dataFacet: dataFacetRef.current }),
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setPageData(result);
      setTotalItems(result.totalItems || 0);

      // Emit status message after successful data load
      emitStatus(`Loaded ${result.totalItems || 0} monitors successfully`);
    } catch (err: unknown) {
      handleError(err, `Failed to fetch ${dataFacetRef.current}`);
    }
  }, [
    clearError,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setTotalItems,
    handleError,
    emitStatus,
  ]);

  const currentData = useMemo(() => {
    return pageData?.monitors || [];
  }, [pageData?.monitors]);
  // EXISTING_CODE
  // === END SECTION 4 ===

  // === SECTION 5: Event Handling ===
  // EXISTING_CODE
  useHotkeys([['mod+shift+c', () => handleCleanAll()]]);
  // EXISTING_CODE
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: Record<string, unknown>) => {
      if (payload?.collection === 'monitors') {
        const eventDataFacet = payload.dataFacet as types.DataFacet | undefined;
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
      await Reload(getCurrentDataFacet() as types.DataFacet, '', '');
      await fetchData();
    } catch (err: unknown) {
      handleError(err, `Failed to reload ${getCurrentDataFacet()}`);
    }
  }, [getCurrentDataFacet, fetchData, handleError]);

  useHotkeys([['mod+r', handleReload]]);
  // === END SECTION 5 ===

  // === SECTION 6: CRUD Operations ===
  // EXISTING_CODE

  // Handle CRUD actions for monitors
  const handleDelete = useCallback(
    (address: Address) => {
      clearError();
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
          types.Payload.createFrom({
            dataFacet: dataFacetRef.current,
            address: address,
          }),
          crud.Operation.DELETE,
          {} as types.Monitor,
        )
          .then(async () => {
            const result = await GetMonitorsPage(
              types.Payload.createFrom({ dataFacet: dataFacetRef.current }),
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
          });
      } catch (err: unknown) {
        handleError(err, `Failed to delete monitor ${address}`);
      }
    },
    [
      clearError,
      pageData?.monitors,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      handleError,
      failure,
    ],
  );

  const handleUndelete = useCallback(
    (address: Address) => {
      clearError();
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
          types.Payload.createFrom({
            dataFacet: dataFacetRef.current,
            address: address,
          }),
          crud.Operation.UNDELETE,
          {} as types.Monitor,
        )
          .then(async () => {
            const result = await GetMonitorsPage(
              types.Payload.createFrom({ dataFacet: dataFacetRef.current }),
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
          });
      } catch (err: unknown) {
        handleError(err, `Failed to undelete monitor ${address}`);
      }
    },
    [
      clearError,
      pageData?.monitors,
      handleError,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      failure,
    ],
  );

  const _handleRemove = useCallback(
    (address: Address) => {
      clearError();
      try {
        const original = [...(pageData?.monitors || [])];
        const optimisticValues = original.filter((monitor) => {
          const monitorAddress = getAddressString(monitor.address);
          return monitorAddress !== address;
        });
        setPageData((prev) => {
          if (!prev) return null;
          return new monitors.MonitorsPage({
            ...prev,
            monitors: optimisticValues,
          });
        });
        MonitorsCrud(
          types.Payload.createFrom({
            dataFacet: dataFacetRef.current,
            address: address,
          }),
          crud.Operation.REMOVE,
          {} as types.Monitor,
        )
          .then(async () => {
            const result = await GetMonitorsPage(
              types.Payload.createFrom({ dataFacet: dataFacetRef.current }),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('remove', address);
          })
          .catch((err) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new monitors.MonitorsPage({
                ...prev,
                monitors: original,
              });
            });
            handleError(err, failure('remove', address, err.message));
          });
      } catch (err: unknown) {
        handleError(err, `Failed to remove monitor ${address}`);
      }
    },
    [
      clearError,
      pageData?.monitors,
      handleError,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      failure,
    ],
  );

  // Combined action handler for monitors
  const handleMonitorAction = useCallback(
    (
      address: Address,
      isDeleted: boolean,
      actionType: 'delete' | 'undelete' | 'remove',
    ) => {
      // Add address to processing set
      setProcessingAddresses((prev) => new Set(prev).add(address));

      try {
        switch (actionType) {
          case 'delete':
            handleDelete(address);
            break;
          case 'undelete':
            handleUndelete(address);
            break;
          case 'remove':
            _handleRemove(address);
            break;
        }
      } finally {
        // Clean up processing state after a delay to allow for optimistic updates
        setTimeout(() => {
          setProcessingAddresses((prev) => {
            const newSet = new Set(prev);
            newSet.delete(address);
            return newSet;
          });
        }, 100);
      }
    },
    [handleDelete, handleUndelete, _handleRemove],
  );

  const handleCleanAll = useCallback(async () => {
    clearError();
    try {
      emitCleaningStatus();
      await MonitorsClean(
        types.Payload.createFrom({ dataFacet: dataFacetRef.current }),
        [],
      );
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
  ]);

  // Handle clean selected monitors
  const _handleCleanSelected = useCallback(
    async (addresses: string[]) => {
      clearError();
      try {
        emitCleaningStatus(addresses.length);
        await MonitorsClean(
          types.Payload.createFrom({ dataFacet: dataFacetRef.current }),
          addresses,
        );
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
    ],
  );
  // EXISTING_CODE
  // === END SECTION 6 ===

  // === SECTION 7: Form & UI Handlers ===
  // EXISTING_CODE
  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(getCurrentDataFacet() as types.DataFacet);

    // Add action buttons render function to the actions column
    const actionsOverride: Partial<FormField> = {
      sortable: false,
      editable: false,
      visible: true,
      render: (row: Record<string, unknown>) => {
        const monitor = row as unknown as types.Monitor;
        const isDeleted = Boolean(monitor.deleted);
        const addressStr = getAddressString(monitor.address);
        const isProcessing = processingAddresses.has(addressStr);

        return (
          <div className="action-buttons-container">
            <Action
              icon={isDeleted ? 'Undelete' : 'Delete'}
              onClick={() =>
                handleMonitorAction(
                  addressStr,
                  isDeleted,
                  isDeleted ? 'undelete' : 'delete',
                )
              }
              disabled={isProcessing}
              title={isDeleted ? 'Undelete' : 'Delete'}
              size="sm"
            />
            <Action
              icon="Remove"
              onClick={() =>
                handleMonitorAction(addressStr, isDeleted, 'remove')
              }
              disabled={isProcessing || !isDeleted}
              title="Remove"
              size="sm"
            />
          </div>
        );
      },
    };

    return baseColumns.map((col) =>
      col.key === 'actions' ? { ...col, ...actionsOverride } : col,
    );
  }, [getCurrentDataFacet, handleMonitorAction, processingAddresses]);

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
  // EXISTING_CODE
  // === END SECTION 7 ===

  // === SECTION 8: Tab Configuration ===
  // EXISTING_CODE
  // EXISTING_CODE
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
      activeFacetHook.availableFacets.map((facetConfig: DataFacetConfig) => ({
        label: facetConfig.label,
        value: facetConfig.id,
        content: perTabTable,
      })),
    [activeFacetHook.availableFacets, perTabTable],
  );
  // === END SECTION 8 ===

  // === SECTION 9: Render/JSX ===
  // EXISTING_CODE
  // EXISTING_CODE
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
