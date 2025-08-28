import { ReactElement } from 'react';

import type { DataFacet } from '@hooks';
import { dalledress, project, types } from '@models';

import { Gallery } from './gallery';
import { Generator } from './generator';

export function renderers(
  pageData: dalledress.DalleDressPage | null,
  viewStateKey?: project.ViewStateKey,
  setActiveFacet?: (f: DataFacet) => void,
) {
  return {
    [types.DataFacet.GENERATOR]: () => <Generator pageData={pageData} />,
    [types.DataFacet.GALLERY]: () => (
      <Gallery
        pageData={pageData}
        viewStateKey={viewStateKey}
        setActiveFacet={setActiveFacet}
      />
    ),
  } as Record<types.DataFacet, () => ReactElement>;
}
