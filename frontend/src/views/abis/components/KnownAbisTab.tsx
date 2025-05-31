import React from 'react';

import { BaseTab, FormField } from '@components';
import { types } from '@models';

import { KNOWN_ABI_COLUMNS } from '../columnDefinitions';

interface KnownAbisTabProps {
  data: types.Abi[];
  loading: boolean;
  error: Error | null;
  onSelect?: (address: string) => void;
}

export const KnownAbisTab = ({
  data,
  loading,
  error,
  onSelect,
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
      onAction={handleAction}
      tableKey={{ viewName: 'abis', tabName: 'known' }}
      emptyMessage="No known ABIs found"
      loadingMessage="Loading known ABIs..."
    />
  );
};
