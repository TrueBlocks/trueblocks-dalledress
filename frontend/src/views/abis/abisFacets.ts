import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

export const abisFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.DOWNLOADED,
    label: toProperCase(types.DataFacet.DOWNLOADED),
    isDefault: true,
  },
  {
    id: types.DataFacet.KNOWN,
    label: toProperCase(types.DataFacet.KNOWN),
  },
  {
    id: types.DataFacet.FUNCTIONS,
    label: toProperCase(types.DataFacet.FUNCTIONS),
  },
  {
    id: types.DataFacet.EVENTS,
    label: toProperCase(types.DataFacet.EVENTS),
  },
];

export const ABIS_DEFAULT_FACET = types.DataFacet.DOWNLOADED;
export const ABIS_ROUTE = '/abis' as const;
