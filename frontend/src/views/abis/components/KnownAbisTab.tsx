import React from 'react';

import { BaseTab, FormField } from '@components';
import { sorting, types } from '@models';

import { KNOWN_ABI_COLUMNS } from '../columnDefinitions';

interface KnownAbisTabProps {
  data: types.Abi[];
  loading: boolean;
  error: Error | null;
  sort?: sorting.SortDef | null;
  onSortChange?: (sort: sorting.SortDef | null) => void;
  filter?: string;
  onFilterChange?: (filter: string) => void;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onSelect?: (address: string) => void;
  tableKey: { viewName: string; tabName: string };
}

export const KnownAbisTab = ({
  data,
  loading,
  error,
  sort,
  onSortChange,
  filter,
  onFilterChange,
  onSubmit,
  onSelect,
  tableKey,
}: KnownAbisTabProps) => {
  const handleAction = (item: Record<string, unknown>) => {
    const abi = item as unknown as types.Abi;
    console.log('Action on Known ABI:', abi.address);

    if (onSelect) {
      onSelect(abi.address.toString());
    }
  };

  return (
    <BaseTab
      data={data as unknown as Record<string, unknown>[]}
      columns={
        KNOWN_ABI_COLUMNS as unknown as FormField<Record<string, unknown>>[]
      }
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
