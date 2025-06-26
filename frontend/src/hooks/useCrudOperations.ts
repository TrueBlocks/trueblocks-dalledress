import { useCallback } from 'react';

import { crud, sdk, types } from '@models';
import { getAddressString } from '@utils';

// Generic types for different page data types
export interface PageData {
  totalItems: number;
  // Add other common properties that all page data types share
}

// Action type definition
type ActionType =
  | 'create'
  | 'update'
  | 'delete'
  | 'undelete'
  | 'remove'
  | 'autoname'
  | 'clean'
  | 'reload';

// Configuration interface for the hook - now with proper typing
export interface CrudOperationsConfig<TPageData extends PageData, TItem> {
  // State management
  pageData: TPageData | null;
  setPageData: React.Dispatch<React.SetStateAction<TPageData | null>>;
  setTotalItems: (total: number) => void;

  // API functions
  crudFunction: (
    payload: types.Payload,
    operation: crud.Operation,
    item: TItem,
  ) => Promise<void>;
  getPageFunction: (
    payload: types.Payload,
    offset: number,
    limit: number,
    sort: sdk.SortSpec,
    filter: string,
  ) => Promise<TPageData>;

  // Data facet (can be derived from current facet state)
  dataFacetRef: React.MutableRefObject<types.DataFacet>;

  // Pagination
  pagination: { currentPage: number; pageSize: number };
  sort: sdk.SortSpec;
  filter: string;
  goToPage: (page: number) => void;

  // Error handling and messaging
  clearError: () => void;
  handleError: (error: unknown, message: string) => void;
  emitSuccess: (operation: ActionType, address: string | number) => void;
  failure: (
    operation: ActionType,
    address?: string,
    message?: string,
  ) => string;

  // Action config
  actionConfig: {
    startProcessing: (address: string) => void;
    stopProcessing: (address: string) => void;
  };

  // Collection-specific configuration
  collectionName: string; // e.g., 'abis', 'monitors', 'names'
  itemsProperty: string; // e.g., 'abis', 'monitors', 'names'
  PageClass: new (data: Record<string, unknown>) => TPageData;
  emptyItem: TItem;

  // Payload creation hook
  createPayload: (
    dataFacet: types.DataFacet,
    address?: string,
  ) => types.Payload;
}

export const useCrudOperations = <TPageData extends PageData, TItem>(
  config: CrudOperationsConfig<TPageData, TItem>,
) => {
  const {
    pageData,
    setPageData,
    setTotalItems,
    crudFunction,
    getPageFunction,
    createPayload,
    dataFacetRef,
    pagination,
    sort,
    filter,
    goToPage,
    clearError,
    handleError,
    emitSuccess,
    failure,
    actionConfig,
    collectionName,
    itemsProperty,
    PageClass,
    emptyItem,
  } = config;

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
          return new PageClass({
            ...prev,
            [itemsProperty]: optimisticValues,
          }) as TPageData;
        });

        crudFunction(
          createPayload(dataFacetRef.current, address),
          crud.Operation.REMOVE,
          emptyItem,
        )
          .then(async () => {
            const result = await getPageFunction(
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
              return new PageClass({
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
      PageClass,
      crudFunction,
      createPayload,
      dataFacetRef,
      emptyItem,
      getPageFunction,
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

  return {
    handleRemove,
  };
};
