import { useCallback } from 'react';

import { Logger } from '@app';
import { useTableContext } from '@components';
import { TableKey } from '@contexts';

import { usePagination } from './usePagination';

// UseTableKeysProps defines the props for the useTableKeys hook.
interface UseTableKeysProps {
  itemCount: number;
  currentPage: number;
  totalPages: number;
  tableKey: TableKey;
  expandedRowIndex?: number | null;
  setExpandedRowIndex?: (idx: number | null) => void;
  onEnter?: () => void;
  onEscape?: () => void;
}

// useTableKeys is a custom hook that provides keyboard navigation logic for the table, handling arrow keys, page navigation, and selection.
export const useTableKeys = ({
  itemCount,
  currentPage,
  totalPages,
  tableKey,
  expandedRowIndex = null,
  setExpandedRowIndex = () => {},
  onEnter,
  onEscape,
}: UseTableKeysProps) => {
  const { focusState, selectedRowIndex, setSelectedRowIndex, focusTable } =
    useTableContext();
  const { goToPage } = usePagination(tableKey);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (focusState !== 'table') return;
      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          if (selectedRowIndex < itemCount - 1) {
            setSelectedRowIndex(selectedRowIndex + 1);
          } else if (currentPage < totalPages - 1) {
            goToPage(currentPage + 1);
            setSelectedRowIndex(0);
          }
          break;
        case 'ArrowUp':
          e.preventDefault();
          if (selectedRowIndex > 0) {
            setSelectedRowIndex(selectedRowIndex - 1);
          } else if (currentPage > 0) {
            goToPage(currentPage - 1);
            setSelectedRowIndex(-1);
          }
          break;
        case 'ArrowLeft':
        case 'PageUp':
          e.preventDefault();
          if (currentPage > 0) {
            goToPage(currentPage - 1);
          }
          break;
        case 'ArrowRight':
        case 'PageDown':
          e.preventDefault();
          if (currentPage < totalPages - 1) {
            goToPage(currentPage + 1);
          }
          break;
        case 'Home':
          e.preventDefault();
          if (currentPage !== 0) {
            goToPage(0);
          }
          setSelectedRowIndex(0);
          break;
        case 'End':
          e.preventDefault();
          if (currentPage !== totalPages - 1) {
            goToPage(totalPages - 1);
          }
          setSelectedRowIndex(itemCount - 1);
          break;
        case 'Enter':
          e.preventDefault();
          if (onEnter) onEnter();
          // If setExpandedRowIndex is the default no-op, call Logger for legacy/test compatibility
          if (setExpandedRowIndex.toString() === '() => {}') {
            Logger(
              `Table ${tableKey.viewName}/${tableKey.tabName}: Enter key pressed`,
            );
          } else {
            if (expandedRowIndex === selectedRowIndex) {
              setExpandedRowIndex(null); // collapse if already open
            } else {
              setExpandedRowIndex(selectedRowIndex); // expand selected row
              // (scrolling logic removed)
            }
          }
          break;
        case 'Escape':
          e.preventDefault();
          if (onEscape) onEscape();
          if (expandedRowIndex !== null) {
            setExpandedRowIndex(null);
          }
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
      goToPage,
      tableKey,
      expandedRowIndex,
      setExpandedRowIndex,
      onEnter,
      onEscape,
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
