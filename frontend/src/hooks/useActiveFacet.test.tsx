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
      listKind: 'transactions' as any, // TODO: Fix after backend alignment
      isDefault: true,
    },
    {
      id: 'receipts' as DataFacet,
      label: 'Receipts',
      listKind: 'receipts' as any, // TODO: Fix after backend alignment
    },
    {
      id: 'statements' as DataFacet,
      label: 'Statements',
      listKind: 'statements' as any, // TODO: Fix after backend alignment
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
      effectiveAddress: '0x123',
      effectiveChain: 'mainnet',
      hasActiveProject: true,
      canExport: true,
    } as any);
  });

  describe('facet selection and defaults', () => {
    it('should use default facet when no preference exists', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('transactions');
      expect(result.current.isFacetActive('transactions')).toBe(true);
      expect(result.current.isFacetActive('receipts')).toBe(false);
    });

    it('should use saved preference when available', () => {
      mockLastTab['/exports'] = 'receipts';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('receipts');
      expect(result.current.isFacetActive('receipts')).toBe(true);
    });

    it('should fallback to default if saved preference is invalid', () => {
      mockLastTab['/exports'] = 'invalid-facet';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.activeFacet).toBe('transactions');
    });

    it('should use explicit defaultFacet parameter', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
          defaultFacet: 'statements' as DataFacet,
        }),
      );

      expect(result.current.activeFacet).toBe('statements');
    });
  });

  describe('facet switching', () => {
    it('should change active facet and persist to preferences', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      act(() => {
        result.current.setActiveFacet('receipts' as DataFacet);
      });

      expect(mockSetLastTab).toHaveBeenCalledWith('/exports', 'receipts');
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
        listKind: 'receipts',
      });
    });

    it('should provide available facets', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.availableFacets).toEqual(sampleFacets);
    });

    it('should provide default facet', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.getDefaultFacet()).toBe('transactions');
    });
  });

  describe('backward compatibility', () => {
    it('should provide ListKind for current facet', () => {
      mockLastTab['/exports'] = 'receipts';

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: sampleFacets,
        }),
      );

      expect(result.current.getCurrentListKind()).toBe('receipts');
    });

    it('should fallback to facet id when no listKind configured', () => {
      const facetsWithoutListKind = [
        {
          id: 'custom' as DataFacet,
          label: 'Custom',
          isDefault: true,
        },
      ];

      const { result } = renderHook(() =>
        useActiveFacet({
          viewRoute: '/exports',
          facets: facetsWithoutListKind,
        }),
      );

      expect(result.current.getCurrentListKind()).toBe('custom');
    });
  });
});
