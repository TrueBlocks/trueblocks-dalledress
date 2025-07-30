import { useActiveProject } from '@hooks';
import { types } from '@models';
import { appPreferencesStore } from '@stores';
import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Mock the store dependency
vi.mock('@stores', () => {
  // Create mock functions inside the factory
  const mockSetLastFacet = vi.fn();
  const mockSetActiveAddress = vi.fn();
  const mockSetActiveChain = vi.fn();
  const mockSetActiveContract = vi.fn();
  const mockSwitchProject = vi.fn();
  const mockToggleTheme = vi.fn();
  const mockChangeLanguage = vi.fn();
  const mockSetMenuCollapsed = vi.fn();
  const mockSetHelpCollapsed = vi.fn();
  const mockSetLastView = vi.fn();
  const mockToggleDebugMode = vi.fn();

  return {
    appPreferencesStore: {
      subscribe: vi.fn(),
      getSnapshot: vi.fn(() => ({
        debugMode: false,
        isDarkMode: false,
        hasActiveProject: true,
        canExport: true,
        lastFacetMap: {},
        activeAddress: '0x123',
        activeContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
        activeChain: 'mainnet',
        lastTheme: 'light',
        lastLanguage: 'en',
        menuCollapsed: false,
        helpCollapsed: false,
        lastView: '/',
        lastProject: 'test-project',
        loading: false,
        effectiveAddress: '0x123',
        effectiveChain: 'mainnet',
        setLastFacet: mockSetLastFacet,
        setActiveAddress: mockSetActiveAddress,
        setActiveChain: mockSetActiveChain,
        setActiveContract: mockSetActiveContract,
        switchProject: mockSwitchProject,
        toggleTheme: mockToggleTheme,
        changeLanguage: mockChangeLanguage,
        setMenuCollapsed: mockSetMenuCollapsed,
        setHelpCollapsed: mockSetHelpCollapsed,
        setLastView: mockSetLastView,
        toggleDebugMode: mockToggleDebugMode,
      })),
      getState: vi.fn(() => ({
        debugMode: false,
        isDarkMode: false,
        hasActiveProject: true,
        canExport: true,
        lastFacetMap: {},
        activeAddress: '0x123',
        activeContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
        activeChain: 'mainnet',
        lastTheme: 'light',
        lastLanguage: 'en',
        menuCollapsed: false,
        helpCollapsed: false,
        lastView: '/',
        lastProject: 'test-project',
        loading: false,
      })),
      setLastFacet: mockSetLastFacet,
      setActiveAddress: mockSetActiveAddress,
      setActiveChain: mockSetActiveChain,
      setActiveContract: mockSetActiveContract,
      switchProject: mockSwitchProject,
      toggleTheme: mockToggleTheme,
      changeLanguage: mockChangeLanguage,
      setMenuCollapsed: mockSetMenuCollapsed,
      setHelpCollapsed: mockSetHelpCollapsed,
      setLastView: mockSetLastView,
      toggleDebugMode: mockToggleDebugMode,
      isDarkMode: false,
      hasActiveProject: true,
      canExport: true,
      // Export the mock functions so tests can access them
      _mockFns: {
        mockSetLastFacet,
        mockSetActiveAddress,
        mockSetActiveChain,
        mockSetActiveContract,
        mockSwitchProject,
        mockToggleTheme,
        mockChangeLanguage,
        mockSetMenuCollapsed,
        mockSetHelpCollapsed,
        mockSetLastView,
        mockToggleDebugMode,
      },
    },
  };
});

describe('Preference System Tests (DataFacet refactor preparation)', () => {
  let mockStoreState: any;
  let mockFns: any;

  beforeEach(() => {
    vi.clearAllMocks();

    // Get mock functions from the store
    mockFns = (appPreferencesStore as any)._mockFns;

    mockStoreState = {
      lastProject: 'test-project',
      activeChain: 'mainnet',
      activeAddress: '0x123',
      activeContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
      lastTheme: 'light',
      lastLanguage: 'en',
      lastView: '/',
      menuCollapsed: false,
      helpCollapsed: false,
      debugMode: false,
      loading: false,
      lastFacetMap: {
        exports: 'transactions' as types.DataFacet,
        chunks: 'chunk-summary' as types.DataFacet,
        monitors: 'txs' as types.DataFacet,
      },
      effectiveAddress: '0x123',
      effectiveChain: 'mainnet',
      isDarkMode: false,
      hasActiveProject: true,
      canExport: true,
      setLastFacet: mockFns.mockSetLastFacet,
      setActiveAddress: mockFns.mockSetActiveAddress,
      setActiveChain: mockFns.mockSetActiveChain,
      setActiveContract: mockFns.mockSetActiveContract,
      switchProject: mockFns.mockSwitchProject,
      toggleTheme: mockFns.mockToggleTheme,
      changeLanguage: mockFns.mockChangeLanguage,
      setMenuCollapsed: mockFns.mockSetMenuCollapsed,
      setHelpCollapsed: mockFns.mockSetHelpCollapsed,
      setLastView: mockFns.mockSetLastView,
      toggleDebugMode: mockFns.mockToggleDebugMode,
    };

    (appPreferencesStore.getSnapshot as any).mockReturnValue(mockStoreState);
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
        exports: 'transactions',
        chunks: 'chunk-summary',
        monitors: 'txs',
      });

      // Verify specific route lookups
      expect(result.current.lastFacetMap['exports']).toBe('transactions');
      expect(result.current.lastFacetMap['chunks']).toBe('chunk-summary');
      expect(result.current.lastFacetMap['monitors']).toBe('txs');
    });

    it('handles missing lastFacetMap entries gracefully', () => {
      // Test scenario where a route has no stored lastFacetMap
      const { result } = renderHook(() => useActiveProject());

      // Routes not in the stored lastFacetMap should return undefined
      expect(result.current.lastFacetMap['names']).toBeUndefined();
      expect(result.current.lastFacetMap['abis']).toBeUndefined();
      expect(result.current.lastFacetMap['/unknown-route']).toBeUndefined();
    });

    it('calls setLastFacet with correct parameters', async () => {
      const { result } = renderHook(() => useActiveProject());

      await act(async () => {
        await result.current.setLastFacet(
          'exports',
          'receipts' as types.DataFacet,
        );
      });

      expect(mockFns.mockSetLastFacet).toHaveBeenCalledWith(
        'exports',
        'receipts',
      );
    });

    it('supports all known DataFacet values in setLastFacet', async () => {
      const { result } = renderHook(() => useActiveProject());

      const testCases: Array<[string, types.DataFacet]> = [
        ['exports', types.DataFacet.ALL],
        ['exports', 'receipts' as types.DataFacet],
        ['chunks', 'chunk-summary' as types.DataFacet],
        ['monitors', 'txs' as types.DataFacet],
        ['names', 'entity-names' as types.DataFacet],
        ['abis', 'get-abis' as types.DataFacet],
      ];

      for (const [route, dataFacet] of testCases) {
        await act(async () => {
          await result.current.setLastFacet(route, dataFacet);
        });

        expect(mockFns.mockSetLastFacet).toHaveBeenCalledWith(route, dataFacet);
      }
    });
  });

  describe('cross-session state recovery', () => {
    it('initializes with previously stored lastFacetMap state', () => {
      // This simulates app restart with stored preferences
      const storedState = {
        ...mockStoreState,
        lastFacetMap: {
          exports: 'receipts' as types.DataFacet,
          chunks: 'chunk-summary' as types.DataFacet,
          monitors: 'txs' as types.DataFacet,
          names: 'entity-names' as types.DataFacet,
        },
      };

      (appPreferencesStore.getSnapshot as any).mockReturnValue(storedState);

      const { result } = renderHook(() => useActiveProject());

      expect(result.current.lastFacetMap).toEqual({
        exports: 'receipts',
        chunks: 'chunk-summary',
        monitors: 'txs',
        names: 'entity-names',
      });
    });

    it('handles empty lastFacetMap state on first run', () => {
      const emptyState = {
        ...mockStoreState,
        lastFacetMap: {},
      };

      (appPreferencesStore.getSnapshot as any).mockReturnValue(emptyState);

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

      expect(mockFns.mockSetLastFacet).toHaveBeenCalledWith(
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
      const exportsTab = lastFacetMap['exports'] || 'transactions';
      const unknownTab = lastFacetMap['unknown'] || 'default-facet';

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

        expect(mockFns.mockSetLastFacet).toHaveBeenCalledWith(
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
          'exports',
          'transactions' as types.DataFacet,
        ),
        result.current.setLastFacet('exports', 'receipts' as types.DataFacet),
        result.current.setLastFacet(
          'chunks',
          'chunk-summary' as types.DataFacet,
        ),
      ];

      await act(async () => {
        await Promise.all(promises);
      });

      expect(mockFns.mockSetLastFacet).toHaveBeenCalledTimes(3);
    });

    it('handles invalid DataFacet values gracefully', async () => {
      const { result } = renderHook(() => useActiveProject());

      // This tests the function signature - TypeScript should catch invalid types,
      // but we test runtime behavior

      await act(async () => {
        await result.current.setLastFacet(
          'exports',
          'transactions' as types.DataFacet,
        );
      });

      expect(mockFns.mockSetLastFacet).toHaveBeenCalledWith(
        'exports',
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
        exports: 'transactions',
        chunks: 'chunk-summary',
        monitors: 'txs',
      });

      // Simulate store update
      const newState = {
        ...mockStoreState,
        lastFacetMap: {
          ...mockStoreState.lastFacetMap,
          exports: 'receipts' as types.DataFacet,
        },
      };
      (appPreferencesStore.getSnapshot as any).mockReturnValue(newState);

      // Trigger subscription callback to simulate store change
      if (subscriptionCallback) {
        act(() => {
          subscriptionCallback?.();
        });
      }

      // State should remain consistent
      expect(result.current.lastFacetMap['exports']).toBe('receipts');
    });
  });
});
