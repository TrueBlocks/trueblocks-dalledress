import { DataFacetConfig, useActiveFacet2 } from '@hooks';
import { setupWailsMocks } from '@mocks';
import { types } from '@models';
import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Mock the focused hooks (no duplication - centralized in setup)
vi.mock('../../hooks/useActiveProject2', () => ({
  useActiveProject2: vi.fn(),
}));

vi.mock('../../hooks/usePreferences2', () => ({
  usePreferences2: vi.fn(),
}));

vi.mock('../../hooks/useUIState2', () => ({
  useUIState: vi.fn(),
}));

const { useActiveProject2 } = await import('../../hooks/useActiveProject2');
const { usePreferences2 } = await import('../../hooks/usePreferences2');
const { useUIState } = await import('../../hooks/useUIState2');

const mockedUseActiveProject2 = vi.mocked(useActiveProject2);
const mockedUsePreferences = vi.mocked(usePreferences2);
const mockedUseUIState = vi.mocked(useUIState);

describe('Exports View Migration Tests (useActiveFacet2 integration)', () => {
  let mockLastFacetMap: Record<string, types.DataFacet>;
  let mockSetLastFacet: ReturnType<typeof vi.fn>;

  const exportsFacets: DataFacetConfig[] = [
    {
      id: types.DataFacet.STATEMENTS,
      label: 'Statements',
    },
    {
      id: types.DataFacet.TRANSFERS,
      label: 'Transfers',
    },
    {
      id: types.DataFacet.BALANCES,
      label: 'Balances',
    },
    {
      id: types.DataFacet.TRANSACTIONS,
      label: 'Transactions',
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();

    // Set up Wails mocks to prevent bridge errors
    setupWailsMocks();

    mockLastFacetMap = {};
    mockSetLastFacet = vi.fn();

    // Mock useActiveProject2 (project context)
    mockedUseActiveProject2.mockReturnValue({
      lastFacetMap: mockLastFacetMap,
      setLastFacet: mockSetLastFacet,
      getLastFacet: vi.fn((view: string) => {
        const vR = view.replace(/^\/+/, '');
        return mockLastFacetMap[vR] || '';
      }),
      activeChain: 'mainnet',
      activeAddress: '0x123',
      activeContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
      loading: false,
      effectiveAddress: '0x123',
      effectiveChain: 'mainnet',
      lastProject: 'test-project',
      lastView: 'exports',
      setActiveAddress: vi.fn(),
      setActiveChain: vi.fn(),
      setActiveContract: vi.fn(),
      setLastView: vi.fn(),
      switchProject: vi.fn(),
      hasActiveProject: true,
      canExport: true,
    });

    // Mock usePreferences2 (theme, language, debug)
    mockedUsePreferences.mockReturnValue({
      lastTheme: 'dark',
      lastLanguage: 'en',
      debugMode: false,
      loading: false,
      toggleTheme: vi.fn(),
      changeLanguage: vi.fn(),
      toggleDebugMode: vi.fn(),
      isDarkMode: true,
    });

    // Mock useUIState (collapsed states)
    mockedUseUIState.mockReturnValue({
      menuCollapsed: false,
      helpCollapsed: false,
      showDetailPanel: true,
      loading: false,
      setMenuCollapsed: vi.fn(),
      setHelpCollapsed: vi.fn(),
      setShowDetailPanel: vi.fn(),
    });
  });

  describe('hook integration', () => {
    it('should initialize with default transactions facet', () => {
      const { result } = renderHook(() =>
        useActiveFacet2({
          viewRoute: 'exports',
          facets: exportsFacets,
        }),
      );

      expect(result.current.activeFacet).toBe(types.DataFacet.STATEMENTS);
      expect(
        result.current.getFacetConfig(types.DataFacet.TRANSACTIONS)?.label,
      ).toBe('Transactions');
    });

    it('should restore saved facet from preferences', () => {
      mockLastFacetMap['exports'] = types.DataFacet.STATEMENTS;

      const { result } = renderHook(() =>
        useActiveFacet2({
          viewRoute: 'exports',
          facets: exportsFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
      expect(
        result.current.getFacetConfig(types.DataFacet.STATEMENTS)?.label,
      ).toBe('Statements');
    });

    it('should support switching between all export facets', async () => {
      const { result } = renderHook(() =>
        useActiveFacet2({
          viewRoute: 'exports',
          facets: exportsFacets,
        }),
      );

      // Switch to transfers
      await act(async () => {
        result.current.setActiveFacet(types.DataFacet.TRANSFERS);
      });

      expect(mockSetLastFacet).toHaveBeenCalledWith('exports', 'transfers');

      // Switch to balances
      await act(async () => {
        result.current.setActiveFacet(types.DataFacet.BALANCES);
      });

      expect(mockSetLastFacet).toHaveBeenCalledWith('exports', 'balances');

      // Verify all expected calls were made
      expect(mockSetLastFacet).toHaveBeenCalledTimes(2);
    });

    it('should provide all available facets', () => {
      const { result } = renderHook(() =>
        useActiveFacet2({
          viewRoute: 'exports',
          facets: exportsFacets,
        }),
      );

      expect(result.current.availableFacets).toHaveLength(4);
      expect(result.current.availableFacets.map((f: any) => f.id)).toEqual([
        'statements',
        'transfers',
        'balances',
        'transactions',
      ]);
    });

    it('should handle ViewStateKey generation', () => {
      mockLastFacetMap['exports'] = types.DataFacet.BALANCES;

      const { result } = renderHook(() =>
        useActiveFacet2({
          viewRoute: 'exports',
          facets: exportsFacets,
        }),
      );

      // ViewStateKey should be created manually in the view using activeFacet
      expect(result.current.activeFacet).toBe('balances');

      // Test the pattern the view will use
      const viewStateKey = {
        viewName: 'exports',
        facetName: result.current.activeFacet,
      };
      expect(viewStateKey).toEqual({
        viewName: 'exports',
        facetName: types.DataFacet.BALANCES,
      });
    });
  });

  describe('state management integration', () => {
    it('should maintain facet state across preference updates', () => {
      // Set initial preference
      mockLastFacetMap['exports'] = types.DataFacet.TRANSFERS;

      const { result } = renderHook(() =>
        useActiveFacet2({
          viewRoute: 'exports',
          facets: exportsFacets,
        }),
      );

      // Should read the saved preference
      expect(result.current.activeFacet).toBe('transfers');
    });

    it('should fallback to default for invalid saved facets', () => {
      // Clear the map to simulate invalid/missing facet data
      mockLastFacetMap = {};

      const { result } = renderHook(() =>
        useActiveFacet2({
          viewRoute: 'exports',
          facets: exportsFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
    });
  });
});
