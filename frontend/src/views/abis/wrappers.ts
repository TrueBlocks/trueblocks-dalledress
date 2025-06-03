import { GetAbisPage, Reload, RemoveAbi } from '@app';
import { abis, sorting, types } from '@models';

// Wrappers for the Wails-generated functions to provide type safety and convenience

export const getAbisPage = (
  kind: types.ListKind,
  first: number,
  pageSize: number,
  sort: sorting.SortDef,
  filter: string,
): Promise<abis.AbisPage> => {
  return GetAbisPage(kind, first, pageSize, sort, filter);
};

export const removeAbi = (address: string): Promise<void> => {
  return RemoveAbi(address);
};

export const reload = (): Promise<void> => {
  return Reload();
};
