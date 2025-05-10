import { useCallback } from 'react';

import { useViewContext } from '@contexts';

// usePagination is a custom hook that manages pagination state and handlers for a table view/tab.
export const usePagination = (viewName: string, tabName: string) => {
  const { getViewPagination, updateViewPagination } = useViewContext();
  const pagination = getViewPagination(viewName, tabName);
  const goToPage = useCallback(
    (page: number) => {
      updateViewPagination(viewName, tabName, { currentPage: page });
    },
    [viewName, tabName, updateViewPagination],
  );
  const changePageSize = useCallback(
    (size: number) => {
      updateViewPagination(viewName, tabName, {
        currentPage: 0,
        pageSize: size,
      });
    },
    [viewName, tabName, updateViewPagination],
  );
  const setTotalItems = useCallback(
    (total: number) => {
      updateViewPagination(viewName, tabName, { totalItems: total });
    },
    [viewName, tabName, updateViewPagination],
  );
  return {
    pagination,
    goToPage,
    changePageSize,
    setTotalItems,
  };
};
