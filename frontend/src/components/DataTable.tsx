import { createDataTable } from '@trueblocks/ui';
import type { PersistTableState } from '@trueblocks/ui';
import { GetTableState, SetTableState } from '../../wailsjs/go/app/App';

export type { Column, DataTableProps } from '@trueblocks/ui';

export const DataTable = createDataTable(
  (name) => GetTableState(name) as Promise<Partial<PersistTableState>>,
  (name, state) => SetTableState(name, state as Record<string, unknown>),
);
