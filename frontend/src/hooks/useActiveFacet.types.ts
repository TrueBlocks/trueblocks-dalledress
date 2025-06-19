import { types } from '@models';

/**
 * Data facet types that will replace ListKind string usage
 * This represents the new architecture for tab/facet selection
 */
export type DataFacet =
  // Exports view facets
  | 'statements'
  | 'transfers'
  | 'balances'
  | 'transactions'
  // Chunks view facets
  | 'stats'
  | 'index'
  | 'blooms'
  | 'manifest'
  | 'chunk-summary'
  // Names view facets
  | 'entity-names'
  // Monitors view facets
  | 'txs'
  // ABIs view facets
  | 'get-abis'
  // Future extensibility
  | string;

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
  /** ListKind mapping for backward compatibility */
  listKind?: types.ListKind;
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

  /** Backward compatibility: get ListKind for current facet */
  getCurrentListKind: () => types.ListKind | string;
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
