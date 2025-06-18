import { useEffect, useState } from 'react';

import { DashboardCard, StatusIndicator } from '@components';
import { useIcons } from '@hooks';
import { Badge, Button, Group, Stack, Text } from '@mantine/core';
import { Log } from '@utils';

interface ExportsSummary {
  availableTypes: string[];
  recentExportCount: number;
  lastExportTime?: Date;
}

interface ExportsPanelProps {
  onViewAll?: () => void;
  onNewExport?: () => void;
}

export const ExportsPanel = ({ onViewAll, onNewExport }: ExportsPanelProps) => {
  const [summary, setSummary] = useState<ExportsSummary>({
    availableTypes: ['Transactions', 'Statements', 'Transfers', 'Balances'],
    recentExportCount: 0,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { Exports } = useIcons();

  useEffect(() => {
    const fetchExportsSummary = async () => {
      try {
        setLoading(true);
        setSummary({
          availableTypes: [
            'Transactions',
            'Statements',
            'Transfers',
            'Balances',
          ],
          recentExportCount: 0,
        });
        setError(null);
      } catch (err) {
        Log(`Error fetching exports summary: ${err}`);
        setError('Failed to load exports');
      } finally {
        setLoading(false);
      }
    };

    fetchExportsSummary();
  }, []);

  return (
    <DashboardCard
      title="Exports"
      subtitle={`${summary.availableTypes.length} types available`}
      icon={<Exports size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator status="healthy" label="Export System" size="xs" />
        </div>

        <div>
          <Text size="xs" c="dimmed" mb={4}>
            Available Types
          </Text>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
            {summary.availableTypes.map((type) => (
              <Badge key={type} size="xs" variant="outline">
                {type}
              </Badge>
            ))}
          </div>
        </div>

        {summary.recentExportCount > 0 && (
          <div>
            <Text size="xs" c="dimmed">
              Recent exports: {summary.recentExportCount}
            </Text>
          </div>
        )}

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="light" onClick={onNewExport}>
            New Export
          </Button>
          <Button size="xs" variant="outline" onClick={onViewAll}>
            View All
          </Button>
        </Group>
      </Stack>
    </DashboardCard>
  );
};
