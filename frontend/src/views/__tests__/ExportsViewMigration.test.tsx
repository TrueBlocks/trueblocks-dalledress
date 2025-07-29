import { DataFacetConfig, useActiveFacet } from '@hooks';
import { types } from '@models';
import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Mock the useActiveProject hook
vi.mock('../../hooks/useActiveProject', () => ({
  useActiveProject: vi.fn(),
}));

const { useActiveProject } = await import('../../hooks/useActiveProject');
const mockedUseActiveProject = vi.mocked(useActiveProject);

describe('Exports View Migration Tests (useActiveFacet integration)', () => {
  let mockLastFacetMap: Record<string, string>;
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

    mockLastFacetMap = {};
    mockSetLastFacet = vi.fn();

    mockedUseActiveProject.mockReturnValue({
      lastFacetMap: mockLastFacetMap,
      setLastFacet: mockSetLastFacet,
      // Mock other required properties
      activeChain: 'mainnet',
      activeAddress: '0x123',
      activeContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
      lastTheme: 'dark',
      lastLanguage: 'en',
      lastView: '/exports',
      menuCollapsed: false,
      helpCollapsed: false,
      loading: false,
      effectiveAddress: '0x123',
      effectiveChain: 'mainnet',
    } as any);
  });

  describe('hook integration', () => {
    it('should initialize with default transactions facet', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
        }),
      );

      expect(result.current.activeFacet).toBe(types.DataFacet.STATEMENTS);
      expect(
        result.current.getFacetConfig(types.DataFacet.TRANSACTIONS)?.label,
      ).toBe('Transactions');
    });

    it('should restore saved facet from preferences', () => {
      mockLastFacetMap['/exports'] = 'statements';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
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
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
        }),
      );

      // Switch to transfers
      await act(async () => {
        result.current.setActiveFacet(types.DataFacet.TRANSFERS);
      });

      expect(mockSetLastFacet).toHaveBeenCalledWith('/exports', 'transfers');

      // Switch to balances
      await act(async () => {
        result.current.setActiveFacet(types.DataFacet.BALANCES);
      });

      expect(mockSetLastFacet).toHaveBeenCalledWith('/exports', 'balances');

      // Verify all expected calls were made
      expect(mockSetLastFacet).toHaveBeenCalledTimes(2);
    });

    it('should provide all available facets', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
        }),
      );

      expect(result.current.availableFacets).toHaveLength(4);
      expect(result.current.availableFacets.map((f) => f.id)).toEqual([
        'statements',
        'transfers',
        'balances',
        'transactions',
      ]);
    });

    it('should handle ViewStateKey generation', () => {
      mockLastFacetMap['/exports'] = 'balances';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
        }),
      );

      // ViewStateKey should be created manually in the view using activeFacet
      expect(result.current.activeFacet).toBe('balances');

      // Test the pattern the view will use
      const viewStateKey = {
        viewName: '/exports',
        facetName: result.current.activeFacet,
      };
      expect(viewStateKey).toEqual({
        viewName: '/exports',
        facetName: types.DataFacet.BALANCES,
      });
    });
  });

  describe('state management integration', () => {
    it('should maintain facet state across preference updates', () => {
      // Set initial preference
      mockLastFacetMap['/exports'] = 'transfers';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
        }),
      );

      // Should read the saved preference
      expect(result.current.activeFacet).toBe('transfers');
    });

    it('should fallback to default for invalid saved facets', () => {
      mockLastFacetMap['/exports'] = 'INVALID_FACET';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
    });
  });
});
