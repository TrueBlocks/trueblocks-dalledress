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

describe('Names View + useActiveFacet Integration Tests', () => {
  let mockLastTab: Record<string, string>;
  let mockSetLastTab: ReturnType<typeof vi.fn>;

  const namesFacets: DataFacetConfig[] = [
    {
      id: 'all' as DataFacet,
      label: 'All',
      listKind: 'ALL' as any,
      isDefault: true,
    },
    {
      id: 'custom' as DataFacet,
      label: 'Custom',
      listKind: 'CUSTOM' as any,
    },
    {
      id: 'prefund' as DataFacet,
      label: 'Prefund',
      listKind: 'PREFUND' as any,
    },
    {
      id: 'regular' as DataFacet,
      label: 'Regular',
      listKind: 'REGULAR' as any,
    },
    {
      id: 'baddress' as DataFacet,
      label: 'Bad Addresses',
      listKind: 'BADDRESS' as any,
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
      effectiveChain: 'mainnet',
      effectiveAddress: '0x123',
      projects: [],
      currentProject: { name: 'test-project', chain: 'mainnet' },
    } as any);
  });

  describe('useActiveFacet hook behavior with Names facets', () => {
    it('returns correct default values for Names view', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      expect(result.current.activeFacet).toBe('all');
      expect(result.current.getCurrentListKind()).toBe('ALL');
      expect(result.current.availableFacets).toHaveLength(5);
    });

    it('correctly maps facet IDs to ListKind values', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'custom' as DataFacet,
          viewRoute: '/names',
        }),
      );

      expect(result.current.getCurrentListKind()).toBe('CUSTOM');
    });

    it('maintains facet consistency across state changes', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      // Initial state
      expect(result.current.activeFacet).toBe('all');

      // Change to prefund facet
      act(() => {
        result.current.setActiveFacet('prefund' as DataFacet);
      });

      // Verify the preference was set correctly
      expect(mockSetLastTab).toHaveBeenCalledWith('/names', 'PREFUND');

      // Update mock to simulate preference persistence
      mockLastTab['/names'] = 'PREFUND';

      // Re-render hook with updated preferences
      const { result: result2 } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      expect(result2.current.activeFacet).toBe('prefund');
      expect(result2.current.getCurrentListKind()).toBe('PREFUND');
    });

    it('respects saved preferences for Names view', () => {
      mockLastTab['/names'] = 'custom';

      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      expect(result.current.activeFacet).toBe('custom');
      expect(result.current.getCurrentListKind()).toBe('CUSTOM');
    });

    it('persists facet changes to preferences', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      act(() => {
        result.current.setActiveFacet('baddress' as DataFacet);
      });

      expect(mockSetLastTab).toHaveBeenCalledWith('/names', 'BADDRESS');
    });

    it('handles all Names view facets correctly', () => {
      // Test each facet individually with fresh hook instances
      const facetMappings = [
        { facet: 'all', listKind: 'ALL' },
        { facet: 'custom', listKind: 'CUSTOM' },
        { facet: 'prefund', listKind: 'PREFUND' },
        { facet: 'regular', listKind: 'REGULAR' },
        { facet: 'baddress', listKind: 'BADDRESS' },
      ];

      facetMappings.forEach(({ facet, listKind }) => {
        // Set up mock with this facet as saved preference
        mockLastTab['/names'] = listKind;

        const { result } = renderHook(() =>
          useActiveFacet({
            facets: namesFacets,
            defaultFacet: 'all' as DataFacet,
            viewRoute: '/names',
          }),
        );

        expect(result.current.activeFacet).toBe(facet);
        expect(result.current.getCurrentListKind()).toBe(listKind);
      });
    });
  });

  describe('error handling and edge cases', () => {
    it('handles invalid facet selection gracefully', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      act(() => {
        result.current.setActiveFacet('invalid-facet' as DataFacet);
      });

      // Should remain on default facet
      expect(result.current.activeFacet).toBe('all');
    });

    it('handles missing default facet gracefully', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets.filter((f) => !f.isDefault),
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      // Should fall back to first available facet
      expect(result.current.availableFacets.length).toBeGreaterThan(0);
      expect(result.current.activeFacet).toBeTruthy();
    });
  });

  describe('integration with Names view constants and patterns', () => {
    it('provides all expected Names facets', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      const facetIds = result.current.availableFacets.map((f) => f.id);
      expect(facetIds).toEqual([
        'all',
        'custom',
        'prefund',
        'regular',
        'baddress',
      ]);

      const facetLabels = result.current.availableFacets.map((f) => f.label);
      expect(facetLabels).toEqual([
        'All',
        'Custom',
        'Prefund',
        'Regular',
        'Bad Addresses',
      ]);
    });

    it('matches expected Names view route pattern', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          defaultFacet: 'all' as DataFacet,
          viewRoute: '/names',
        }),
      );

      // Test that route is correctly handled internally
      act(() => {
        result.current.setActiveFacet('custom' as DataFacet);
      });

      expect(mockSetLastTab).toHaveBeenCalledWith('/names', 'CUSTOM');
    });
  });
});
