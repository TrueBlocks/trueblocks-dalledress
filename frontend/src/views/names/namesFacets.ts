import { types } from '@models';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

/**
 * Data facet configuration for the Names view
 * This replaces the direct ListKind usage and provides
 * a mapping between DataFacet and ListKind for backward compatibility
 */
export const namesFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.ALL,
    label: 'All',
    listKind: types.ListKind.ALL,
    isDefault: true,
  },
  {
    id: types.DataFacet.CUSTOM,
    label: 'Custom',
    listKind: types.ListKind.CUSTOM,
  },
  {
    id: types.DataFacet.PREFUND,
    label: 'Prefund',
    listKind: types.ListKind.PREFUND,
  },
  {
    id: types.DataFacet.REGULAR,
    label: 'Regular',
    listKind: types.ListKind.REGULAR,
  },
  {
    id: types.DataFacet.BADDRESS,
    label: 'Baddress',
    listKind: types.ListKind.BADDRESS,
  },
];

/**
 * The default facet for the Names view
 */
export const NAMES_DEFAULT_FACET = types.DataFacet.ALL;

/**
 * Route identifier for the Names view
 */
export const NAMES_ROUTE = '/names' as const;
