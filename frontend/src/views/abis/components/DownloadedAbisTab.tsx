import React from 'react';

import { BaseTab, FormField } from '@components';
import { sorting, types } from '@models';

import { ABI_COLUMNS } from '../columnDefinitions';

interface DownloadedAbisTabProps {
  data: types.Abi[];
  loading: boolean;
  error: Error | null;
  sort?: sorting.SortDef | null;
  onSortChange?: (sort: sorting.SortDef | null) => void;
  filter?: string;
  onFilterChange?: (filter: string) => void;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onDelete?: (address: string) => void;
  onHistory?: (address: string) => void;
  tableKey: { viewName: string; tabName: string };
}

export const DownloadedAbisTab = ({
  data,
  loading,
  error,
  sort,
  onSortChange,
  filter,
  onFilterChange,
  onSubmit,
  onDelete,
  onHistory,
  tableKey,
}: DownloadedAbisTabProps) => {
  const handleAction = (item: Record<string, unknown>) => {
    const abi = item as unknown as types.Abi;
    console.log('Action on ABI:', abi.address);

    if (onDelete) {
      onDelete(abi.address.toString());
    }
    if (onHistory) {
      onHistory(abi.address.toString());
    }
  };

  return (
    <BaseTab
      data={data as unknown as Record<string, unknown>[]}
      columns={ABI_COLUMNS as unknown as FormField<Record<string, unknown>>[]}
      loading={loading}
      error={error}
      sort={sort}
      onSortChange={onSortChange}
      filter={filter}
      onFilterChange={onFilterChange}
      onSubmit={onSubmit}
      onAction={handleAction}
      tableKey={tableKey}
    />
  );
};
