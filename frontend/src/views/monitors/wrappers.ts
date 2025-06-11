// MONITORS_ROUTE
import { GetMonitorsPage, MonitorsClean, Reload } from '@app';
import { monitors, sdk, types } from '@models';

export const getMonitorsPage = (
  kind: types.ListKind,
  first: number,
  pageSize: number,
  sort: sdk.SortSpec,
  filter: string,
): Promise<monitors.MonitorsPage> =>
  GetMonitorsPage(kind, first, pageSize, sort, filter);

export const cleanMonitors = (addresses: string[]): Promise<void> =>
  MonitorsClean(addresses);

export const reload = (): Promise<void> => Reload();

// MONITORS_ROUTE
