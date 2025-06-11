// ABIS_ROUTE
import { GetAbisPage, Reload } from '@app';
import { abis, sdk, types } from '@models';

export const getAbisPage = (
  kind: types.ListKind,
  first: number,
  pageSize: number,
  sort: sdk.SortSpec,
  filter: string,
): Promise<abis.AbisPage> => GetAbisPage(kind, first, pageSize, sort, filter);

export const reload = (): Promise<void> => Reload();

// ABIS_ROUTE
