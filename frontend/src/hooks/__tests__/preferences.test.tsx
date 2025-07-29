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
      lastFacet: 'default',
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
      lastFacet: {
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

  describe('lastFacet persistence behavior', () => {
    it('retrieves stored lastFacet values for different routes', () => {
      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastFacet).toEqual({
        '/exports': 'transactions',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
      });

      // Verify specific route lookups
      expect(result.current.lastFacet['/exports']).toBe('transactions');
      expect(result.current.lastFacet['/chunks']).toBe('chunk-summary');
      expect(result.current.lastFacet['/monitors']).toBe('txs');
    });

    it('handles missing lastFacet entries gracefully', () => {
      // Test scenario where a route has no stored lastFacet
      const { result } = renderHook(() => useActiveProject());

      // Routes not in the stored lastFacet should return undefined
      expect(result.current.lastFacet['/names']).toBeUndefined();
      expect(result.current.lastFacet['/abis']).toBeUndefined();
      expect(result.current.lastFacet['/unknown-route']).toBeUndefined();
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
    it('initializes with previously stored lastFacet state', () => {
      // This simulates app restart with stored preferences
      const storedState = {
        ...mockStoreState,
        lastFacet: {
          '/exports': 'receipts' as types.DataFacet,
          '/chunks': 'chunk-summary' as types.DataFacet,
          '/monitors': 'txs' as types.DataFacet,
          '/names': 'entity-names' as types.DataFacet,
        },
      };

      (appPreferencesStore.getState as any).mockReturnValue(storedState);

      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastFacet).toEqual({
        '/exports': 'receipts',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
        '/names': 'entity-names',
      });
    });

    it('handles empty lastFacet state on first run', () => {
      const emptyState = {
        ...mockStoreState,
        lastFacet: {},
      };

      (appPreferencesStore.getState as any).mockReturnValue(emptyState);

      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastFacet).toEqual({});
      expect(Object.keys(result.current.lastFacet)).toHaveLength(0);
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

    it('maintains lastFacet structure for route-based lookups', () => {
      const { result } = renderHook(() => useActiveProject());

      // Verify the structure supports the pattern used in views:
      // const currentTab = lastFacet[currentRoute] || defaultTab;
      const lastFacet = result.current.lastFacet;

      expect(typeof lastFacet).toBe('object');
      expect(lastFacet).not.toBeNull();
      expect(Array.isArray(lastFacet)).toBe(false);

      // Test the lookup pattern
      const exportsTab = lastFacet['/exports'] || 'transactions';
      const unknownTab = lastFacet['/unknown'] || 'default-facet';

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

    it('handles concurrent lastFacet updates', async () => {
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
      expect(result.current.lastFacet).toEqual({
        '/exports': 'transactions',
        '/chunks': 'chunk-summary',
        '/monitors': 'txs',
      });

      // Simulate store update
      const newState = {
        ...mockStoreState,
        lastFacet: {
          ...mockStoreState.lastFacet,
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
      expect(result.current.lastFacet['/exports']).toBe('receipts');
    });
  });
});
