import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const exportsFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATEMENTS,
    label: toProperCase(types.DataFacet.STATEMENTS),
  },
  {
    id: types.DataFacet.TRANSFERS,
    label: toProperCase(types.DataFacet.TRANSFERS),
  },
  {
    id: types.DataFacet.BALANCES,
    label: toProperCase(types.DataFacet.BALANCES),
  },
  {
    id: types.DataFacet.TRANSACTIONS,
    label: toProperCase(types.DataFacet.TRANSACTIONS),
    isDefault: true,
  },
];

export const EXPORTS_DEFAULT_FACET = types.DataFacet.TRANSACTIONS;
export const EXPORTS_ROUTE = '/exports' as const;
