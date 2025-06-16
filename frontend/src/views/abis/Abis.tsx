import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { AbisCrud, GetAbisPage, Reload } from '@app';
import { Action, BaseTab, FormField, usePagination } from '@components';
import { TableKey, useAppContext, useFiltering, useSorting } from '@contexts';
import { useActionMsgs, useEvent } from '@hooks';
import { TabView } from '@layout';
import { useHotkeys } from '@mantine/hooks';
import { abis, crud, msgs, types } from '@models';
import { getAddressString, useErrorHandler } from '@utils';

import { Address } from '../../types/address';
import { getColumns } from './';

export const Abis = () => {
  const { lastTab } = useAppContext();
  const { emitSuccess } = useActionMsgs('abis');
  const [pageData, setPageData] = useState<abis.AbisPage | null>(null);
  const [state, setState] = useState<types.LoadState>();
  const [processingAddresses, setProcessingAddresses] = useState<Set<string>>(
    new Set(),
  );

  const [listKind, setListKind] = useState<types.ListKind>(
    lastTab[ABIS_ROUTE] || ABIS_DEFAULT_LIST,
  );
  const tableKey = useMemo(
    (): TableKey => ({ viewName: ABIS_ROUTE, tabName: listKind }),
    [listKind],
  );

  const { error, handleError, clearError } = useErrorHandler();
  const { pagination, setTotalItems, goToPage } = usePagination(tableKey);
  const { sort } = useSorting(tableKey);
  const { filter } = useFiltering(tableKey);

  const listKindRef = useRef(listKind);
  const renderCnt = useRef(0);
  // renderCnt.current++;

  useEffect(() => {
    listKindRef.current = listKind;
  }, [listKind]);

  const fetchData = useCallback(async () => {
    clearError();
    try {
      const result = await GetAbisPage(
        listKindRef.current,
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setState(result.state);
      setPageData(result);
      setTotalItems(result.totalItems || 0);
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
  ]);

  const currentData = useMemo(() => {
    // For ABI-based tabs (Downloaded, Known), use abis data
    if (
      listKind === types.ListKind.DOWNLOADED ||
      listKind === types.ListKind.KNOWN
    ) {
      return pageData?.abis || [];
    }
    // For function-based tabs (Functions, Events), use functions data
    return pageData?.functions || [];
  }, [pageData?.abis, pageData?.functions, listKind]);

  useEffect(() => {
    const currentTab = lastTab[ABIS_ROUTE];
    if (currentTab && currentTab !== listKindRef.current) {
      setListKind(currentTab);
    }
  }, [lastTab]);

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
        Reload(listKind).then(() => {
          fetchData();
        });
      },
    ],
  ]);

  // Optimistic delete action with simple last-record navigation
  const handleAction = useCallback(
    (address: Address) => {
      clearError();
      try {
        const original = [...(pageData?.abis || [])];

        // Check if we're removing the only row on the current page
        const isOnlyRowOnPage = original.length === 1;

        const optimisticValues = original.filter((abi) => {
          const abiAddress = getAddressString(abi.address);
          return abiAddress !== address;
        });

        // Apply optimistic update immediately without setting loading state
        setPageData((prev) => {
          if (!prev) return null;
          return new abis.AbisPage({
            ...prev,
            abis: optimisticValues,
          });
        });
        // Update total items count optimistically
        const currentTotal = pageData?.totalItems || 0;
        setTotalItems(Math.max(0, currentTotal - 1));

        AbisCrud(
          listKindRef.current,
          crud.Operation.REMOVE,
          {} as types.Abi,
          address,
        )
          .then(async () => {
            // Fetch fresh data to confirm the removal after successful backend operation
            const result = await GetAbisPage(
              listKindRef.current,
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setState(result.state);
            setPageData(result);
            setTotalItems(result.totalItems || 0);

            // If we removed the only row on the page, navigate to the last available record
            if (isOnlyRowOnPage && result.totalItems > 0) {
              const newTotalPages = Math.ceil(
                result.totalItems / pagination.pageSize,
              );
              const lastPageIndex = Math.max(0, newTotalPages - 1);

              // Navigate to the last page if we're not already there
              if (lastPageIndex !== pagination.currentPage) {
                goToPage(lastPageIndex);

                // After navigation, simulate pressing End key to select last row
                // This mimics the logic from useTableKeys.ts line 85
                setTimeout(() => {
                  // Simulate End key: setSelectedRowIndex(itemCount - 1)
                  // We need to send a custom event to the table to select the last row
                  const endKeyEvent = new KeyboardEvent('keydown', {
                    key: 'End',
                    bubbles: true,
                    cancelable: true,
                  });

                  // Find the table element and dispatch the End key event
                  const tableElement = document.querySelector('.data-table');
                  if (tableElement) {
                    tableElement.dispatchEvent(endKeyEvent);
                  }
                }, 200); // Give time for page navigation and data loading
              }
            }

            emitSuccess('delete', address);
          })
          .catch((err) => {
            // Revert optimistic update on error
            setState(types.LoadState.ERROR);
            setPageData((prev) => {
              if (!prev) return null;
              return new abis.AbisPage({
                ...prev,
                abis: original,
              });
            });
            // Revert total items count on error
            setTotalItems(pageData?.totalItems || 0);
            handleError(err, 'handleAction');
          });
      } finally {
        // Always clean up the processing state if needed
      }
    },
    [
      clearError,
      pageData?.abis,
      pageData?.totalItems,
      setTotalItems,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      emitSuccess,
      goToPage,
      handleError,
    ],
  );

  // Handle remove action for ABIs
  const handleRemove = useCallback(
    (address: Address) => {
      // Add address to processing set
      setProcessingAddresses((prev) => new Set(prev).add(address));

      try {
        handleAction(address);
      } finally {
        // Clean up processing state after a delay
        setTimeout(() => {
          setProcessingAddresses((prev) => {
            const newSet = new Set(prev);
            newSet.delete(address);
            return newSet;
          });
        }, 100);
      }
    },
    [handleAction],
  );

  const handleSubmit = useCallback((_formData: Record<string, unknown>) => {
    // Log(`Table submitted: ${formData}`);
  }, []);

  const currentColumns = useMemo(() => {
    const baseColumns = getColumns(pageData?.kind || ABIS_DEFAULT_LIST);

    // Only add actions for ABI-based tabs (Downloaded, Known), not for Functions/Events tabs
    const shouldShowActions =
      (pageData?.kind || ABIS_DEFAULT_LIST) === types.ListKind.DOWNLOADED ||
      (pageData?.kind || ABIS_DEFAULT_LIST) === types.ListKind.KNOWN;

    if (!shouldShowActions) {
      // For Functions/Events tabs, filter out the actions column
      return baseColumns.filter((col) => col.key !== 'actions');
    }

    // Add action buttons render function to the actions column for ABI tabs
    const actionsOverride: Partial<FormField> = {
      sortable: false,
      editable: false,
      visible: true,
      render: (row: Record<string, unknown>) => {
        const abi = row as unknown as types.Abi;
        const addressStr = getAddressString(abi.address);
        const isProcessing = processingAddresses.has(addressStr);
        const canRemove = pageData?.kind === types.ListKind.DOWNLOADED;

        return (
          <div className="action-buttons-container">
            <Action
              icon="Remove"
              onClick={() => handleRemove(addressStr)}
              disabled={isProcessing || !canRemove}
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
  }, [pageData?.kind, handleRemove, processingAddresses]);

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
        label: types.ListKind.DOWNLOADED,
        value: types.ListKind.DOWNLOADED,
        content: perTabTable,
      },
      {
        label: types.ListKind.KNOWN,
        value: types.ListKind.KNOWN,
        content: perTabTable,
      },
      {
        label: types.ListKind.FUNCTIONS,
        value: types.ListKind.FUNCTIONS,
        content: perTabTable,
      },
      {
        label: types.ListKind.EVENTS,
        value: types.ListKind.EVENTS,
        content: perTabTable,
      },
    ],
    [perTabTable],
  );

  return (
    <div className="mainView">
      {(state as string) === '' && <div>{`state: ${state}`}</div>}
      <TabView tabs={tabs} route={ABIS_ROUTE} />
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

const ABIS_DEFAULT_LIST = types.ListKind.DOWNLOADED;
const ABIS_ROUTE = '/abis';
