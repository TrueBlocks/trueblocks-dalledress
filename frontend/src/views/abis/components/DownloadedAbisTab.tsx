import { BaseTab, FormField } from '@components';
import { types } from '@models';

import { ABI_COLUMNS } from '../columnDefinitions';

interface DownloadedAbisTabProps {
  data: types.Abi[];
  loading: boolean;
  error: Error | null;
  onDelete?: (address: string) => void;
  onHistory?: (address: string) => void;
}

export const DownloadedAbisTab = ({
  data,
  loading,
  error,
  onDelete,
  onHistory,
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
      onAction={handleAction}
      tableKey={{ viewName: 'abis', tabName: 'downloaded' }}
      emptyMessage="No downloaded ABIs found"
      loadingMessage="Loading downloaded ABIs..."
    />
  );
};
