import { types } from '@models';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

/**
 * Data facet configuration for the Abis view
 * This replaces the direct ListKind usage and provides
 * a mapping between DataFacet and ListKind for backward compatibility
 */
export const abisFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.DOWNLOADED,
    label: 'Downloaded',
    listKind: types.ListKind.DOWNLOADED,
    isDefault: true,
  },
  {
    id: types.DataFacet.KNOWN,
    label: 'Known',
    listKind: types.ListKind.KNOWN,
  },
  {
    id: types.DataFacet.FUNCTIONS,
    label: 'Functions',
    listKind: types.ListKind.FUNCTIONS,
  },
  {
    id: types.DataFacet.EVENTS,
    label: 'Events',
    listKind: types.ListKind.EVENTS,
  },
];

/**
 * The default facet for the Abis view
 */
export const ABIS_DEFAULT_FACET = types.DataFacet.DOWNLOADED;

/**
 * Route identifier for the Abis view
 */
export const ABIS_ROUTE = '/abis' as const;
