import { ReactNode } from 'react';

import { act, renderHook } from '@testing-library/react';

import { ViewContextProvider, useViewContext } from '../ViewContext';
import { ViewStateKey } from '../index';

// Minimal wrapper for testing
const createWrapper = () => {
  const Wrapper = ({ children }: { children: ReactNode }) => (
    <ViewContextProvider>{children}</ViewContextProvider>
  );
  Wrapper.displayName = 'TestWrapper';
  return Wrapper;
};

// Helper to create mock ViewStateKey - using correct property names
const createViewStateKey = (view: string, tab: string): ViewStateKey => ({
  viewName: view,
  tabName: tab,
});

describe('ViewContext Object Identity', () => {
  test('getPagination returns same object reference for repeated calls with same viewStateKey', () => {
    const wrapper = createWrapper();
    const { result } = renderHook(() => useViewContext(), { wrapper });

    const viewStateKey = createViewStateKey('test-view', 'test-tab');

    // Call getPagination multiple times with same viewStateKey
    const pagination1 = result.current.getPagination(viewStateKey);
    const pagination2 = result.current.getPagination(viewStateKey);
    const pagination3 = result.current.getPagination(viewStateKey);

    // They should be the same object reference
    expect(pagination1).toBe(pagination2);
    expect(pagination2).toBe(pagination3);
    expect(pagination1).toBe(pagination3);
  });

  test('getPagination returns different objects for different viewKeys', () => {
    const wrapper = createWrapper();
    const { result } = renderHook(() => useViewContext(), { wrapper });

    const viewStateKey1 = createViewStateKey('view1', 'tab1');
    const viewStateKey2 = createViewStateKey('view2', 'tab2');

    // First, update one of the keys to create actual pagination data
    act(() => {
      result.current.updatePagination(viewStateKey1, { currentPage: 1 });
    });

    const pagination1 = result.current.getPagination(viewStateKey1); // Has data
    const pagination2 = result.current.getPagination(viewStateKey2); // Returns initial state

    // Different keys should return different objects (one has data, one is initial)
    expect(pagination1).not.toBe(pagination2);
    expect(pagination1.currentPage).toBe(1); // Updated value
    expect(pagination2.currentPage).toBe(0); // Initial value

    // But calling same key again should return same object
    const pagination1Again = result.current.getPagination(viewStateKey1);
    expect(pagination1).toBe(pagination1Again);
  });

  test('getPagination function reference is stable', () => {
    const wrapper = createWrapper();
    const { result, rerender } = renderHook(() => useViewContext(), {
      wrapper,
    });

    const getPagination1 = result.current.getPagination;

    // Force re-render
    rerender();

    const getPagination2 = result.current.getPagination;

    // Function reference should be stable
    expect(getPagination1).toBe(getPagination2);
  });

  test('updatePagination preserves object identity for unchanged keys', () => {
    const wrapper = createWrapper();
    const { result } = renderHook(() => useViewContext(), { wrapper });

    const viewStateKey1 = createViewStateKey('view1', 'tab1');
    const viewStateKey2 = createViewStateKey('view2', 'tab2');

    // Get initial pagination objects
    const pagination1Before = result.current.getPagination(viewStateKey1);
    const pagination2Before = result.current.getPagination(viewStateKey2);

    // Update pagination for viewStateKey1 only
    act(() => {
      result.current.updatePagination(viewStateKey1, { currentPage: 1 });
    });

    // Get pagination objects after update
    const pagination1After = result.current.getPagination(viewStateKey1);
    const pagination2After = result.current.getPagination(viewStateKey2);

    // viewKey1 should have new object (it was updated)
    expect(pagination1Before).not.toBe(pagination1After);
    expect(pagination1After.currentPage).toBe(1);

    // viewKey2 should have same object (it wasn't updated)
    expect(pagination2Before).toBe(pagination2After);
  });

  test('initial pagination state returns same object reference', () => {
    const wrapper = createWrapper();
    const { result } = renderHook(() => useViewContext(), { wrapper });

    const viewStateKey1 = createViewStateKey('new-view1', 'new-tab1');
    const viewStateKey2 = createViewStateKey('new-view2', 'new-tab2');

    // Get pagination for non-existent keys (should return initial state)
    const pagination1 = result.current.getPagination(viewStateKey1);
    const pagination2 = result.current.getPagination(viewStateKey2);

    // Should return same initial state object reference
    expect(pagination1).toBe(pagination2);
    expect(pagination1.currentPage).toBe(0);
    expect(pagination1.pageSize).toBe(15);
    expect(pagination1.totalItems).toBe(0);
  });
});
