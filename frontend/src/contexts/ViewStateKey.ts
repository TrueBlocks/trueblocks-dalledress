import { project } from '@models';

export type ViewStateKey = project.ViewStateKey;

export const viewStateKeyToString = (key: ViewStateKey): string => {
  return `${key.viewName}/${key.facetName}/`;
};
