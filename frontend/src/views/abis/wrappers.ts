import { GetAbisPage, Reload, RemoveAbi } from '@app';
import { abis, sdk, types } from '@models';

// Type-safe wrappers for Wails-generated backend functions
export const getAbisPage = (
  kind: types.ListKind,
  first: number,
  pageSize: number,
  sort: sdk.SortSpec,
  filter: string,
): Promise<abis.AbisPage> => GetAbisPage(kind, first, pageSize, sort, filter);

export const removeAbi = (address: string): Promise<void> => RemoveAbi(address);

export const reload = (): Promise<void> => Reload();
