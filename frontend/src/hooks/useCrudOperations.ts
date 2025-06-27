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
  emptyItem: TItem;

  // Optional operations for specific collections
  cleanFunc?: (payload: types.Payload, addresses: string[]) => Promise<void>;
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
    dataFacetRef,
    actionConfig,
    pageClass,
    emptyItem,
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
          emptyItem,
        )
          .then(async () => {
            const result = await pageFunc(
              createPayload(dataFacetRef.current),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setPageData(result);
            setTotalItems(result.totalItems || 0);

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
      emptyItem,
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

  const handleDelete = useCallback(
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
            return { ...item, deleted: true } as TItem;
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
          crud.Operation.DELETE,
          emptyItem,
        )
          .then(async () => {
            const result = await pageFunc(
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
            handleError(err, failure('delete', address, errorMessage));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(
          err,
          `Failed to delete ${collectionName.slice(0, -1)} ${address}`,
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
      emptyItem,
      pageFunc,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
      emitSuccess,
      handleError,
      failure,
      collectionName,
    ],
  );

  const handleUndelete = useCallback(
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
            return { ...item, deleted: false } as TItem;
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
          crud.Operation.UNDELETE,
          emptyItem,
        )
          .then(async () => {
            const result = await pageFunc(
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
            handleError(err, failure('undelete', address, errorMessage));
          })
          .finally(() => {
            setTimeout(() => {
              actionConfig.stopProcessing(address);
            }, 100);
          });
      } catch (err: unknown) {
        handleError(
          err,
          `Failed to undelete ${collectionName.slice(0, -1)} ${address}`,
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
      emptyItem,
      pageFunc,
      pagination.currentPage,
      pagination.pageSize,
      sort,
      filter,
      setTotalItems,
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
          emptyItem,
        )
          .then(async () => {
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
      emptyItem,
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

  return {
    handleRemove,
    handleDelete,
    handleUndelete,
    handleAutoname,
    handleClean,
    handleCleanOne,
  };
};
