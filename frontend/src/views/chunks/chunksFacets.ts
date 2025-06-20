import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

import { DataFacetConfig } from '../../hooks/useActiveFacet.types';

export const chunksFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.STATS,
    label: toProperCase(types.DataFacet.STATS),
    isDefault: true,
  },
  {
    id: types.DataFacet.INDEX,
    label: toProperCase(types.DataFacet.INDEX),
  },
  {
    id: types.DataFacet.BLOOMS,
    label: toProperCase(types.DataFacet.BLOOMS),
  },
  {
    id: types.DataFacet.MANIFEST,
    label: toProperCase(types.DataFacet.MANIFEST),
  },
];

export const CHUNKS_DEFAULT_FACET = types.DataFacet.STATS;
export const CHUNKS_ROUTE = '/chunks' as const;
