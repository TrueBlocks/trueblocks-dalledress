import React from 'react';

import { FormField, Table, TableProvider } from '@components';
import { ViewStateKey } from '@contexts';

import './BaseTab.css';

interface BaseTabProps<T extends Record<string, unknown>> {
  data: T[];
  columns: FormField<T>[];
  viewStateKey: ViewStateKey;
  loading: boolean;
  error: Error | null;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onAction?: (item: T) => void;
}

export function BaseTab<T extends Record<string, unknown>>({
  data,
  columns,
  loading,
  error: _error,
  onSubmit,
  onAction: _onAction,
  viewStateKey,
}: BaseTabProps<T>) {
  // Always render table structure - let Table component handle all states
  return (
    <TableProvider>
      <div className="tableContainer">
        <Table
          data={data as Record<string, unknown>[]}
          columns={columns as FormField<Record<string, unknown>>[]}
          viewStateKey={viewStateKey}
          loading={loading}
          onSubmit={onSubmit || (() => {})}
        />
      </div>
    </TableProvider>
  );
}
