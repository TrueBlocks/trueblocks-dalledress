import * as Contexts from '@contexts';
import { act, renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

import { usePagination } from '../usePagination';

const baseMockContext = {
  currentView: '',
  setCurrentView: vi.fn(),
  viewPagination: {},
};

describe('usePagination', () => {
  it('returns pagination state and handlers', () => {
    const getPagination = vi.fn().mockReturnValue({
      currentPage: 2,
      pageSize: 25,
      totalItems: 100,
    });
    const updatePagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getPagination,
      updatePagination,
    });

    const tableKey = { viewName: 'view', tabName: 'tab' };
    const { result } = renderHook(() => usePagination(tableKey));
    expect(result.current.pagination).toEqual({
      currentPage: 2,
      pageSize: 25,
      totalItems: 100,
    });
    expect(result.current.pagination.currentPage).toBe(2);
    expect(result.current.pagination.pageSize).toBe(25);
    expect(result.current.pagination.totalItems).toBe(100);
  });

  it('goToPage calls updatePagination', () => {
    const getPagination = vi
      .fn()
      .mockReturnValue({ currentPage: 0, pageSize: 10, totalItems: 10 });
    const updatePagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getPagination,
      updatePagination,
    });

    const tableKey = { viewName: 'view', tabName: 'tab' };
    const { result } = renderHook(() => usePagination(tableKey));
    act(() => {
      result.current.goToPage(3);
    });
    expect(updatePagination).toHaveBeenCalledWith(tableKey, {
      currentPage: 3,
    });
  });

  it('changePageSize calls updatePagination with currentPage 0', () => {
    const getPagination = vi
      .fn()
      .mockReturnValue({ currentPage: 0, pageSize: 10, totalItems: 10 });
    const updatePagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getPagination,
      updatePagination,
    });

    const tableKey = { viewName: 'view', tabName: 'tab' };
    const { result } = renderHook(() => usePagination(tableKey));
    act(() => {
      result.current.changePageSize(50);
    });
    expect(updatePagination).toHaveBeenCalledWith(
      { viewName: 'view', tabName: 'tab' },
      { currentPage: 0, pageSize: 50 },
    );
  });

  it('setTotalItems calls updatePagination', () => {
    const getPagination = vi
      .fn()
      .mockReturnValue({ currentPage: 0, pageSize: 10, totalItems: 10 });
    const updatePagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getPagination,
      updatePagination,
    });

    const tableKey = { viewName: 'view', tabName: 'tab' };
    const { result } = renderHook(() => usePagination(tableKey));
    act(() => {
      result.current.setTotalItems(123);
    });
    expect(updatePagination).toHaveBeenCalledWith(
      { viewName: 'view', tabName: 'tab' },
      { totalItems: 123 },
    );
  });

  it('uses memoized callbacks', () => {
    const getPagination = vi
      .fn()
      .mockReturnValue({ currentPage: 0, pageSize: 10, totalItems: 10 });
    const updatePagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getPagination,
      updatePagination,
    });

    const tableKey = { viewName: 'view', tabName: 'tab' };
    const { result, rerender } = renderHook(() => usePagination(tableKey));

    const initialGoToPage = result.current.goToPage;
    const initialChangePageSize = result.current.changePageSize;
    const initialSetTotalItems = result.current.setTotalItems;

    rerender();

    expect(result.current.goToPage).toBe(initialGoToPage);
    expect(result.current.changePageSize).toBe(initialChangePageSize);
    expect(result.current.setTotalItems).toBe(initialSetTotalItems);
  });
});
