import { useEffect, useState } from 'react';

import { GetChunksPage } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useIcons } from '@hooks';
import { Badge, Button, Group, Stack, Text } from '@mantine/core';
import { types } from '@models';
import { Log } from '@utils';

interface ChunksSummary {
  totalChunks: number;
  totalAddresses: number;
  totalBlocks: number;
  indexHealth: 'healthy' | 'warning' | 'error' | 'loading';
}

interface ChunksPanelProps {
  onViewAll?: () => void;
}

export const ChunksPanel = ({ onViewAll }: ChunksPanelProps) => {
  const [summary, setSummary] = useState<ChunksSummary>({
    totalChunks: 0,
    totalAddresses: 0,
    totalBlocks: 0,
    indexHealth: 'loading',
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { Chunks } = useIcons();

  useEffect(() => {
    const fetchChunksSummary = async () => {
      try {
        setLoading(true);
        const chunksData = await GetChunksPage(
          types.ListKind.STATS,
          0,
          100,
          { fields: ['range'], orders: [true] },
          '',
        );

        if (chunksData?.chunksStats) {
          const stats = chunksData.chunksStats.reduce(
            (acc, chunk) => {
              acc.totalChunks++;
              acc.totalAddresses += chunk.nAddrs || 0;
              acc.totalBlocks += chunk.nBlocks || 0;
              return acc;
            },
            { totalChunks: 0, totalAddresses: 0, totalBlocks: 0 },
          );

          setSummary({
            ...stats,
            indexHealth: stats.totalChunks > 0 ? 'healthy' : 'warning',
          });
        }
        setError(null);
      } catch (err) {
        Log(`Error fetching chunks summary: ${err}`);
        setError('Failed to load chunks');
        setSummary((prev) => ({ ...prev, indexHealth: 'error' }));
      } finally {
        setLoading(false);
      }
    };

    fetchChunksSummary();
  }, []);

  const formatNumber = (num: number) => {
    if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`;
    if (num >= 1000) return `${(num / 1000).toFixed(1)}K`;
    return num.toString();
  };

  return (
    <DashboardCard
      title="Chunks"
      subtitle="Index statistics"
      icon={<Chunks size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator
            status={summary.indexHealth}
            label="Index Health"
            count={summary.totalChunks}
            size="xs"
          />
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
          <Badge size="sm" variant="light" color="blue">
            Addresses: {formatNumber(summary.totalAddresses)}
          </Badge>
          <Badge size="sm" variant="light" color="green">
            Blocks: {formatNumber(summary.totalBlocks)}
          </Badge>
        </div>

        <Text size="xs" c="dimmed">
          View stats, index, blooms, and manifest data
        </Text>

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="outline" onClick={onViewAll}>
            View Details
          </Button>
        </Group>
      </Stack>
    </DashboardCard>
  );
};
