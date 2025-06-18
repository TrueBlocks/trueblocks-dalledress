import { useEffect, useState } from 'react';

import { GetNamesPage } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useEvent, useIcons } from '@hooks';
import { Badge, Button, Group, Stack } from '@mantine/core';
import { msgs, types } from '@models';
import { Log } from '@utils';

interface NamesSummary {
  total: number;
  custom: number;
  prefund: number;
  regular: number;
  baddress: number;
}

interface NamesPanelProps {
  onViewAll?: () => void;
  onAddName?: () => void;
}

export const NamesPanel = ({ onViewAll, onAddName }: NamesPanelProps) => {
  const [summary, setSummary] = useState<NamesSummary>({
    total: 0,
    custom: 0,
    prefund: 0,
    regular: 0,
    baddress: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { Names } = useIcons();

  const fetchNamesSummary = async (showLoading = true) => {
    try {
      if (showLoading) setLoading(true);
      const allNames = await GetNamesPage(
        types.ListKind.ALL,
        0,
        200000,
        { fields: ['address'], orders: [true] },
        '',
      );

      if (allNames?.names) {
        const summary = allNames.names.reduce(
          (acc, name) => {
            acc.total++;
            if (name.isCustom || name.tags?.includes('Custom')) acc.custom++;
            else if (name.isPrefund || name.tags?.includes('Prefund'))
              acc.prefund++;
            else if (name.tags?.includes('Baddress')) acc.baddress++;
            else acc.regular++;
            return acc;
          },
          { total: 0, custom: 0, prefund: 0, regular: 0, baddress: 0 },
        );
        setSummary(summary);
      }
      setError(null);
    } catch (err) {
      Log(`Error fetching names summary: ${err}`);
      setError('Failed to load names');
    } finally {
      if (showLoading) setLoading(false);
    }
  };

  useEffect(() => {
    fetchNamesSummary();
  }, []);

  // Listen for data changes and refresh silently (no loading state)
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: types.DataLoadedPayload) => {
      // Update when any name-related data is loaded
      if (
        payload?.listKind === types.ListKind.ALL ||
        payload?.listKind === types.ListKind.CUSTOM ||
        payload?.listKind === types.ListKind.PREFUND ||
        payload?.listKind === types.ListKind.REGULAR ||
        payload?.listKind === types.ListKind.BADDRESS
      ) {
        fetchNamesSummary(false);
      }
    },
  );

  return (
    <DashboardCard
      title="Names"
      subtitle={`${summary.total} addresses`}
      icon={<Names size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator
            status={summary.total > 0 ? 'healthy' : 'inactive'}
            label="Name Database"
            count={summary.total}
            size="xs"
          />
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
          {summary.custom > 0 && (
            <Badge size="sm" variant="light" color="blue">
              Custom: {summary.custom}
            </Badge>
          )}
          {summary.regular > 0 && (
            <Badge size="sm" variant="light" color="green">
              Regular: {summary.regular}
            </Badge>
          )}
          {summary.prefund > 0 && (
            <Badge size="sm" variant="light" color="orange">
              Prefund: {summary.prefund}
            </Badge>
          )}
          {summary.baddress > 0 && (
            <Badge size="sm" variant="light" color="red">
              Bad: {summary.baddress}
            </Badge>
          )}
        </div>

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="light" onClick={onAddName}>
            Add Name
          </Button>
          <Button size="xs" variant="outline" onClick={onViewAll}>
            View All
          </Button>
        </Group>
      </Stack>
    </DashboardCard>
  );
};
