import { useEffect, useState } from 'react';

import { GetMonitorsPage } from '@app';
import { DashboardCard, StatusIndicator } from '@components';
import { useEvent, useIcons } from '@hooks';
import { Badge, Button, Group, Stack } from '@mantine/core';
import { msgs, types } from '@models';
import { Log } from '@utils';

interface MonitorsSummary {
  total: number;
  active: number;
  staged: number;
  deleted: number;
  empty: number;
}

interface MonitorsPanelProps {
  onViewAll?: () => void;
  onAddMonitor?: () => void;
}

export const MonitorsPanel = ({
  onViewAll,
  onAddMonitor,
}: MonitorsPanelProps) => {
  const [summary, setSummary] = useState<MonitorsSummary>({
    total: 0,
    active: 0,
    staged: 0,
    deleted: 0,
    empty: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { Monitors } = useIcons();

  const fetchMonitorsSummary = async (showLoading = true) => {
    try {
      if (showLoading) setLoading(true);
      const monitorsData = await GetMonitorsPage(
        types.ListKind.MONITORS,
        0,
        1000,
        { fields: ['address'], orders: [true] },
        '',
      );

      if (monitorsData?.monitors) {
        const summary = monitorsData.monitors.reduce(
          (acc, monitor) => {
            acc.total++;
            if (monitor.deleted) acc.deleted++;
            else if (monitor.isStaged) acc.staged++;
            else acc.active++;

            if (monitor.isEmpty) acc.empty++;
            return acc;
          },
          { total: 0, active: 0, staged: 0, deleted: 0, empty: 0 },
        );
        setSummary(summary);
      }
      setError(null);
    } catch (err) {
      Log(`Error fetching monitors summary: ${err}`);
      setError('Failed to load monitors');
    } finally {
      if (showLoading) setLoading(false);
    }
  };

  useEffect(() => {
    fetchMonitorsSummary();
  }, []);

  // Listen for data changes and refresh silently (no loading state)
  useEvent(
    msgs.EventType.DATA_LOADED,
    (_message: string, payload?: types.DataLoadedPayload) => {
      if (payload?.listKind === types.ListKind.MONITORS) {
        fetchMonitorsSummary(false);
      }
    },
  );

  const getHealthStatus = () => {
    if (summary.total === 0) return 'inactive';
    if (summary.deleted > summary.active) return 'warning';
    return 'healthy';
  };

  return (
    <DashboardCard
      title="Monitors"
      subtitle={`${summary.total} addresses`}
      icon={<Monitors size={20} />}
      loading={loading}
      error={error}
      onClick={onViewAll}
    >
      <Stack gap="sm">
        <div>
          <StatusIndicator
            status={getHealthStatus()}
            label="Monitor Health"
            count={summary.active}
            size="xs"
          />
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
          {summary.active > 0 && (
            <Badge size="sm" variant="light" color="green">
              Active: {summary.active}
            </Badge>
          )}
          {summary.staged > 0 && (
            <Badge size="sm" variant="light" color="blue">
              Staged: {summary.staged}
            </Badge>
          )}
          {summary.deleted > 0 && (
            <Badge size="sm" variant="light" color="red">
              Deleted: {summary.deleted}
            </Badge>
          )}
          {summary.empty > 0 && (
            <Badge size="sm" variant="light" color="gray">
              Empty: {summary.empty}
            </Badge>
          )}
        </div>

        <Group gap="xs" mt="auto">
          <Button size="xs" variant="light" onClick={onAddMonitor}>
            Add Monitor
          </Button>
          <Button size="xs" variant="outline" onClick={onViewAll}>
            View All
          </Button>
        </Group>
      </Stack>
    </DashboardCard>
  );
};
