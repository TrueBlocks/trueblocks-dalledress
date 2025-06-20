import { useCallback, useMemo } from 'react';

import { types } from '@models';

import {
  DataFacet,
  DataFacetConfig,
  UseActiveFacetParams,
  UseActiveFacetReturn,
} from './useActiveFacet.types';
import { useActiveProject } from './useActiveProject';

/**
 * Hook for managing active data facets in views
 *
 * This hook provides a clean API for:
 * - Managing the currently active facet in a view
 * - Persisting facet selection to user preferences
 * - Providing facet configuration and metadata
 * - Backward compatibility with existing ListKind usage
 *
 * @param params Configuration for the hook
 * @returns Hook interface for facet management
 */
export const useActiveFacet = (
  params: UseActiveFacetParams,
): UseActiveFacetReturn => {
  const { viewRoute, facets, defaultFacet } = params;
  const { lastTab, setLastTab } = useActiveProject();

  // Determine the default facet for this view
  const computedDefaultFacet = useMemo((): DataFacet => {
    if (defaultFacet) return defaultFacet;

    const defaultConfig = facets.find((f) => f.isDefault);
    if (defaultConfig) return defaultConfig.id;

    // Fallback to first facet
    return facets[0]?.id || 'transactions';
  }, [defaultFacet, facets]);

  // Get the currently active facet from preferences or default
  const activeFacet = useMemo((): DataFacet => {
    const saved = lastTab[viewRoute];

    if (saved) {
      // First try to find facet by ID
      if (facets.some((f) => f.id === saved)) {
        return saved as DataFacet;
      }
      // Then try to find by ListKind
      const facetByListKind = facets.find((f) => f.listKind === saved);
      if (facetByListKind) {
        return facetByListKind.id;
      }
    }
    return computedDefaultFacet;
  }, [lastTab, viewRoute, facets, computedDefaultFacet]);

  // Function to change the active facet
  const setActiveFacet = useCallback(
    (facet: DataFacet): void => {
      // Find the corresponding ListKind for this facet
      const config = facets.find((f) => f.id === facet);
      const listKindValue = config?.listKind || (facet as types.ListKind);
      setLastTab(viewRoute, listKindValue);
    },
    [setLastTab, viewRoute, facets],
  );

  // Get configuration for a specific facet
  const getFacetConfig = useCallback(
    (facet: DataFacet): DataFacetConfig | undefined => {
      return facets.find((f) => f.id === facet);
    },
    [facets],
  );

  // Check if a facet is active
  const isFacetActive = useCallback(
    (facet: DataFacet): boolean => {
      return activeFacet === facet;
    },
    [activeFacet],
  );

  // Get the default facet for the view
  const getDefaultFacet = useCallback((): DataFacet => {
    return computedDefaultFacet;
  }, [computedDefaultFacet]);

  // Backward compatibility: get ListKind for current facet
  const getCurrentListKind = useCallback(() => {
    const config = getFacetConfig(activeFacet);
    return config?.listKind || activeFacet;
  }, [activeFacet, getFacetConfig]);

  return {
    activeFacet,
    setActiveFacet,
    availableFacets: facets,
    getFacetConfig,
    isFacetActive,
    getDefaultFacet,
    getCurrentListKind,
  };
};
