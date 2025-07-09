import { useCallback, useMemo } from 'react';

import { ViewStateKey } from '@contexts';
import { useWalletGatedAction } from '@hooks';
import { crud, sdk, types } from '@models';
import { Log, isDebugMode, useErrorHandler } from '@utils';

import { useActionMsgs } from './useActionMsgs';

const debug = isDebugMode();

// Helper function for getting address string - moved outside component for stability
const getAddressString = (address: unknown): string => {
  if (typeof address === 'string') return address;
  if (address && typeof address === 'object' && 'address' in address) {
    return String((address as { address: unknown }).address);
  }
  return String(address || '');
};

export type ActionType =
  | 'publish'
  | 'pin'
  | 'add'
  | 'delete'
  | 'remove'
  | 'autoname'
  | 'clean'
  | 'update';

// Action definition with level, wallet requirements
export interface ActionDefinition {
  type: ActionType;
  level: 'row' | 'header' | 'both';
  requiresWallet: boolean;
  title: string;
  icon: string;
}

export interface CollectionActionsConfig<TPageData, TItem> {
  // Collection identifier
  collection: string;

  // Current context
  viewStateKey: ViewStateKey;

  // Required hooks from parent component
  pagination: ReturnType<
    typeof import('@components').usePagination
  >['pagination'];
  goToPage: ReturnType<typeof import('@components').usePagination>['goToPage'];
  sort: ReturnType<typeof import('@contexts').useSorting>['sort'];
  filter: ReturnType<typeof import('@contexts').useFiltering>['filter'];

  // Enabled actions for this collection
  enabledActions: ActionDefinition['type'][];

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

  // Collection-specific classes and items
  pageClass: new (data: Record<string, unknown>) => TPageData;
  updateItem: TItem;

  // Optional functions
  cleanFunc?: (payload: types.Payload, addresses: string[]) => Promise<void>;
  postFunc?: (item: TItem) => TItem;

  // Helper functions
  createPayload: (facet: types.DataFacet, address?: string) => types.Payload;
  getCurrentDataFacet: () => types.DataFacet;
}

// Predefined action definitions
const ACTION_DEFINITIONS: Record<string, ActionDefinition> = {
  publish: {
    type: 'publish',
    level: 'header',
    requiresWallet: true,
    title: 'Publish',
    icon: 'Publish',
  },
  pin: {
    type: 'pin',
    level: 'header',
    requiresWallet: false,
    title: 'Pin',
    icon: 'Pin',
  },
  add: {
    type: 'add',
    level: 'header',
    requiresWallet: false,
    title: 'Add',
    icon: 'Add',
  },
  delete: {
    type: 'delete',
    level: 'row',
    requiresWallet: false,
    title: 'Delete',
    icon: 'Delete',
  },
  remove: {
    type: 'remove',
    level: 'row',
    requiresWallet: false,
    title: 'Remove',
    icon: 'Remove',
  },
  autoname: {
    type: 'autoname',
    level: 'row',
    requiresWallet: false,
    title: 'Auto-generate name',
    icon: 'Autoname',
  },
  clean: {
    type: 'clean',
    level: 'header',
    requiresWallet: false,
    title: 'Clean',
    icon: 'Clean',
  },
  update: {
    type: 'update',
    level: 'row',
    requiresWallet: false,
    title: 'Update',
    icon: 'Edit',
  },
};

export const useActions = <TPageData extends { totalItems: number }, TItem>(
  config: CollectionActionsConfig<TPageData, TItem>,
) => {
  const {
    collection,
    pagination,
    goToPage,
    sort,
    filter,
    enabledActions,
    pageData,
    setPageData,
    setTotalItems,
    crudFunc,
    pageFunc,
    postFunc,
    createPayload,
    getCurrentDataFacet,
    pageClass,
    updateItem,
    cleanFunc,
  } = config;

  const { isWalletConnected, createWalletGatedAction } = useWalletGatedAction();
  const { emitSuccess } = useActionMsgs(
    collection as 'names' | 'monitors' | 'abis',
  );

  const currentDataFacet = useMemo(
    () => getCurrentDataFacet(),
    [getCurrentDataFacet],
  );

  // Items property name (derived from collection name)
  const itemsProperty = collection;

  // Filter actions based on enabled actions
  const availableActions = useMemo(() => {
    return enabledActions
      .map((actionType) => {
        const baseAction = ACTION_DEFINITIONS[actionType];
        if (!baseAction) return null;
        return baseAction;
      })
      .filter((action): action is ActionDefinition => {
        return !!action;
      });
  }, [enabledActions]);

  // Separate actions by level
  const headerActions = useMemo(
    () =>
      availableActions.filter(
        (action) => action.level === 'header' || action.level === 'both',
      ),
    [availableActions],
  );

  const rowActions = useMemo(
    () =>
      availableActions.filter(
        (action) => action.level === 'row' || action.level === 'both',
      ),
    [availableActions],
  );

  const { clearError, handleError } = useErrorHandler();

  // Handler implementations
  const handleAdd = useCallback(() => {
    Log(`Adding ${collection}`);
    // TODO: Implement add functionality
  }, [collection]);

  const handlePublish = createWalletGatedAction((walletAddress: string) => {
    Log(`Publishing ${collection} with wallet:`, walletAddress);
    // TODO: Implement publish functionality with connected wallet
  }, collection);

  const handlePin = useCallback(() => {
    Log(`Pinning ${collection}`);
    // TODO: Implement pin functionality
  }, [collection]);

  const handleToggle = useCallback(
    (address: string) => {
      clearError();

      try {
        const items = pageData
          ? (pageData as Record<string, unknown>)[itemsProperty]
          : null;
        if (!Array.isArray(items)) {
          handleError(
            new Error(`No valid items array found for ${collection}`),
            `handleToggle for ${address}`,
          );
          return;
        }
        const original = [...(items as TItem[])];

        const currentItem = original.find((item: TItem) => {
          const itemAddress = getAddressString(
            (item as Record<string, unknown>).address,
          );
          return itemAddress === address;
        });

        if (!currentItem) {
          handleError(
            new Error(`Item ${address} not found in ${collection}`),
            `handleToggle for ${address}`,
          );
          return;
        }

        const isCurrentlyDeleted = Boolean(
          (currentItem as Record<string, unknown>)?.deleted,
        );
        const newDeletedState = !isCurrentlyDeleted;
        const operation = newDeletedState
          ? crud.Operation.DELETE
          : crud.Operation.UNDELETE;

        const optimisticValues = original.map((item: TItem) => {
          const itemAddress = getAddressString(
            (item as Record<string, unknown>).address,
          );
          if (itemAddress === address) {
            return {
              ...item,
              deleted: newDeletedState,
            } as TItem;
          }
          return item;
        });

        setPageData((prev) => {
          if (!prev) return prev;
          return new pageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunc(createPayload(currentDataFacet, address), operation, {
          ...updateItem,
          address,
        } as TItem)
          .then(() => {
            // Success handled by visual state change
            if (debug) {
              const action = newDeletedState ? 'delete' : 'undelete';
              emitSuccess(action as 'delete' | 'undelete', address);
            }
          })
          .catch((err: unknown) => {
            handleError(err, `Failed to toggle delete for ${address}`);
            setPageData((prev) => {
              if (!prev) return prev;
              return new pageClass({
                ...prev,
                [itemsProperty]: original,
              }) as TPageData;
            });
          });
      } catch (err: unknown) {
        handleError(err, `Error in handleToggle for ${address}`);
      }
    },
    [
      pageData,
      itemsProperty,
      setPageData,
      pageClass,
      crudFunc,
      createPayload,
      updateItem,
      collection,
      currentDataFacet,
      clearError,
      handleError,
      emitSuccess,
    ],
  );

  const handleRemove = useCallback(
    (address: string) => {
      clearError();

      try {
        const items = pageData
          ? (pageData as Record<string, unknown>)[itemsProperty]
          : null;
        if (!Array.isArray(items)) {
          handleError(
            new Error(`No valid items array found for ${collection}`),
            `handleRemove for ${address}`,
          );
          return;
        }
        const original = [...(items as TItem[])];
        const isOnlyRowOnPage = original.length === 1;

        const optimisticValues = original.filter((item: TItem) => {
          const itemAddress = getAddressString(
            (item as Record<string, unknown>).address,
          );
          return itemAddress !== address;
        });

        setPageData((prev) => {
          if (!prev) return prev;
          return new pageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunc(
          createPayload(currentDataFacet, address),
          crud.Operation.REMOVE,
          { ...updateItem, address } as TItem,
        )
          .then(async () => {
            const result = await pageFunc(
              createPayload(currentDataFacet),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );

            if (isOnlyRowOnPage && result.totalItems > 0) {
              const totalPages = Math.ceil(
                result.totalItems / pagination.pageSize,
              );
              if (pagination.currentPage >= totalPages) {
                goToPage(Math.max(0, totalPages - 1));
                // Continue to set data - don't return early!
              }
            }

            setPageData(result);
            setTotalItems(result.totalItems || 0);
            if (debug) {
              emitSuccess('remove', address);
            }
          })
          .catch((err: unknown) => {
            handleError(err, `Failed to remove ${address}`);
            setPageData((prev) => {
              if (!prev) return prev;
              return new pageClass({
                ...prev,
                [itemsProperty]: original,
              }) as TPageData;
            });
          });
      } catch (err: unknown) {
        handleError(err, `Error in handleRemove for ${address}`);
      }
    },
    [
      pageData,
      itemsProperty,
      setPageData,
      pageClass,
      crudFunc,
      createPayload,
      updateItem,
      pageFunc,
      pagination,
      sort,
      filter,
      setTotalItems,
      goToPage,
      collection,
      currentDataFacet,
      clearError,
      handleError,
      emitSuccess,
    ],
  );

  const handleAutoname = useCallback(
    (address: string) => {
      clearError();

      try {
        const items = pageData
          ? (pageData as Record<string, unknown>)[itemsProperty]
          : null;
        if (!Array.isArray(items)) {
          handleError(
            new Error(`No valid items array found for ${collection}`),
            `handleAutoname for ${address}`,
          );
          return;
        }
        const original = [...(items as TItem[])];

        const optimisticValues = original.map((item: TItem) => {
          const itemAddress = getAddressString(
            (item as Record<string, unknown>).address,
          );
          if (itemAddress === address) {
            return {
              ...item,
              name: 'Generating...',
            } as TItem;
          }
          return item;
        });

        setPageData((prev) => {
          if (!prev) return prev;
          return new pageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunc(
          createPayload(currentDataFacet, address),
          crud.Operation.AUTONAME,
          { ...updateItem, address } as TItem,
        )
          .then(async () => {
            // Refresh to get the actual generated name
            const result = await pageFunc(
              createPayload(currentDataFacet),
              pagination.currentPage * pagination.pageSize,
              pagination.pageSize,
              sort,
              filter,
            );
            setPageData(result);
            setTotalItems(result.totalItems || 0);
            if (debug) {
              emitSuccess('autoname', address);
            }
          })
          .catch((err: unknown) => {
            handleError(err, `Failed to auto-name ${address}`);
            setPageData((prev) => {
              if (!prev) return prev;
              return new pageClass({
                ...prev,
                [itemsProperty]: original,
              }) as TPageData;
            });
          });
      } catch (err: unknown) {
        handleError(err, `Error in handleAutoname for ${address}`);
      }
    },
    [
      pageData,
      itemsProperty,
      setPageData,
      pageClass,
      crudFunc,
      createPayload,
      updateItem,
      pageFunc,
      pagination,
      sort,
      filter,
      setTotalItems,
      collection,
      currentDataFacet,
      clearError,
      handleError,
      emitSuccess,
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
          ? (pageData as Record<string, unknown>)[itemsProperty]
          : null;
        if (!Array.isArray(items)) {
          handleError(
            new Error(`No valid items array found for ${collection}`),
            `handleUpdate for ${addressStr}`,
          );
          return;
        }
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
          optimisticValues = [...original];
          optimisticValues[existingItemIndex] = processedItem as TItem;
        } else {
          optimisticValues = [processedItem as TItem, ...original];
        }

        setPageData((prev) => {
          if (!prev) return prev;
          return new pageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunc(
          createPayload(currentDataFacet, addressStr),
          crud.Operation.UPDATE,
          processedItem as TItem,
        )
          .then(() => {
            // Success handled by visual state change
            if (debug) {
              emitSuccess('update', addressStr);
            }
          })
          .catch((err: unknown) => {
            setPageData((prev) => {
              if (!prev) return prev;
              return new pageClass({
                ...prev,
                [itemsProperty]: original,
              }) as TPageData;
            });
            handleError(err, `Failed to update ${addressStr}`);
          });
      } catch (err: unknown) {
        handleError(err, `Error in handleUpdate for ${addressStr}`);
      }
    },
    [
      pageData,
      itemsProperty,
      setPageData,
      pageClass,
      crudFunc,
      createPayload,
      postFunc,
      collection,
      currentDataFacet,
      clearError,
      handleError,
      emitSuccess,
    ],
  );

  const handleClean = useCallback(async () => {
    if (!cleanFunc) return;

    clearError();

    try {
      await cleanFunc(createPayload(currentDataFacet), []);
      const result = await pageFunc(
        createPayload(currentDataFacet),
        pagination.currentPage * pagination.pageSize,
        pagination.pageSize,
        sort,
        filter,
      );
      setPageData(result);
      setTotalItems(result.totalItems || 0);
      if (debug) {
        emitSuccess('clean');
      }
    } catch (err: unknown) {
      handleError(err, `Failed to clean ${collection}`);
    }
  }, [
    cleanFunc,
    clearError,
    createPayload,
    currentDataFacet,
    pageFunc,
    pagination.currentPage,
    pagination.pageSize,
    sort,
    filter,
    setPageData,
    setTotalItems,
    handleError,
    collection,
    emitSuccess,
  ]);

  // TODO: Implement handleCleanOne if needed for cleaning specific addresses
  // const handleCleanOne = useCallback(
  //   async (addresses: string[]) => {
  //     if (!cleanFunc) return;
  //
  //     clearError();
  //     const firstAddress = addresses.length > 0 ? addresses[0] : null;
  //     const processingKey = firstAddress || 'clean-one';
  //     actionConfig.startProcessing(processingKey);
  //
  //     try {
  //       // Emit cleaning status
  //       if (collection === 'monitors') {
  //         emitStatus(`Cleaning ${addresses.length} monitor(s)...`);
  //       }
  //       await cleanFunc(createPayload(dataFacetRef.current), addresses);
  //       const result = await pageFunc(
  //         createPayload(dataFacetRef.current),
  //         pagination.currentPage * pagination.pageSize,
  //         pagination.pageSize,
  //         sort,
  //         filter,
  //       );
  //       setPageData(result);
  //       setTotalItems(result.totalItems || 0);
  //       emitSuccess('clean', addresses.length);
  //     } catch (err: unknown) {
  //       const errorMessage = err instanceof Error ? err.message : String(err);
  //       handleError(err, failure('clean', undefined, errorMessage));
  //     } finally {
  //       setTimeout(() => {
  //         actionConfig.stopProcessing(processingKey);
  //       }, 100);
  //     }
  //   },
  //   [
  //     cleanFunc,
  //     clearError,
  //     collection,
  //     createPayload,
  //     dataFacetRef,
  //     pageFunc,
  //     pagination.currentPage,
  //     pagination.pageSize,
  //     sort,
  //     filter,
  //     setPageData,
  //     setTotalItems,
  //     emitSuccess,
  //     handleError,
  //     failure,
  //     emitStatus,
  //     actionConfig,
  //   ],
  // );

  // Map action types to handlers
  const handlers = useMemo(
    () => ({
      handleAdd,
      handlePublish,
      handlePin,
      handleToggle,
      handleRemove,
      handleAutoname,
      handleUpdate,
      handleClean,
    }),
    [
      handleAdd,
      handlePublish,
      handlePin,
      handleToggle,
      handleRemove,
      handleAutoname,
      handleUpdate,
      handleClean,
      // TODO: Add handleCleanOne when implemented
      // handleCleanOne,
    ],
  );

  return {
    handlers,
    config: {
      headerActions,
      rowActions,
      collection,
      getCurrentDataFacet,
      isWalletConnected,
    },
  };
};
