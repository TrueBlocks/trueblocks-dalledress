import { useCallback, useMemo } from 'react';

import { NavigationTarget, TableKey, useViewContext } from '@contexts';

import { calculateNavigationTarget } from './navigationUtils';

// usePagination is a custom hook that manages pagination state and handlers for a table view/tab.
export const usePagination = (tableKey: TableKey) => {
  const { getPagination, updatePagination } = useViewContext();
  const pagination = getPagination(tableKey);

  const goToPage = useCallback(
    (page: number) => {
      updatePagination(tableKey, { currentPage: page });
    },
    [tableKey, updatePagination],
  );

  const changePageSize = useCallback(
    (size: number) => {
      updatePagination(tableKey, {
        currentPage: 0,
        pageSize: size,
      });
    },
    [tableKey, updatePagination],
  );

  const setTotalItems = useCallback(
    (total: number) => {
      updatePagination(tableKey, { totalItems: total });
    },
    [tableKey, updatePagination],
  );

  // Calculate navigation target for post-deletion positioning
  const calculatePostDeletionTarget = useCallback(
    (
      deletingRowId: string,
      currentPageData: Record<string, unknown>[],
    ): NavigationTarget | null => {
      return calculateNavigationTarget(
        deletingRowId,
        currentPageData,
        pagination.currentPage,
        pagination.pageSize,
        pagination.totalItems,
      );
    },
    [pagination.currentPage, pagination.pageSize, pagination.totalItems],
  );

  // Apply navigation target (navigate to calculated position)
  const applyNavigationTarget = useCallback(
    (target: NavigationTarget | null) => {
      if (!target) return;

      switch (target.type) {
        case 'page':
          if (target.page !== undefined) {
            updatePagination(tableKey, {
              currentPage: target.page,
            });
          }
          break;
        case 'row':
          // For row navigation, we'll let the table component handle it automatically
          // by selecting the last available row when data changes
          break;
        case 'none':
          // No action needed
          break;
      }
    },
    [tableKey, updatePagination],
  );

  return useMemo(
    () => ({
      pagination,
      goToPage,
      changePageSize,
      setTotalItems,
      calculatePostDeletionTarget,
      applyNavigationTarget,
    }),
    [
      pagination,
      goToPage,
      changePageSize,
      setTotalItems,
      calculatePostDeletionTarget,
      applyNavigationTarget,
    ],
  );
};
