import { useEffect, useState } from 'react';

import { GetAbisPage } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useIcons } from '@hooks';
import { Badge, Button, Group, Stack, Text } from '@mantine/core';
import { types } from '@models';
import { Log } from '@utils';

interface AbisSummary {
  totalAbis: number;
  totalFunctions: number;
  recentlyAdded: number;
}

interface AbisPanelProps {
  onViewAll?: () => void;
  onAddAbi?: () => void;
}

export const AbisPanel = ({ onViewAll, onAddAbi }: AbisPanelProps) => {
  const [summary, setSummary] = useState<AbisSummary>({
    totalAbis: 0,
    totalFunctions: 0,
    recentlyAdded: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { ABIs } = useIcons();

  useEffect(() => {
    const fetchAbisSummary = async () => {
      try {
        setLoading(true);
        const abisData = await GetAbisPage(
          types.ListKind.DOWNLOADED,
          0,
          1000,
          { fields: ['address'], orders: [true] },
          '',
        );

        if (abisData?.abis) {
          setSummary({
            totalAbis: abisData.abis.length,
            totalFunctions: 0,
            recentlyAdded: 0,
          });
        }

        const functionsData = await GetAbisPage(
          types.ListKind.FUNCTIONS,
          0,
          1000,
          { fields: ['name'], orders: [true] },
          '',
        );

        if (functionsData?.functions) {
          setSummary((prev) => ({
            ...prev,
            totalFunctions: functionsData.functions?.length || 0,
          }));
        }

        setError(null);
      } catch (err) {
        Log(`Error fetching ABIs summary: ${err}`);
        setError('Failed to load ABIs');
      } finally {
        setLoading(false);
      }
    };

    fetchAbisSummary();
  }, []);

  return (
    <DashboardCard
      title="ABIs"
      subtitle={`${summary.totalAbis} contracts`}
      icon={<ABIs size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator
            status={summary.totalAbis > 0 ? 'healthy' : 'inactive'}
            label="ABI Database"
            count={summary.totalAbis}
            size="xs"
          />
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
          <Badge size="sm" variant="light" color="blue">
            Functions: {summary.totalFunctions}
          </Badge>
          {summary.recentlyAdded > 0 && (
            <Badge size="sm" variant="light" color="green">
              Recent: {summary.recentlyAdded}
            </Badge>
          )}
        </div>

        <Text size="xs" c="dimmed">
          Smart contract interfaces and function signatures
        </Text>

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="light" onClick={onAddAbi}>
            Add ABI
          </Button>
          <Button size="xs" variant="outline" onClick={onViewAll}>
            View All
          </Button>
        </Group>
      </Stack>
    </DashboardCard>
  );
};
