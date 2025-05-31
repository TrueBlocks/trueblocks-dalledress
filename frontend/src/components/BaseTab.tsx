import React from 'react';

import { FormField, Table, TableProvider } from '@components';
import { sorting } from '@models';

interface BaseTabProps<T extends Record<string, unknown>> {
  data: T[];
  columns: FormField<T>[];
  tableKey: { viewName: string; tabName: string };
  loading: boolean;
  error: Error | null;
  sort?: sorting.SortDef | null;
  onSortChange?: (sort: sorting.SortDef | null) => void;
  filter?: string;
  onFilterChange?: (filter: string) => void;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onAction?: (item: T) => void;
}

export function BaseTab<T extends Record<string, unknown>>({
  data,
  columns,
  loading,
  error: _error,
  sort,
  onSortChange,
  filter,
  onFilterChange,
  onSubmit,
  onAction: _onAction,
  tableKey,
}: BaseTabProps<T>) {
  // Always render table structure - let Table component handle all states
  return (
    <TableProvider>
      <div className="tableContainer">
        <Table
          data={data as Record<string, unknown>[]}
          columns={columns as FormField<Record<string, unknown>>[]}
          tableKey={tableKey}
          loading={loading}
          sort={sort}
          onSortChange={onSortChange}
          filter={filter}
          onFilterChange={onFilterChange}
          onSubmit={onSubmit || (() => {})}
        />
      </div>
    </TableProvider>
  );
}
