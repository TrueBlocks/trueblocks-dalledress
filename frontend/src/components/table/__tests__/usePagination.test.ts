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
    const getViewPagination = vi.fn().mockReturnValue({
      currentPage: 2,
      pageSize: 25,
      totalItems: 100,
    });
    const updateViewPagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getViewPagination,
      updateViewPagination,
    });

    const { result } = renderHook(() => usePagination('view', 'tab'));
    expect(result.current.pagination).toEqual({
      currentPage: 2,
      pageSize: 25,
      totalItems: 100,
    });
    expect(result.current.currentPage).toBe(2);
    expect(result.current.pageSize).toBe(25);
    expect(result.current.totalItems).toBe(100);
  });

  it('goToPage calls updateViewPagination', () => {
    const getViewPagination = vi
      .fn()
      .mockReturnValue({ currentPage: 0, pageSize: 10, totalItems: 10 });
    const updateViewPagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getViewPagination,
      updateViewPagination,
    });
    const { result } = renderHook(() => usePagination('view', 'tab'));
    act(() => {
      result.current.goToPage(3);
    });
    expect(updateViewPagination).toHaveBeenCalledWith('view', 'tab', {
      currentPage: 3,
    });
  });

  it('changePageSize calls updateViewPagination with currentPage 0', () => {
    const getViewPagination = vi
      .fn()
      .mockReturnValue({ currentPage: 0, pageSize: 10, totalItems: 10 });
    const updateViewPagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getViewPagination,
      updateViewPagination,
    });
    const { result } = renderHook(() => usePagination('view', 'tab'));
    act(() => {
      result.current.changePageSize(50);
    });
    expect(updateViewPagination).toHaveBeenCalledWith('view', 'tab', {
      currentPage: 0,
      pageSize: 50,
    });
  });

  it('setTotalItems calls updateViewPagination', () => {
    const getViewPagination = vi
      .fn()
      .mockReturnValue({ currentPage: 0, pageSize: 10, totalItems: 10 });
    const updateViewPagination = vi.fn();
    vi.spyOn(Contexts, 'useViewContext').mockReturnValue({
      ...baseMockContext,
      getViewPagination,
      updateViewPagination,
    });
    const { result } = renderHook(() => usePagination('view', 'tab'));
    act(() => {
      result.current.setTotalItems(123);
    });
    expect(updateViewPagination).toHaveBeenCalledWith('view', 'tab', {
      totalItems: 123,
    });
  });
});
