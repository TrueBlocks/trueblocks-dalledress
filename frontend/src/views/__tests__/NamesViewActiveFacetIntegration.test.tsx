import { DataFacet, DataFacetConfig, useActiveFacet } from '@hooks';
import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Mock the useActiveProject hook
vi.mock('../../hooks/useActiveProject', () => ({
  useActiveProject: vi.fn(),
}));

const { useActiveProject } = await import('../../hooks/useActiveProject');
const mockedUseActiveProject = vi.mocked(useActiveProject);

describe('Names View + useActiveFacet Integration Tests', () => {
  let mockLastFacetMap: Record<string, string>;
  let mockSetLastFacet: ReturnType<typeof vi.fn>;
  const namesFacets: DataFacetConfig[] = [
    {
      id: 'all' as DataFacet,
      label: 'All',
    },
    {
      id: 'custom' as DataFacet,
      label: 'Custom',
    },
    {
      id: 'prefund' as DataFacet,
      label: 'Prefund',
    },
    {
      id: 'regular' as DataFacet,
      label: 'Regular',
    },
    {
      id: 'baddress' as DataFacet,
      label: 'Bad Addresses',
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();

    // Mock successful project state
    mockLastFacetMap = {};
    mockSetLastFacet = vi.fn();

    mockedUseActiveProject.mockReturnValue({
      lastFacetMap: mockLastFacetMap,
      setLastFacet: mockSetLastFacet,
      getLastFacet: vi.fn((view: string) => mockLastFacetMap[view] || ''),
      // Mock other required properties
      activeChain: 'mainnet',
      activeAddress: '0x123',
      activeContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
      projects: [],
      currentProject: { name: 'test-project', chain: 'mainnet' },
    } as any);
  });

  describe('useActiveFacet hook behavior with Names facets', () => {
    it('returns correct default values for Names view', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          viewRoute: 'names',
        }),
      );

      expect(result.current.activeFacet).toBe('all');
      expect(result.current.getCurrentDataFacet()).toBe('all');
      expect(result.current.availableFacets).toHaveLength(5);
    });

    it('maintains facet consistency across state changes', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          viewRoute: 'names',
        }),
      );

      // Initial state (should be first facet since no saved preference)
      expect(result.current.activeFacet).toBe('all');

      // Change to prefund facet
      act(() => {
        result.current.setActiveFacet('prefund' as DataFacet);
      });

      // Verify the preference was set correctly
      expect(mockSetLastFacet).toHaveBeenCalledWith('names', 'prefund');

      // Update mock to simulate preference persistence
      mockLastFacetMap['names'] = 'prefund';

      // Re-render hook with updated preferences
      const { result: result2 } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          viewRoute: 'names',
        }),
      );

      expect(result2.current.activeFacet).toBe('prefund');
      expect(result2.current.getCurrentDataFacet()).toBe('prefund');
    });

    it('respects saved preferences for Names view', () => {
      mockLastFacetMap['names'] = 'custom';

      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          viewRoute: 'names',
        }),
      );

      expect(result.current.activeFacet).toBe('custom');
      expect(result.current.getCurrentDataFacet()).toBe('custom');
    });

    it('persists facet changes to preferences', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          viewRoute: 'names',
        }),
      );

      act(() => {
        result.current.setActiveFacet('baddress' as DataFacet);
      });

      expect(mockSetLastFacet).toHaveBeenCalledWith('names', 'baddress');
    });

    it('handles all Names view facets correctly', () => {
      // Test each facet individually with fresh hook instances
      const facetMappings = [
        { facet: 'all', dataFacet: 'all' },
        { facet: 'custom', dataFacet: 'custom' },
        { facet: 'prefund', dataFacet: 'prefund' },
        { facet: 'regular', dataFacet: 'regular' },
        { facet: 'baddress', dataFacet: 'baddress' },
      ];

      facetMappings.forEach(({ facet, dataFacet }) => {
        // Set up mock with this facet as saved preference
        mockLastFacetMap['names'] = dataFacet;

        const { result } = renderHook(() =>
          useActiveFacet({
            facets: namesFacets,
            viewRoute: 'names',
          }),
        );

        expect(result.current.activeFacet).toBe(facet);
        expect(result.current.getCurrentDataFacet()).toBe(dataFacet);
      });
    });
  });

  describe('error handling and edge cases', () => {
    it('handles invalid facet selection gracefully', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          viewRoute: 'names',
        }),
      );

      act(() => {
        result.current.setActiveFacet('invalid-facet' as DataFacet);
      });

      // Should remain on default facet
      expect(result.current.activeFacet).toBe('all');
    });
  });

  describe('integration with Names view constants and patterns', () => {
    it('provides all expected Names facets', () => {
      const { result } = renderHook(() =>
        useActiveFacet({
          facets: namesFacets,
          viewRoute: 'names',
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
          viewRoute: 'names',
        }),
      );

      // Test that route is correctly handled internally
      act(() => {
        result.current.setActiveFacet('custom' as DataFacet);
      });

      expect(mockSetLastFacet).toHaveBeenCalledWith('names', 'custom');
    });
  });
});
