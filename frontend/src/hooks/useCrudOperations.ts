import { useCallback, useMemo } from 'react';

import { usePagination } from '@components';
import { ViewStateKey, useFiltering, useSorting } from '@contexts';
import { crud, sdk, types } from '@models';
import { getAddressString, useErrorHandler } from '@utils';

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
  getCurrentDataFacet: () => string;

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

  // Data facet (keeping for now)
  dataFacetRef: React.MutableRefObject<types.DataFacet>;

  // Action config (keeping for now)
  actionConfig: {
    startProcessing: (address: string) => void;
    stopProcessing: (address: string) => void;
  };

  // Collection-specific classes
  PageClass: new (data: Record<string, unknown>) => TPageData;
  emptyItem: TItem;
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
    crudFunction,
    getPageFunction,
    dataFacetRef,
    actionConfig,
    PageClass,
    emptyItem,
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
