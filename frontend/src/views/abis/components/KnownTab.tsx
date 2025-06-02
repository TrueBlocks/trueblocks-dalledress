import React from 'react';

import { BaseTab, FormField } from '@components';
import { types } from '@models';

import { KNOWN_ABI_COLUMNS } from '../columnDefinitions';

interface KnownAbisTabProps {
  data: types.Abi[];
  loading: boolean;
  error: Error | null;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onSelect?: (address: string) => void;
  tableKey: { viewName: string; tabName: string };
}

export const KnownTab = ({
  data,
  loading,
  error,
  onSubmit,
  onSelect,
  tableKey,
}: KnownAbisTabProps) => {
  const handleAction = (item: Record<string, unknown>) => {
    const abi = item as unknown as types.Abi;

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
      onSubmit={onSubmit}
      onAction={handleAction}
      tableKey={tableKey}
    />
  );
};
