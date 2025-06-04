import { useCallback, useMemo } from 'react';

import { TableKey, useViewContext } from '@contexts';

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

  return useMemo(
    () => ({
      pagination,
      goToPage,
      changePageSize,
      setTotalItems,
    }),
    [pagination, goToPage, changePageSize, setTotalItems],
  );
};
