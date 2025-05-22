// Now import modules
import { useTableContext } from '@components';
import { TableKey } from '@contexts';
import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { usePagination } from '../usePagination';
import { useTableKeys } from '../useTableKeys';

// Mock dependencies before imports
vi.mock('@components', () => ({
  useTableContext: vi.fn(),
}));

vi.mock('../usePagination', () => ({
  usePagination: vi.fn(),
}));

// Mock the Log function
vi.mock('@utils', () => ({
  Log: vi.fn(),
}));

// Helper function to create mock keyboard events
function mockEvent(
  key: string,
  extra: Partial<React.KeyboardEvent> = {},
): React.KeyboardEvent {
  return {
    key,
    preventDefault: vi.fn(),
    ...extra,
  } as unknown as React.KeyboardEvent;
}

describe('useTableKeys', () => {
  const mockTableContext = {
    focusState: 'table',
    selectedRowIndex: 1,
    setSelectedRowIndex: vi.fn(),
    focusTable: vi.fn(),
    focusControls: vi.fn(),
    tableRef: { current: null },
  };

  const mockGoToPage = vi.fn();
  const tableKey: TableKey = { viewName: 'test-view', tabName: 'test-tab' };

  // Setup for each test
  beforeEach(() => {
    vi.clearAllMocks();
    (useTableContext as unknown as ReturnType<typeof vi.fn>).mockReturnValue(
      mockTableContext,
    );
    (usePagination as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      pagination: { currentPage: 0, pageSize: 10, totalItems: 100 },
      goToPage: mockGoToPage,
    });
  });

  // Group 1: Focus state behavior
  describe('Focus state behavior', () => {
    it('should do nothing when focusState is not "table"', () => {
      (useTableContext as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
        ...mockTableContext,
        focusState: 'controls',
      });

      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('ArrowDown'));
      });

      expect(mockTableContext.setSelectedRowIndex).not.toHaveBeenCalled();
      expect(mockGoToPage).not.toHaveBeenCalled();
    });

    it('should call focusTable when requestFocus is called', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.requestFocus();
      });

      expect(mockTableContext.focusTable).toHaveBeenCalled();
    });
  });

  // Group 2: Vertical navigation (ArrowUp/ArrowDown) tests
  describe('Vertical navigation', () => {
    it('should handle ArrowDown key - within items', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('ArrowDown'));
      });

      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(2);
      expect(mockGoToPage).not.toHaveBeenCalled();
    });

    it('should handle ArrowDown key - navigate to next page', () => {
      mockTableContext.selectedRowIndex = 4;

      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('ArrowDown'));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(1);
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(0);
    });

    it('should handle ArrowUp key - within items', () => {
      mockTableContext.selectedRowIndex = 2;

      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('ArrowUp'));
      });

      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(1);
      expect(mockGoToPage).not.toHaveBeenCalled();
    });

    it('should handle ArrowUp key - navigate to previous page', () => {
      mockTableContext.selectedRowIndex = 0;

      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('ArrowUp'));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(0);
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(-1);
    });

    it('should not navigate down at last row of last page', () => {
      mockTableContext.selectedRowIndex = 4;

      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1, // Last page
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('ArrowDown'));
      });

      expect(mockGoToPage).not.toHaveBeenCalled();
      expect(mockTableContext.setSelectedRowIndex).not.toHaveBeenCalled();
    });

    it('should not navigate up at first row of first page', () => {
      mockTableContext.selectedRowIndex = 0;

      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0, // First page
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('ArrowUp'));
      });

      expect(mockGoToPage).not.toHaveBeenCalled();
      expect(mockTableContext.setSelectedRowIndex).not.toHaveBeenCalled();
    });
  });

  // Group 3: Page navigation (PageUp/PageDown) tests
  describe('Page navigation', () => {
    it('should handle PageUp key', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('PageUp'));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(0);
    });

    it('should handle PageDown key', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('PageDown'));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(1);
    });

    it('should not navigate to previous page when already on the first page', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0, // First page
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('PageUp'));
      });

      expect(mockGoToPage).not.toHaveBeenCalled();
    });

    it('should not navigate to next page when already on the last page', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1, // Last page
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('PageDown'));
      });

      expect(mockGoToPage).not.toHaveBeenCalled();
    });
  });

  // Group 4: Home/End navigation tests
  describe('Home/End navigation', () => {
    it('should handle Home key', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('Home'));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(0);
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(0);
    });

    it('should handle Home key with Ctrl modifier', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('Home', { ctrlKey: true }));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(0);
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(0);
    });

    it('should handle Home key with Meta modifier', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('Home', { metaKey: true }));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(0);
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(0);
    });

    it('should handle End key', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('End'));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(1);
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(4);
    });

    it('should handle End key with Ctrl modifier', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('End', { ctrlKey: true }));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(1);
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(4);
    });

    it('should handle End key with Meta modifier', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('End', { metaKey: true }));
      });

      expect(mockGoToPage).toHaveBeenCalledWith(1);
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(4);
    });

    it('should not change page when Home is pressed on first page', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0, // Already on first page
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('Home'));
      });

      expect(mockGoToPage).not.toHaveBeenCalled();
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(0);
    });

    it('should not change page when End is pressed on last page', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 1, // Already on last page
          totalPages: 2,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('End'));
      });

      expect(mockGoToPage).not.toHaveBeenCalled();
      expect(mockTableContext.setSelectedRowIndex).toHaveBeenCalledWith(4);
    });
  });

  describe('Special cases', () => {
    it('should handle empty data set correctly', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 0, // Empty table
          currentPage: 0,
          totalPages: 1,
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('ArrowDown'));
      });

      expect(mockTableContext.setSelectedRowIndex).not.toHaveBeenCalled();
      expect(mockGoToPage).not.toHaveBeenCalled();
    });

    it('should handle single page dataset correctly', () => {
      const { result } = renderHook(() =>
        useTableKeys({
          itemCount: 5,
          currentPage: 0,
          totalPages: 1, // Only one page
          tableKey,
        }),
      );

      act(() => {
        result.current.handleKeyDown(mockEvent('PageDown'));
      });

      expect(mockGoToPage).not.toHaveBeenCalled();
    });
  });
});
