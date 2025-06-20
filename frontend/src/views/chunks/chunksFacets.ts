import { types } from '@models';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

/**
 * Data facet configuration for the Chunks view
 * This replaces the direct ListKind usage and provides
 * a mapping between DataFacet and ListKind for backward compatibility
 */
export const chunksFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATS,
    label: 'Stats',
    listKind: types.ListKind.STATS,
    isDefault: true,
  },
  {
    id: types.DataFacet.INDEX,
    label: 'Index',
    listKind: types.ListKind.INDEX,
  },
  {
    id: types.DataFacet.BLOOMS,
    label: 'Blooms',
    listKind: types.ListKind.BLOOMS,
  },
  {
    id: types.DataFacet.MANIFEST,
    label: 'Manifest',
    listKind: types.ListKind.MANIFEST,
  },
];

/**
 * The default facet for the Chunks view
 */
export const CHUNKS_DEFAULT_FACET = types.DataFacet.STATS;

/**
 * Route identifier for the Chunks view
 */
export const CHUNKS_ROUTE = '/chunks' as const;
