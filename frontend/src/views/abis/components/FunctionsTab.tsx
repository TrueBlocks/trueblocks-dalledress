import React from 'react';

import { BaseTab, FormField } from '@components';
import { types } from '@models';

import { FUNCTION_COLUMNS } from '../columnDefinitions';

interface FunctionsTabProps {
  data: types.Function[];
  loading: boolean;
  error: Error | null;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onSelect?: (encoding: string) => void;
  tableKey: { viewName: string; tabName: string };
}

export const FunctionsTab = ({
  data,
  loading,
  error,
  onSubmit,
  onSelect,
  tableKey,
}: FunctionsTabProps) => {
  const handleAction = (item: Record<string, unknown>) => {
    const func = item as unknown as types.Function;

    if (onSelect) {
      onSelect(func.encoding.toString());
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
