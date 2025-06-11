// MONITORS_ROUTE
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { MonitorsCrud } from '@app';
import { Action, BaseTab, FormField, usePagination } from '@components';
import { TableKey, useAppContext, useFiltering, useSorting } from '@contexts';
import { useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { crud, facets, monitors, msgs, types } from '@models';
import { getAddressString, useEmitters, useErrorHandler } from '@utils';

import { Address } from '../../types/address';
import {
  ACTION_MESSAGES,
  MONITORS_DEFAULT_LIST,
  MONITORS_ROUTE,
  cleanMonitors,
  getColumns,
  getMonitorsPage,
  reload,
} from './';

export const Monitors = () => {
  const { lastTab } = useAppContext();
  const [pageData, setPageData] = useState<monitors.MonitorsPage | null>(null);
  const [state, setState] = useState<facets.LoadState>();
  const [processingAddresses, setProcessingAddresses] = useState<Set<string>>(
    new Set(),
  );

  const [listKind, setListKind] = useState<types.ListKind>(
    lastTab[MONITORS_ROUTE] || MONITORS_DEFAULT_LIST,
  );
  const tableKey = useMemo(
    (): TableKey => ({ viewName: MONITORS_ROUTE, tabName: listKind }),
    [listKind],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter } = useFiltering(tableKey);
  const { emitStatus } = useEmitters();

  const listKindRef = useRef(listKind);
  const renderCnt = useRef(0);

  useEffect(() => {
    listKindRef.current = listKind;
  }, [listKind]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await getMonitorsPage(
        listKindRef.current,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setState(result.state);
      setPageData(result);
      setTotalItems(result.totalItems || 0);

      // Emit status message after successful data load
      emitStatus(`Loaded ${result.totalItems || 0} monitors successfully`);
    } catch (err: unknown) {
      handleError(err, `Failed to fetch ${listKindRef.current}`);
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

  useEffect(() => {
    const currentTab = lastTab[MONITORS_ROUTE];
    if (currentTab && currentTab !== listKind) {
      setListKind(currentTab);
    }
  }, [lastTab, listKind]);

  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: types.DataLoadedPayload) => {
      if (payload?.listKind === listKindRef.current) {
        fetchData();
      }
    },
  );

  useEffect(() => {
    fetchData();
  }, [fetchData, listKind]);

  useHotkeys([
    [
      'mod+r',
      () => {
        reload().then(() => {
          fetchData();
        });
      },
    ],
    [
      'mod+shift+c',
      () => {
        handleCleanAll();
      },
    ],
  ]);

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
        setState(facets.LoadState.PENDING);
        setPageData((prev) => {
          if (!prev) return null;
          return new monitors.MonitorsPage({
            ...prev,
            monitors: optimisticValues,
          });
        });
        MonitorsCrud(
          listKindRef.current,
          crud.Operation.DELETE,
          {} as types.Monitor,
          address,
        )
          .then(async () => {
            const result = await getMonitorsPage(
              listKindRef.current,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setState(result.state);
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitStatus(ACTION_MESSAGES.DELETE_SUCCESS(address));
          })
          .catch((err) => {
            setState(facets.LoadState.ERROR);
            setPageData((prev) => {
              if (!prev) return null;
              return new monitors.MonitorsPage({
                ...prev,
                monitors: original,
              });
            });
            handleError(
              err,
              ACTION_MESSAGES.DELETE_FAILURE(address, err.message),
            );
          });
      } catch (err: unknown) {
        handleError(err, `Failed to delete monitor ${address}`);
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
      emitStatus,
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
        setState(facets.LoadState.PENDING);
        setPageData((prev) => {
          if (!prev) return null;
          return new monitors.MonitorsPage({
            ...prev,
            monitors: optimisticValues,
          });
        });
        MonitorsCrud(
          listKindRef.current,
          crud.Operation.UNDELETE,
          {} as types.Monitor,
          address,
        )
          .then(async () => {
            const result = await getMonitorsPage(
              listKindRef.current,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setState(result.state);
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitStatus(ACTION_MESSAGES.UNDELETE_SUCCESS(address));
          })
          .catch((err) => {
            setState(facets.LoadState.ERROR);
            setPageData((prev) => {
              if (!prev) return null;
              return new monitors.MonitorsPage({
                ...prev,
                monitors: original,
              });
            });
            handleError(
              err,
              ACTION_MESSAGES.UNDELETE_FAILURE(address, err.message),
            );
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
      emitStatus,
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
        setState(facets.LoadState.PENDING);
        setPageData((prev) => {
          if (!prev) return null;
          return new monitors.MonitorsPage({
            ...prev,
            monitors: optimisticValues,
          });
        });
        MonitorsCrud(
          listKindRef.current,
          crud.Operation.REMOVE,
          {} as types.Monitor,
          address,
        )
          .then(async () => {
            const result = await getMonitorsPage(
              listKindRef.current,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setState(result.state);
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitStatus(ACTION_MESSAGES.REMOVE_SUCCESS(address));
          })
          .catch((err) => {
            setState(facets.LoadState.ERROR);
            setPageData((prev) => {
              if (!prev) return null;
              return new monitors.MonitorsPage({
                ...prev,
                monitors: original,
              });
            });
            handleError(
              err,
              ACTION_MESSAGES.REMOVE_FAILURE(address, err.message),
            );
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
      emitStatus,
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

  // Handle clean all monitors
  const handleCleanAll = useCallback(async () => {
    clearError();
    try {
      emitStatus('Cleaning all monitors...');
      await cleanMonitors([]);
      await fetchData();
      emitStatus(ACTION_MESSAGES.CLEAN_SUCCESS(0));
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : String(err);
      handleError(err, ACTION_MESSAGES.CLEAN_FAILURE(errorMessage));
    }
  }, [clearError, fetchData, emitStatus, handleError]);

  // Handle clean selected monitors
  const _handleCleanSelected = useCallback(
    async (addresses: string[]) => {
      clearError();
      try {
        emitStatus(`Cleaning ${addresses.length} monitor(s)...`);
        await cleanMonitors(addresses);
        await fetchData();
        emitStatus(ACTION_MESSAGES.CLEAN_SUCCESS(addresses.length));
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : String(err);
        handleError(err, ACTION_MESSAGES.CLEAN_FAILURE(errorMessage));
      }
    },
    [clearError, fetchData, emitStatus, handleError],
  );

  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(listKind);

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
  }, [listKind, handleMonitorAction, processingAddresses]);

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

  const perTabTable = useMemo(
    () => (
      <BaseTab
        data={currentData as unknown as Record<string, unknown>[]}
        columns={currentColumns}
        loading={!!pageData?.isFetching}
        error={error}
        onSubmit={handleSubmit}
        tableKey={tableKey}
      />
    ),
    [
      currentData,
      currentColumns,
      pageData?.isFetching,
      error,
      handleSubmit,
      tableKey,
    ],
  );

  const tabs = useMemo(
    () => [
      {
        label: 'Monitors',
        value: types.ListKind.MONITORS,
        content: perTabTable,
      },
    ],
    [perTabTable],
  );

  return (
    <div className="mainView">
      {(state as string) === '' && <div>{`state: ${state}`}</div>}
      <TabView tabs={tabs} route={MONITORS_ROUTE} />
      {error && (
        <div>
          <h3>{`Error fetching ${listKind}`}</h3>
          <p>{error.message}</p>
        </div>
      )}
      {renderCnt.current > 0 && <div>{`renderCnt: ${renderCnt.current}`}</div>}
    </div>
  );
};

// MONITORS_ROUTE
