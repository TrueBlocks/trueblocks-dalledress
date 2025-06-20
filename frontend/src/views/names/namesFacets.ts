import { types } from '@models';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

/**
 * Data facet configuration for the Names view
 * This replaces the direct ListKind usage and provides
 * a mapping between DataFacet and ListKind for backward compatibility
 */
export const namesFacets: DataFacetConfig[] = [
  {
    id: 'all',
    label: 'All',
    listKind: types.ListKind.ALL,
    isDefault: true,
  },
  {
    id: 'custom',
    label: 'Custom',
    listKind: types.ListKind.CUSTOM,
  },
  {
    id: 'prefund',
    label: 'Prefund',
    listKind: types.ListKind.PREFUND,
  },
  {
    id: 'regular',
    label: 'Regular',
    listKind: types.ListKind.REGULAR,
  },
  {
    id: 'baddress',
    label: 'Baddress',
    listKind: types.ListKind.BADDRESS,
  },
];

/**
 * The default facet for the Names view
 */
export const NAMES_DEFAULT_FACET = 'all' as const;

/**
 * Route identifier for the Names view
 */
export const NAMES_ROUTE = '/names' as const;
