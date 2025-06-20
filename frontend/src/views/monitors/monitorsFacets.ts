import { DataFacetConfig } from '@hooks';
import { types } from '@models';
import { toProperCase } from 'src/utils/toProper';

export const monitorsFacets: DataFacetConfig[] = [
  {
    id: types.DataFacet.MONITORS,
    label: toProperCase(types.DataFacet.MONITORS),
    isDefault: true,
  },
];

export const MONITORS_DEFAULT_FACET = types.DataFacet.MONITORS;
export const MONITORS_ROUTE = '/monitors' as const;
