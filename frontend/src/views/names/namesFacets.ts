import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

export const namesFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.ALL,
    label: toProperCase(types.DataFacet.ALL),
    isDefault: true,
  },
  {
    id: types.DataFacet.CUSTOM,
    label: toProperCase(types.DataFacet.CUSTOM),
  },
  {
    id: types.DataFacet.PREFUND,
    label: toProperCase(types.DataFacet.PREFUND),
  },
  {
    id: types.DataFacet.REGULAR,
    label: toProperCase(types.DataFacet.REGULAR),
  },
  {
    id: types.DataFacet.BADDRESS,
    label: toProperCase(types.DataFacet.BADDRESS),
  },
];

export const NAMES_DEFAULT_FACET = types.DataFacet.ALL;
export const NAMES_ROUTE = '/names' as const;
