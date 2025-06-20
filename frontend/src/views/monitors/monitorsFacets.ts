import { types } from '@models';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

/**
 * Data facet configuration for the Monitors view
 * This replaces the direct ListKind usage and provides
 * a mapping between DataFacet and ListKind for backward compatibility
 */
export const monitorsFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.MONITORS,
    label: 'Monitors',
    listKind: types.ListKind.MONITORS,
    isDefault: true,
  },
];

/**
 * The default facet for the Monitors view
 */
export const MONITORS_DEFAULT_FACET = types.DataFacet.MONITORS;

/**
 * Route identifier for the Monitors view
 */
export const MONITORS_ROUTE = '/monitors' as const;
