import { useCallback, useMemo } from 'react';

import { usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { crud, sdk, types } from '@models';
import { getAddressString, useEmitters, useErrorHandler } from '@utils';

import { EntityType, useActionMsgs } from './useActionMsgs';
import { usePayload } from './usePayload';

// Generic types for different page data types
export interface PageData {
  totalItems: number;
  // Add other common properties that all page data types share
}

// Minimal configuration interface for the hook
export interface CrudOperationsConfig<TPageData extends PageData, TItem> {
  // Collection identifier (e.g., 'abis', 'monitors', 'names')
  collectionName: string;

  // Function to get current data facet
  getCurrentDataFacet: () => types.DataFacet;

  // State management
  pageData: TPageData | null;
  setPageData: React.Dispatch<React.SetStateAction<TPageData | null>>;
  setTotalItems: (total: number) => void;

  // API functions
  crudFunc: (
    payload: types.Payload,
    operation: crud.Operation,
    item: TItem,
  ) => Promise<void>;
  pageFunc: (
    payload: types.Payload,
    offset: number,
    limit: number,
    sort: sdk.SortSpec,
    filter: string,
  ) => Promise<TPageData>;

  // Data facet (keeping for now)
  dataFacetRef: React.MutableRefObject<types.DataFacet>;

  // Action config (keeping for now)
  actionConfig: {
    startProcessing: (address: string) => void;
    stopProcessing: (address: string) => void;
  };

  // Collection-specific classes
  pageClass: new (data: Record<string, unknown>) => TPageData;
  updateItem: TItem;

  // Optional operations for specific collections
  cleanFunc?: (payload: types.Payload, addresses: string[]) => Promise<void>;

  // Optional post-processing function for collection-specific logic
  postFunc?: (item: TItem) => TItem;
}

export const useCrudOperations = <TPageData extends PageData, TItem>(
  config: CrudOperationsConfig<TPageData, TItem>,
) => {
  const {
    collectionName,
    getCurrentDataFacet,
    pageData,
    setPageData,
    setTotalItems,
    crudFunc,
    pageFunc,
    postFunc,
    dataFacetRef,
    actionConfig,
    pageClass,
    updateItem,
    cleanFunc,
  } = config;

  // Derive itemsProperty from collectionName (they're always the same)
  const itemsProperty = collectionName;

  // Create viewStateKey internally
  const viewStateKey = useMemo(
    (): ViewStateKey => ({
      viewName: collectionName,
      tabName: getCurrentDataFacet(),
    }),
    [collectionName, getCurrentDataFacet],
  );

  // Internal hooks - derive pagination, sorting, filtering from viewStateKey
  const { pagination, goToPage } = usePagination(viewStateKey);
  const { sort } = useSorting(viewStateKey);
  const { filter } = useFiltering(viewStateKey);

  // Internal hooks that were previously passed as props
  const createPayload = usePayload();
  const { clearError, handleError } = useErrorHandler();
  const { emitSuccess, failure } = useActionMsgs(collectionName as EntityType);
  const { emitStatus } = useEmitters();

  const handleRemove = useCallback(
    (address: string) => {
      clearError();
      actionConfig.startProcessing(address);

      try {
        // Get the items from the page data
        const items = pageData
          ? (pageData as Record<string, unknown>)[itemsProperty] || []
          : [];
        const original = [...(items as TItem[])];
        const isOnlyRowOnPage = original.length === 1;

        const optimisticValues = original.filter((item: TItem) => {
          const itemAddress = getAddressString(
            (item as Record<string, unknown>).address,
          );
          return itemAddress !== address;
        });

        setPageData((prev) => {
          if (!prev) return null;
          return new pageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunc(
          createPayload(dataFacetRef.current, address),
          crud.Operation.REMOVE,
          { ...updateItem, address } as TItem,
        )
          .then(async () => {
            const result = await pageFunc(
              createPayload(dataFacetRef.current),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );

            // Handle pagination adjustment when removing the last item on a page
            if (isOnlyRowOnPage && result.totalItems > 0) {
              const newTotalPages = Math.ceil(
                result.totalItems / pagination.pageSize,
              );
              const lastPageIndex = Math.max(0, newTotalPages - 1);

              if (lastPageIndex !== pagination.currentPage) {
                goToPage(lastPageIndex);
              }
            }

            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('remove', address);
          })
          .catch((err: unknown) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new pageClass({
                ...prev,
                [itemsProperty]: original,
              }) as TPageData;
            });
            const errorMessage =
              err instanceof Error ? err.message : String(err);
            handleError(err, failure('remove', address, errorMessage));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(
          err,
          `Failed to remove ${collectionName.slice(0, -1)} ${address}`,
        );
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
      pageData,
      itemsProperty,
      setPageData,
      pageClass,
      crudFunc,
      createPayload,
      dataFacetRef,
      updateItem,
      pageFunc,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      goToPage,
      emitSuccess,
      handleError,
      failure,
      collectionName,
    ],
  );

  const handleToggle = useCallback(
    (address: string) => {
      clearError();

      try {
        // Get the items from the page data
        const items = pageData
          ? (pageData as Record<string, unknown>)[itemsProperty] || []
          : [];
        const original = [...(items as TItem[])];

        const currentItem = original.find((item: TItem) => {
          const itemAddress = getAddressString(
            (item as Record<string, unknown>).address,
          );
          return itemAddress === address;
        });

        const isCurrentlyDeleted = Boolean(
          (currentItem as Record<string, unknown>)?.deleted,
        );
        const newDeletedState = !isCurrentlyDeleted;
        const operation = newDeletedState
          ? crud.Operation.DELETE
          : crud.Operation.UNDELETE;
        const operationName = newDeletedState ? 'delete' : 'undelete';

        const optimisticValues = original.map((item: TItem) => {
          const itemAddress = getAddressString(
            (item as Record<string, unknown>).address,
          );
          if (itemAddress === address) {
            return {
              ...item,
              deleted: newDeletedState,
              processing: true,
            } as TItem;
          }
          return item;
        });

        setPageData((prev) => {
          if (!prev) return null;
          return new pageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunc(createPayload(dataFacetRef.current, address), operation, {
          ...updateItem,
          address,
        } as TItem)
          .then(() => {
            setPageData((prev) => {
              if (!prev) return null;
              const items =
                (prev as Record<string, unknown>)[itemsProperty] || [];
              const updatedItems = (items as TItem[]).map((item: TItem) => {
                const itemAddress = getAddressString(
                  (item as Record<string, unknown>).address,
                );
                if (itemAddress === address) {
                  const updatedItem = { ...item };
                  delete (updatedItem as Record<string, unknown>).processing;
                  return updatedItem;
                }
                return item;
              });
              return new pageClass({
                ...prev,
                [itemsProperty]: updatedItems,
              }) as TPageData;
            });
            emitSuccess(operationName, address);
          })
          .catch((err: unknown) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new pageClass({
                ...prev,
                [itemsProperty]: original,
              }) as TPageData;
            });
            const errorMessage =
              err instanceof Error ? err.message : String(err);
            handleError(err, failure(operationName, address, errorMessage));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(
          err,
          `Failed to toggle delete for ${collectionName.slice(0, -1)} ${address}`,
        );
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
      pageData,
      itemsProperty,
      setPageData,
      pageClass,
      crudFunc,
      createPayload,
      dataFacetRef,
      updateItem,
      emitSuccess,
      handleError,
      failure,
      collectionName,
    ],
  );

  const handleAutoname = useCallback(
    (address: string) => {
      clearError();
      actionConfig.startProcessing(address);

      try {
        // Get the items from the page data
        const items = pageData
          ? (pageData as Record<string, unknown>)[itemsProperty] || []
          : [];
        const original = [...(items as TItem[])];

        const optimisticValues = original.map((item: TItem) => {
          const itemAddress = getAddressString(
            (item as Record<string, unknown>).address,
          );
          if (itemAddress === address) {
            return { ...item, name: 'Generating...' } as TItem;
          }
          return item;
        });

        setPageData((prev) => {
          if (!prev) return null;
          return new pageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunc(
          createPayload(dataFacetRef.current, address),
          crud.Operation.AUTONAME,
          { ...updateItem, address } as TItem,
        )
          .then(async () => {
            // For autoname, we need to refresh to get the actual generated name
            // The optimistic "Generating..." should be replaced with the real name
            // Since autoname happens in the SDK, we must refresh from backend
            const result = await pageFunc(
              createPayload(dataFacetRef.current),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            emitSuccess('autoname', address);
          })
          .catch((err: unknown) => {
            setPageData((prev) => {
              if (!prev) return null;
              return new pageClass({
                ...prev,
                [itemsProperty]: original,
              }) as TPageData;
            });
            const errorMessage =
              err instanceof Error ? err.message : String(err);
            handleError(err, failure('autoname', address, errorMessage));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(err, `Failed to autoname address ${address}`);
        actionConfig.stopProcessing(address);
      }
    },
    [
      clearError,
      actionConfig,
      pageData,
      itemsProperty,
      setPageData,
      pageClass,
      crudFunc,
      createPayload,
      dataFacetRef,
      updateItem,
      pageFunc,
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

  const handleClean = useCallback(async () => {
    if (!cleanFunc) return;

    clearError();
    actionConfig.startProcessing('clean');

    try {
      // Emit cleaning status
      if (collectionName === 'monitors') {
        emitStatus('Cleaning all monitors...');
      }
      await cleanFunc(createPayload(dataFacetRef.current), []);
      const result = await pageFunc(
        createPayload(dataFacetRef.current),
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setPageData(result);
      setTotalItems(result.totalItems || 0);
      emitSuccess('clean', 0);
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : String(err);
      handleError(err, failure('clean', undefined, errorMessage));
    } finally {
      setTimeout(() => {
        actionConfig.stopProcessing('clean');
      }, 100);
    }
  }, [
    cleanFunc,
    clearError,
    collectionName,
    createPayload,
    dataFacetRef,
    pageFunc,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setPageData,
    setTotalItems,
    emitSuccess,
    handleError,
    failure,
    emitStatus,
    actionConfig,
  ]);

  const handleCleanOne = useCallback(
    async (addresses: string[]) => {
      if (!cleanFunc) return;

      clearError();
      const firstAddress = addresses.length > 0 ? addresses[0] : null;
      const processingKey = firstAddress || 'clean-one';
      actionConfig.startProcessing(processingKey);

      try {
        // Emit cleaning status
        if (collectionName === 'monitors') {
          emitStatus(`Cleaning ${addresses.length} monitor(s)...`);
        }
        await cleanFunc(createPayload(dataFacetRef.current), addresses);
        const result = await pageFunc(
          createPayload(dataFacetRef.current),
          pagination.currentPage * pagination.pageSize,
          pagination.pageSize,
          sort,
          filter,
        );
        setPageData(result);
        setTotalItems(result.totalItems || 0);
        emitSuccess('clean', addresses.length);
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : String(err);
        handleError(err, failure('clean', undefined, errorMessage));
      } finally {
        setTimeout(() => {
          actionConfig.stopProcessing(processingKey);
        }, 100);
      }
    },
    [
      cleanFunc,
      clearError,
      collectionName,
      createPayload,
      dataFacetRef,
      pageFunc,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setPageData,
      setTotalItems,
      emitSuccess,
      handleError,
      failure,
      emitStatus,
      actionConfig,
    ],
  );

  const handleUpdate = useCallback(
    (data: Record<string, unknown>) => {
      const item = data as unknown as TItem;
      const addressStr = getAddressString(
        (item as Record<string, unknown>).address,
      );

      clearError();

      try {
        const processedItem = postFunc ? postFunc({ ...item }) : { ...item };

        const items = pageData
          ? (pageData as Record<string, unknown>)[itemsProperty] || []
          : [];
        const original = [...(items as TItem[])];

        // Optimistic UI Update
        let optimisticValues: TItem[];
        const existingItemIndex = original.findIndex((originalItem: TItem) => {
          const itemAddress = getAddressString(
            (originalItem as Record<string, unknown>).address,
          );
          return itemAddress === addressStr;
        });

        if (existingItemIndex !== -1) {
          optimisticValues = original.map((originalItem: TItem, index) =>
            index === existingItemIndex
              ? ({ ...originalItem, ...processedItem } as TItem)
              : originalItem,
          );
        } else {
          optimisticValues = [processedItem as TItem, ...original];
        }

        setPageData((prev) => {
          if (!prev) return null;
          return new pageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunc(
          createPayload(dataFacetRef.current, addressStr),
          crud.Operation.UPDATE,
          processedItem as TItem,
        )
          .then(() => {
            // For updates, we keep the optimistic update since it succeeded
            // No need to refresh the entire dataset - the optimistic update is now the truth
            emitSuccess('update', addressStr);
          })
          .catch((err: unknown) => {
            // Only on error do we revert to the original data
            setPageData((prev) => {
              if (!prev) return null;
              return new pageClass({
                ...prev,
                [itemsProperty]: original,
              }) as TPageData;
            });
            const errorMessage =
              err instanceof Error ? err.message : String(err);
            handleError(err, failure('update', addressStr, errorMessage));
          });
      } catch (err: unknown) {
        handleError(
          err,
          `Failed to update ${collectionName.slice(0, -1)} ${addressStr}`,
        );
      }
    },
    [
      clearError,
      collectionName,
      pageData,
      itemsProperty,
      setPageData,
      pageClass,
      crudFunc,
      createPayload,
      dataFacetRef,
      postFunc,
      emitSuccess,
      handleError,
      failure,
    ],
  );

  return {
    handleRemove,
    handleToggle,
    handleAutoname,
    handleClean,
    handleCleanOne,
    handleUpdate,
  };
};
