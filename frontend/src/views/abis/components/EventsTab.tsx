import React from 'react';

import { BaseTab, FormField } from '@components';
import { types } from '@models';

import { FUNCTION_COLUMNS } from '../columnDefinitions';

interface EventsTabProps {
  data: types.Function[];
  loading: boolean;
  error: Error | null;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onSelect?: (encoding: string) => void;
  tableKey: { viewName: string; tabName: string };
}

export const EventsTab = ({
  data,
  loading,
  error,
  onSubmit,
  onSelect,
  tableKey,
}: EventsTabProps) => {
  const handleAction = (item: Record<string, unknown>) => {
    const event = item as unknown as types.Function;

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
      onSubmit={onSubmit}
      onAction={handleAction}
      tableKey={tableKey}
    />
  );
};
