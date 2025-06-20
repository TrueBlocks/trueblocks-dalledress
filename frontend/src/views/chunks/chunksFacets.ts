import { types } from '@models';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

/**
 * Data facet configuration for the Chunks view
 * This replaces the direct ListKind usage and provides
 * a mapping between DataFacet and ListKind for backward compatibility
 */
export const chunksFacets: DataFacetConfig[] = [
  {
    id: 'stats',
    label: 'Stats',
    listKind: types.ListKind.STATS,
    isDefault: true,
  },
  {
    id: 'index',
    label: 'Index',
    listKind: types.ListKind.INDEX,
  },
  {
    id: 'blooms',
    label: 'Blooms',
    listKind: types.ListKind.BLOOMS,
  },
  {
    id: 'manifest',
    label: 'Manifest',
    listKind: types.ListKind.MANIFEST,
  },
];

/**
 * The default facet for the Chunks view
 */
export const CHUNKS_DEFAULT_FACET = 'stats' as const;

/**
 * Route identifier for the Chunks view
 */
export const CHUNKS_ROUTE = '/chunks' as const;
