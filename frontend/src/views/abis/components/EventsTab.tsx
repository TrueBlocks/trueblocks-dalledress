import React from 'react';

import { BaseTab, FormField } from '@components';
import { types } from '@models';

import { FUNCTION_COLUMNS } from '../columnDefinitions';

interface EventsTabProps {
  data: types.Function[];
  loading: boolean;
  error: Error | null;
  onSelect?: (encoding: string) => void;
}

export const EventsTab = ({
  data,
  loading,
  error,
  onSelect,
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
      onAction={handleAction}
      tableKey={{ viewName: 'abis', tabName: 'events' }}
      emptyMessage="No events found"
      loadingMessage="Loading events..."
    />
  );
};
