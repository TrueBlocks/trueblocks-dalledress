import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { useActiveFacet } from '../useActiveFacet';
import { DataFacet, DataFacetConfig } from '../useActiveFacet.types';

// Mock the useActiveProject hook
vi.mock('../useActiveProject', () => ({
  useActiveProject: vi.fn(),
}));

const { useActiveProject } = await import('../useActiveProject');
const mockedUseActiveProject = vi.mocked(useActiveProject);

describe('useActiveFacet Hook Tests (DataFacet implementation)', () => {
  let mockLastTab: Record<string, string>;
  let mockSetLastTab: ReturnType<typeof vi.fn>;

  const sampleFacets: DataFacetConfig[] = [
    {
      id: 'transactions' as DataFacet,
      label: 'Transactions',
      listKind: 'TRANSACTIONS' as any,
      isDefault: true,
    },
    {
      id: 'receipts' as DataFacet,
      label: 'Receipts',
      listKind: 'RECEIPTS' as any,
    },
    {
      id: 'statements' as DataFacet,
      label: 'Statements',
      listKind: 'STATEMENTS' as any,
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

  describe('initialization and defaults', () => {
    it('should initialize with default facet when no preference saved', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('transactions');
      expect(result.current.getDefaultFacet()).toBe('transactions');
    });

    it('should use provided default facet override', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
          defaultFacet: 'statements' as DataFacet,
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
    });

    it('should restore saved preference via ListKind mapping', () => {
      mockLastTab['/exports'] = 'STATEMENTS';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
    });
  });

  describe('facet switching', () => {
    it('should change active facet and persist ListKind to preferences', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      act(() => {
        result.current.setActiveFacet('receipts' as DataFacet);
      });

      expect(mockSetLastTab).toHaveBeenCalledWith('/exports', 'RECEIPTS');
    });
  });

  describe('facet configuration', () => {
    it('should return correct facet configuration', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      const config = result.current.getFacetConfig('receipts' as DataFacet);
      expect(config).toEqual({
        id: 'receipts',
        label: 'Receipts',
        listKind: 'RECEIPTS',
      });
    });

    it('should provide available facets', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.availableFacets).toHaveLength(3);
      expect(result.current.availableFacets).toEqual(sampleFacets);
    });

    it('should check if facet is active correctly', () => {
      mockLastTab['/exports'] = 'STATEMENTS';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.isFacetActive('statements' as DataFacet)).toBe(
        true,
      );
      expect(result.current.isFacetActive('receipts' as DataFacet)).toBe(false);
    });
  });

  describe('backward compatibility', () => {
    it('should provide correct ListKind for current facet', () => {
      mockLastTab['/exports'] = 'STATEMENTS';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.getCurrentListKind()).toBe('STATEMENTS');
    });

    it('should fallback to facet ID when no ListKind mapping exists', () => {
      const facetsWithoutListKind: DataFacetConfig[] = [
        {
          id: 'custom-facet' as DataFacet,
          label: 'Custom Facet',
          // No listKind property
        },
      ];

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: facetsWithoutListKind,
        }),
      );

      expect(result.current.getCurrentListKind()).toBe('custom-facet');
    });
  });

  describe('edge cases', () => {
    it('should handle empty facets array', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: [],
        }),
      );

      expect(result.current.activeFacet).toBe('transactions'); // fallback
      expect(result.current.availableFacets).toEqual([]);
    });

    it('should handle invalid saved preference', () => {
      mockLastTab['/exports'] = 'INVALID_LISTKIND';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('transactions'); // fallback to default
    });
  });
});
