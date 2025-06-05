// ADD_ROUTE
import { AbisCrudByAddr, GetAbisPage, Reload } from '@app';
import { abis, crud, sdk, types } from '@models';

export const getAbisPage = (
  kind: types.ListKind,
  first: number,
  pageSize: number,
  sort: sdk.SortSpec,
  filter: string,
): Promise<abis.AbisPage> => GetAbisPage(kind, first, pageSize, sort, filter);

export const abisCrud = (address: string): Promise<void> =>
  AbisCrudByAddr(crud.Operation.REMOVE, address);

export const reload = (): Promise<void> => Reload();

// ADD_ROUTE
