import React from 'react';

import { BaseTab, FormField } from '@components';
import { types } from '@models';

import { FUNCTION_COLUMNS } from '../columnDefinitions';

interface FunctionsTabProps {
  data: types.Function[];
  loading: boolean;
  error: Error | null;
  onSelect?: (encoding: string) => void;
}

export const FunctionsTab = ({
  data,
  loading,
  error,
  onSelect,
}: FunctionsTabProps) => {
  const handleAction = (item: Record<string, unknown>) => {
    const func = item as unknown as types.Function;
    console.log('Action on Function:', func.encoding);

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
      onAction={handleAction}
      tableKey={{ viewName: 'abis', tabName: 'functions' }}
      emptyMessage="No functions found"
      loadingMessage="Loading functions..."
    />
  );
};
