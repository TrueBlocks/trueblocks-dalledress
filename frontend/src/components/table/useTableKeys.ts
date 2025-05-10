import { useCallback } from 'react';

import { Logger } from '@app';
import { useTableContext } from '@components';

// UseTableKeysProps defines the props for the useTableKeys hook.
interface UseTableKeysProps {
  itemCount: number;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

// useTableKeys is a custom hook that provides keyboard navigation logic for the table, handling arrow keys, page navigation, and selection.
export const useTableKeys = ({
  itemCount,
  currentPage,
  totalPages,
  onPageChange,
}: UseTableKeysProps) => {
  const { focusState, selectedRowIndex, setSelectedRowIndex, focusTable } =
    useTableContext();

  // Debounce ref for page navigation

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (focusState !== 'table') return;
      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          if (selectedRowIndex < itemCount - 1) {
            setSelectedRowIndex(selectedRowIndex + 1);
          } else if (currentPage < totalPages - 1) {
            onPageChange(currentPage + 1);
            setSelectedRowIndex(0);
          }
          break;
        case 'ArrowUp':
          e.preventDefault();
          if (selectedRowIndex > 0) {
            setSelectedRowIndex(selectedRowIndex - 1);
          } else if (currentPage > 0) {
            onPageChange(currentPage - 1);
            setSelectedRowIndex(-1);
          }
          break;
        case 'ArrowLeft':
        case 'PageUp':
          e.preventDefault();
          if (currentPage > 0) {
            onPageChange(currentPage - 1);
          }
          break;
        case 'ArrowRight':
        case 'PageDown':
          e.preventDefault();
          if (currentPage < totalPages - 1) {
            onPageChange(currentPage + 1);
          }
          break;
        case 'Home':
          e.preventDefault();
          onPageChange(0);
          setSelectedRowIndex(0);
          break;
        case 'End':
          e.preventDefault();
          onPageChange(totalPages - 1);
          setSelectedRowIndex(itemCount - 1);
          break;
        case 'Enter':
          e.preventDefault();
          Logger('Table: Enter key pressed');
          break;
      }
    },
    [
      focusState,
      selectedRowIndex,
      itemCount,
      currentPage,
      totalPages,
      setSelectedRowIndex,
      onPageChange,
    ],
  );

  const requestFocus = useCallback(() => {
    if (focusTable) {
      focusTable();
    }
  }, [focusTable]);

  return {
    handleKeyDown,
    requestFocus,
  };
};
