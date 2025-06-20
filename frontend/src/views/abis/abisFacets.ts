import { types } from '@models';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

/**
 * Data facet configuration for the Abis view
 * This replaces the direct ListKind usage and provides
 * a mapping between DataFacet and ListKind for backward compatibility
 */
export const abisFacets: DataFacetConfig[] = [
  {
    id: 'downloaded',
    label: 'Downloaded',
    listKind: types.ListKind.DOWNLOADED,
    isDefault: true,
  },
  {
    id: 'known',
    label: 'Known',
    listKind: types.ListKind.KNOWN,
  },
  {
    id: 'functions',
    label: 'Functions',
    listKind: types.ListKind.FUNCTIONS,
  },
  {
    id: 'events',
    label: 'Events',
    listKind: types.ListKind.EVENTS,
  },
];

/**
 * The default facet for the Abis view
 */
export const ABIS_DEFAULT_FACET = 'downloaded' as const;

/**
 * Route identifier for the Abis view
 */
export const ABIS_ROUTE = '/abis' as const;
