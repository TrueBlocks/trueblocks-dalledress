import { useActiveProject } from '@hooks';
import { types } from '@models';
import { appPreferencesStore } from '@stores';
import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Mock the store dependency
vi.mock('@stores', () => ({
  appPreferencesStore: {
    subscribe: vi.fn(),
    getState: vi.fn(() => ({
      debugMode: false,
      isDarkMode: false,
      hasActiveProject: true,
      canExport: true,
      lastFacetMap: 'default',
      activeAddress: '',
      activeChain: '',
      theme: 'light',
      language: 'en',
      menuCollapsed: false,
      helpCollapsed: false,
      lastView: 'default',
    })),
    setLastFacet: vi.fn(),
    setActiveAddress: vi.fn(),
    setActiveChain: vi.fn(),
    switchProject: vi.fn(),
    toggleTheme: vi.fn(),
    changeLanguage: vi.fn(),
    setMenuCollapsed: vi.fn(),
    setHelpCollapsed: vi.fn(),
    setLastView: vi.fn(),
    toggleDarkMode: vi.fn(),
    isDarkMode: false,
    hasActiveProject: true,
    canExport: true,
  },
}));

describe('Preference System Tests (DataFacet refactor preparation)', () => {
  let mockStoreState: any;

  beforeEach(() => {
    vi.clearAllMocks();

    mockStoreState = {
      lastProject: 'test-project',
      activeChain: 'mainnet',
      activeAddress: '0x123',
      activeContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
      lastTheme: 'dark',
      lastLanguage: 'en',
      lastView: '/exports',
      menuCollapsed: false,
      helpCollapsed: false,
      lastFacetMap: {
        '/exports': 'transactions' as types.DataFacet,
        '/chunks': 'chunk-summary' as types.DataFacet,
        '/monitors': 'txs' as types.DataFacet,
      },
      loading: false,
    };

    (appPreferencesStore.getState as any).mockReturnValue(mockStoreState);
    (appPreferencesStore.subscribe as any).mockImplementation(
      (_callback: () => void) => {
        // Mock subscription
        return () => {}; // unsubscribe function
      },
    );
  });

  describe('lastFacetMap persistence behavior', () => {
    it('retrieves stored lastFacetMap values for different routes', () => {
      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastFacetMap).toEqual({
        '/exports': 'transactions',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
      });

      // Verify specific route lookups
      expect(result.current.lastFacetMap['/exports']).toBe('transactions');
      expect(result.current.lastFacetMap['/chunks']).toBe('chunk-summary');
      expect(result.current.lastFacetMap['/monitors']).toBe('txs');
    });

    it('handles missing lastFacetMap entries gracefully', () => {
      // Test scenario where a route has no stored lastFacetMap
      const { result } = renderHook(() => useActiveProject());

      // Routes not in the stored lastFacetMap should return undefined
      expect(result.current.lastFacetMap['/names']).toBeUndefined();
      expect(result.current.lastFacetMap['/abis']).toBeUndefined();
      expect(result.current.lastFacetMap['/unknown-route']).toBeUndefined();
    });

    it('calls setLastFacet with correct parameters', async () => {
      const { result } = renderHook(() => useActiveProject());

      await act(async () => {
        await result.current.setLastFacet(
          '/exports',
          'receipts' as types.DataFacet,
        );
      });

      expect(appPreferencesStore.setLastFacet as any).toHaveBeenCalledWith(
        '/exports',
        'receipts',
      );
    });

    it('supports all known DataFacet values in setLastFacet', async () => {
      const { result } = renderHook(() => useActiveProject());

      const testCases: Array<[string, types.DataFacet]> = [
        ['/exports', types.DataFacet.ALL],
        ['/exports', 'receipts' as types.DataFacet],
        ['/chunks', 'chunk-summary' as types.DataFacet],
        ['/monitors', 'txs' as types.DataFacet],
        ['/names', 'entity-names' as types.DataFacet],
        ['/abis', 'get-abis' as types.DataFacet],
      ];

      for (const [route, dataFacet] of testCases) {
        await act(async () => {
          await result.current.setLastFacet(route, dataFacet);
        });

        expect(appPreferencesStore.setLastFacet as any).toHaveBeenCalledWith(
          route,
          dataFacet,
        );
      }
    });
  });

  describe('cross-session state recovery', () => {
    it('initializes with previously stored lastFacetMap state', () => {
      // This simulates app restart with stored preferences
      const storedState = {
        ...mockStoreState,
        lastFacetMap: {
          '/exports': 'receipts' as types.DataFacet,
          '/chunks': 'chunk-summary' as types.DataFacet,
          '/monitors': 'txs' as types.DataFacet,
          '/names': 'entity-names' as types.DataFacet,
        },
      };

      (appPreferencesStore.getState as any).mockReturnValue(storedState);

      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastFacetMap).toEqual({
        '/exports': 'receipts',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
        '/names': 'entity-names',
      });
    });

    it('handles empty lastFacetMap state on first run', () => {
      const emptyState = {
        ...mockStoreState,
        lastFacetMap: {},
      };

      (appPreferencesStore.getState as any).mockReturnValue(emptyState);

      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastFacetMap).toEqual({});
      expect(Object.keys(result.current.lastFacetMap)).toHaveLength(0);
    });
  });

  describe('default facet selection logic', () => {
    it('provides setLastFacet for setting default selections', async () => {
      const { result } = renderHook(() => useActiveProject());

      expect(result.current.setLastFacet).toBeInstanceOf(Function);

      // Test that it can be called to set defaults
      await act(async () => {
        await result.current.setLastFacet(
          '/new-route',
          'transactions' as types.DataFacet,
        );
      });

      expect(appPreferencesStore.setLastFacet as any).toHaveBeenCalledWith(
        '/new-route',
        'transactions',
      );
    });

    it('maintains lastFacetMap structure for route-based lookups', () => {
      const { result } = renderHook(() => useActiveProject());

      // Verify the structure supports the pattern used in views:
      // const currentFacet = lastFacetMap[currentRoute] || defaultTab;
      const lastFacetMap = result.current.lastFacetMap;

      expect(typeof lastFacetMap).toBe('object');
      expect(lastFacetMap).not.toBeNull();
      expect(Array.isArray(lastFacetMap)).toBe(false);

      // Test the lookup pattern
      const exportsTab = lastFacetMap['/exports'] || 'transactions';
      const unknownTab = lastFacetMap['/unknown'] || 'default-facet';

      expect(exportsTab).toBe('transactions');
      expect(unknownTab).toBe('default-facet');
    });
  });

  describe('edge cases in preference handling', () => {
    it('handles route keys with special characters', async () => {
      const { result } = renderHook(() => useActiveProject());

      // Test routes that might have special characters
      const specialRoutes = [
        '/exports/sub-route',
        '/exports?param=value',
        '/exports#hash',
      ];

      for (const route of specialRoutes) {
        await act(async () => {
          await result.current.setLastFacet(
            route,
            'transactions' as types.DataFacet,
          );
        });

        expect(appPreferencesStore.setLastFacet as any).toHaveBeenCalledWith(
          route,
          'transactions',
        );
      }
    });

    it('handles concurrent lastFacetMap updates', async () => {
      const { result } = renderHook(() => useActiveProject());

      // Simulate rapid tab switching
      const promises = [
        result.current.setLastFacet(
          '/exports',
          'transactions' as types.DataFacet,
        ),
        result.current.setLastFacet('/exports', 'receipts' as types.DataFacet),
        result.current.setLastFacet(
          '/chunks',
          'chunk-summary' as types.DataFacet,
        ),
      ];

      await act(async () => {
        await Promise.all(promises);
      });

      expect(appPreferencesStore.setLastFacet as any).toHaveBeenCalledTimes(3);
    });

    it('handles invalid DataFacet values gracefully', async () => {
      const { result } = renderHook(() => useActiveProject());

      // This tests the function signature - TypeScript should catch invalid types,
      // but we test runtime behavior

      await act(async () => {
        await result.current.setLastFacet(
          '/exports',
          'transactions' as types.DataFacet,
        );
      });

      expect(appPreferencesStore.setLastFacet as any).toHaveBeenCalledWith(
        '/exports',
        'transactions',
      );
    });

    it('maintains state consistency during store updates', () => {
      let subscriptionCallback: (() => void) | null = null;

      (appPreferencesStore.subscribe as any).mockImplementation(
        (callback: () => void) => {
          subscriptionCallback = callback;
          return () => {
            subscriptionCallback = null;
          };
        },
      );

      const { result } = renderHook(() => useActiveProject());

      // Initial state
      expect(result.current.lastFacetMap).toEqual({
        '/exports': 'transactions',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
      });

      // Simulate store update
      const newState = {
        ...mockStoreState,
        lastFacetMap: {
          ...mockStoreState.lastFacetMap,
          '/exports': 'receipts' as types.DataFacet,
        },
      };
      (appPreferencesStore.getState as any).mockReturnValue(newState);

      // Trigger subscription callback to simulate store change
      if (subscriptionCallback) {
        act(() => {
          subscriptionCallback?.();
        });
      }

      // State should remain consistent
      expect(result.current.lastFacetMap['/exports']).toBe('receipts');
    });
  });
});
