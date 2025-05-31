import React from 'react';

import { FormField, Table } from '@components';

interface BaseTabProps<T extends Record<string, unknown>> {
  data: T[];
  columns: FormField<T>[];
  tableKey: { viewName: string; tabName: string };
  loading: boolean;
  error: Error | null;
  onAction?: (item: T) => void;
  pagination?: {
    currentPage: number;
    pageSize: number;
    totalItems: number;
  };
  emptyMessage?: string;
  loadingMessage?: string;
}

export function BaseTab<T extends Record<string, unknown>>({
  data,
  columns,
  loading,
  error,
  onAction: _onAction,
  tableKey,
  pagination: _pagination,
  emptyMessage = 'No data available',
  loadingMessage = 'Loading...',
}: BaseTabProps<T>) {
  // Handle loading state
  if (loading) {
    return (
      <div className="loading-container">
        <p>{loadingMessage}</p>
      </div>
    );
  }

  // Handle error state
  if (error) {
    return (
      <div className="error-message-placeholder">
        <h3>Error loading data</h3>
        <p>{error.message}</p>
      </div>
    );
  }

  // Handle empty state
  if (!data || data.length === 0) {
    return (
      <div className="empty-state">
        <p>{emptyMessage}</p>
      </div>
    );
  }

  // Render the table
  return (
    <Table
      data={data as Record<string, unknown>[]}
      columns={columns as FormField<Record<string, unknown>>[]}
      tableKey={tableKey}
      loading={loading}
      onSubmit={() => {}}
    />
  );
}
