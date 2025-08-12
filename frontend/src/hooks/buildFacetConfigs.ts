import { types } from '@models';

import { DataFacetConfig } from './useActiveFacet';

export function buildFacetConfigs(
  viewConfig: types.ViewConfig,
  labelTransform?: (id: string) => string,
): DataFacetConfig[] {
  return (viewConfig.facetOrder || [])
    .filter((facetId: string) => Boolean(viewConfig.facets[facetId]))
    .map((facetId: string) => {
      const facetCfg = viewConfig.facets[facetId];
      return {
        id: facetId as types.DataFacet,
        label:
          facetCfg?.name ||
          (labelTransform ? labelTransform(facetId) : facetId),
        dividerBefore: facetCfg?.dividerBefore,
      };
    });
}
