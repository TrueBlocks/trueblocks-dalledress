import { types } from '@models';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

/**
 * Data facet configuration for the Exports view
 * This replaces the direct ListKind usage and provides
 * a mapping between DataFacet and ListKind for backward compatibility
 */
export const exportsFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATEMENTS,
    label: 'Statements',
    listKind: types.ListKind.STATEMENTS,
  },
  {
    id: types.DataFacet.TRANSFERS,
    label: 'Transfers',
    listKind: types.ListKind.TRANSFERS,
  },
  {
    id: types.DataFacet.BALANCES,
    label: 'Balances',
    listKind: types.ListKind.BALANCES,
  },
  {
    id: types.DataFacet.TRANSACTIONS,
    label: 'Transactions',
    listKind: types.ListKind.TRANSACTIONS,
    isDefault: true,
  },
];

/**
 * The default facet for the Exports view
 */
export const EXPORTS_DEFAULT_FACET = types.DataFacet.TRANSACTIONS;

/**
 * Route identifier for the Exports view
 */
export const EXPORTS_ROUTE = '/exports' as const;
