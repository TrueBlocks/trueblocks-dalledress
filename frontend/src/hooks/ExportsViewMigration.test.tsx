import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { useActiveFacet } from '../../hooks/useActiveFacet';
import { DataFacet, DataFacetConfig } from '../../hooks/useActiveFacet.types';

// Mock the useActiveProject hook
vi.mock('../../hooks/useActiveProject', () => ({
  useActiveProject: vi.fn(),
}));

const { useActiveProject } = await import('../../hooks/useActiveProject');
const mockedUseActiveProject = vi.mocked(useActiveProject);

describe('Exports View Migration Tests (useActiveFacet integration)', () => {
  let mockLastTab: Record<string, string>;
  let mockSetLastTab: ReturnType<typeof vi.fn>;

  const exportsFacets: DataFacetConfig[] = [
    {
      id: 'statements' as DataFacet,
      label: 'Statements',
      listKind: 'STATEMENTS' as any,
    },
    {
      id: 'transfers' as DataFacet,
      label: 'Transfers',
      listKind: 'TRANSFERS' as any,
    },
    {
      id: 'balances' as DataFacet,
      label: 'Balances',
      listKind: 'BALANCES' as any,
    },
    {
      id: 'transactions' as DataFacet,
      label: 'Transactions',
      listKind: 'TRANSACTIONS' as any,
      isDefault: true,
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();

    mockLastTab = {};
    mockSetLastTab = vi.fn();

    mockedUseActiveProject.mockReturnValue({
      lastTab: mockLastTab,
      setLastTab: mockSetLastTab,
      // Mock other required properties
      lastProject: 'test-project',
      lastChain: 'mainnet',
      lastAddress: '0x123',
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
          defaultFacet: 'transactions',
        }),
      );

      expect(result.current.activeFacet).toBe('transactions');
      expect(result.current.getFacetConfig('transactions')?.label).toBe(
        'Transactions',
      );
      expect(result.current.getFacetConfig('transactions')?.isDefault).toBe(
        true,
      );
    });

    it('should restore saved facet from preferences', () => {
      mockLastTab['/exports'] = 'statements';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
          defaultFacet: 'transactions',
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
      expect(result.current.getFacetConfig('statements')?.label).toBe(
        'Statements',
      );
    });

    it('should support switching between all export facets', async () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
          defaultFacet: 'transactions',
        }),
      );

      // Switch to transfers
      await act(async () => {
        result.current.setActiveFacet('transfers');
      });

      expect(mockSetLastTab).toHaveBeenCalledWith('/exports', 'TRANSFERS');

      // Switch to balances
      await act(async () => {
        result.current.setActiveFacet('balances');
      });

      expect(mockSetLastTab).toHaveBeenCalledWith('/exports', 'BALANCES');

      // Verify all expected calls were made
      expect(mockSetLastTab).toHaveBeenCalledTimes(2);
    });

    it('should provide correct ListKind backward compatibility', () => {
      mockLastTab['/exports'] = 'statements';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
          defaultFacet: 'transactions',
        }),
      );

      expect(result.current.getCurrentListKind()).toBe('STATEMENTS');
    });

    it('should provide all available facets', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
          defaultFacet: 'transactions',
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
      mockLastTab['/exports'] = 'balances';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
          defaultFacet: 'transactions',
        }),
      );

      // ViewStateKey should be created manually in the view using activeFacet
      expect(result.current.activeFacet).toBe('balances');

      // Test the pattern the view will use
      const viewStateKey = {
        viewName: '/exports',
        tabName: result.current.activeFacet,
      };
      expect(viewStateKey).toEqual({
        viewName: '/exports',
        tabName: 'balances',
      });
    });
  });

  describe('state management integration', () => {
    it('should maintain facet state across preference updates', () => {
      // Set initial preference
      mockLastTab['/exports'] = 'transfers';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
          defaultFacet: 'transactions',
        }),
      );

      // Should read the saved preference
      expect(result.current.activeFacet).toBe('transfers');
    });

    it('should fallback to default for invalid saved facets', () => {
      mockLastTab['/exports'] = 'invalid-facet';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: exportsFacets,
          defaultFacet: 'transactions',
        }),
      );

      expect(result.current.activeFacet).toBe('transactions');
    });
  });
});
