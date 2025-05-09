import * as TableContext from '@components';
import { useTableKeys } from '@components';
import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const baseMockContext = {
  focusState: 'table' as const,
  selectedRowIndex: 1,
  setSelectedRowIndex: vi.fn(),
  focusTable: vi.fn(),
  tableRef: { current: null },
  focusControls: vi.fn(),
};

function mockEvent(
  key: string,
  extra: Partial<React.KeyboardEvent> = {},
): React.KeyboardEvent<HTMLTableElement> {
  return {
    key,
    preventDefault: vi.fn(),
    ...extra,
  } as unknown as React.KeyboardEvent<HTMLTableElement>;
}

describe('useTableKeys', () => {
  it('ArrowDown selects next row if not last', () => {
    const setSelectedRowIndex = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue({
      ...baseMockContext,
      selectedRowIndex: 1,
      setSelectedRowIndex,
    });
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 0,
        totalPages: 2,
        onPageChange: vi.fn(),
      }),
    );
    result.current.handleKeyDown(mockEvent('ArrowDown'));
    expect(setSelectedRowIndex).toHaveBeenCalledWith(2);
  });

  it('ArrowDown on last row goes to next page and selects first row', () => {
    const setSelectedRowIndex = vi.fn();
    const onPageChange = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue({
      ...baseMockContext,
      selectedRowIndex: 2,
      setSelectedRowIndex,
    });
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 0,
        totalPages: 2,
        onPageChange,
      }),
    );
    result.current.handleKeyDown(mockEvent('ArrowDown'));
    expect(onPageChange).toHaveBeenCalledWith(1);
    expect(setSelectedRowIndex).toHaveBeenCalledWith(0);
  });

  it('ArrowUp selects previous row if not first', () => {
    const setSelectedRowIndex = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue({
      ...baseMockContext,
      selectedRowIndex: 2,
      setSelectedRowIndex,
    });
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 0,
        totalPages: 2,
        onPageChange: vi.fn(),
      }),
    );
    result.current.handleKeyDown(mockEvent('ArrowUp'));
    expect(setSelectedRowIndex).toHaveBeenCalledWith(1);
  });

  it('ArrowUp on first row goes to previous page and selects -1', () => {
    const setSelectedRowIndex = vi.fn();
    const onPageChange = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue({
      ...baseMockContext,
      selectedRowIndex: 0,
      setSelectedRowIndex,
    });
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 1,
        totalPages: 2,
        onPageChange,
      }),
    );
    result.current.handleKeyDown(mockEvent('ArrowUp'));
    expect(onPageChange).toHaveBeenCalledWith(0);
    expect(setSelectedRowIndex).toHaveBeenCalledWith(-1);
  });

  it('ArrowLeft and PageUp go to previous page', () => {
    const onPageChange = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue(baseMockContext);
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 1,
        totalPages: 2,
        onPageChange,
      }),
    );
    result.current.handleKeyDown(mockEvent('ArrowLeft'));
    result.current.handleKeyDown(mockEvent('PageUp'));
    expect(onPageChange).toHaveBeenCalledWith(0);
  });

  it('ArrowRight and PageDown go to next page', () => {
    const onPageChange = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue(baseMockContext);
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 0,
        totalPages: 2,
        onPageChange,
      }),
    );
    result.current.handleKeyDown(mockEvent('ArrowRight'));
    result.current.handleKeyDown(mockEvent('PageDown'));
    expect(onPageChange).toHaveBeenCalledWith(1);
  });

  it('Home with ctrl/meta goes to first page and selects first row', () => {
    const setSelectedRowIndex = vi.fn();
    const onPageChange = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue({
      ...baseMockContext,
      setSelectedRowIndex,
    });
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 1,
        totalPages: 2,
        onPageChange,
      }),
    );
    result.current.handleKeyDown(mockEvent('Home', { ctrlKey: true }));
    result.current.handleKeyDown(mockEvent('Home', { metaKey: true }));
    expect(onPageChange).toHaveBeenCalledWith(0);
    expect(setSelectedRowIndex).toHaveBeenCalledWith(0);
  });

  it('Home without ctrl/meta selects first row', () => {
    const setSelectedRowIndex = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue({
      ...baseMockContext,
      setSelectedRowIndex,
    });
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 1,
        totalPages: 2,
        onPageChange: vi.fn(),
      }),
    );
    result.current.handleKeyDown(mockEvent('Home'));
    expect(setSelectedRowIndex).toHaveBeenCalledWith(0);
  });

  it('End with ctrl/meta goes to last page and selects last row', () => {
    const setSelectedRowIndex = vi.fn();
    const onPageChange = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue({
      ...baseMockContext,
      setSelectedRowIndex,
    });
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 0,
        totalPages: 2,
        onPageChange,
      }),
    );
    result.current.handleKeyDown(mockEvent('End', { ctrlKey: true }));
    result.current.handleKeyDown(mockEvent('End', { metaKey: true }));
    expect(onPageChange).toHaveBeenCalledWith(1);
    expect(setSelectedRowIndex).toHaveBeenCalledWith(2);
  });

  it('End without ctrl/meta selects last row', () => {
    const setSelectedRowIndex = vi.fn();
    vi.spyOn(TableContext, 'useTableContext').mockReturnValue({
      ...baseMockContext,
      setSelectedRowIndex,
    });
    const { result } = renderHook(() =>
      useTableKeys({
        itemCount: 3,
        currentPage: 1,
        totalPages: 2,
        onPageChange: vi.fn(),
      }),
    );
    result.current.handleKeyDown(mockEvent('End'));
    expect(setSelectedRowIndex).toHaveBeenCalledWith(2);
  });
});
