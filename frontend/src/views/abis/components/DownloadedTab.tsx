import React from 'react';

import { BaseTab, FormField } from '@components';
import { types } from '@models';

import { ABI_COLUMNS } from '../columnDefinitions';

interface DownloadedAbisTabProps {
  data: types.Abi[];
  loading: boolean;
  error: Error | null;
  onSubmit?: (formData: Record<string, unknown>) => void;
  onDelete?: (address: string) => void;
  onHistory?: (address: string) => void;
  tableKey: { viewName: string; tabName: string };
}

export const DownloadedTab = ({
  data,
  loading,
  error,
  onSubmit,
  onDelete,
  onHistory,
  tableKey,
}: DownloadedAbisTabProps) => {
  const handleAction = (item: Record<string, unknown>) => {
    const abi = item as unknown as types.Abi;

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
      onSubmit={onSubmit}
      onAction={handleAction}
      tableKey={tableKey}
    />
  );
};
