import { types } from '@models';

/**
 * Re-export the generated DataFacet enum from the backend
 * This replaces direct DataFacet usage in frontend components
 */
export type DataFacet = types.DataFacet;

/**
 * Configuration for a data facet
 */
export interface DataFacetConfig {
  /** The facet identifier */
  id: DataFacet;
  /** Display name for the facet */
  label: string;
  /** Whether this facet is the default for its view */
  isDefault?: boolean;
}

/**
 * Return type for useActiveFacet hook
 */
export interface UseActiveFacetReturn {
  /** Currently active facet for the view */
  activeFacet: DataFacet;

  /** Function to change the active facet */
  setActiveFacet: (facet: DataFacet) => void;

  /** Available facets for the current view */
  availableFacets: DataFacetConfig[];

  /** Get configuration for a specific facet */
  getFacetConfig: (facet: DataFacet) => DataFacetConfig | undefined;

  /** Check if a facet is active */
  isFacetActive: (facet: DataFacet) => boolean;

  /** Get the default facet for the view */
  getDefaultFacet: () => DataFacet;

  /** Backend API compatibility: get types.DataFacet for current facet */
  getCurrentDataFacet: () => types.DataFacet | string;
}

/**
 * Hook parameters
 */
export interface UseActiveFacetParams {
  /** The view route (e.g., '/exports', '/names') */
  viewRoute: string;

  /** Available facets for this view */
  facets: DataFacetConfig[];

  /** Optional: override default facet selection */
  defaultFacet?: DataFacet;
}
