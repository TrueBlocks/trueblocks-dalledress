import React from 'react';

import { BaseTab, FormField } from '@components';
import { sorting, types } from '@models';

import { FUNCTION_COLUMNS } from '../columnDefinitions';

interface EventsTabProps {
  data: types.Function[];
  loading: boolean;
  error: Error | null;
  sort?: sorting.SortDef | null;
  onSortChange?: (sort: sorting.SortDef | null) => void;
  filter?: string;
  onFilterChange?: (filter: string) => void;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onSelect?: (encoding: string) => void;
  tableKey: { viewName: string; tabName: string };
}

export const EventsTab = ({
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
}: EventsTabProps) => {
  const handleAction = (item: Record<string, unknown>) => {
    const event = item as unknown as types.Function;
    console.log('Action on Event:', event.encoding);

    if (onSelect) {
      onSelect(event.encoding.toString());
    }
  };

  return (
    <BaseTab
      data={data as unknown as Record<string, unknown>[]}
      columns={
        FUNCTION_COLUMNS as unknown as FormField<Record<string, unknown>>[]
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
