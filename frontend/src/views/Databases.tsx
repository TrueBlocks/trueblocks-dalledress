import { useCallback, useEffect, useMemo, useState } from 'react';
import { Group, Stack, Tabs, Text, Title } from '@mantine/core';
import { Column, DataTable, usePersistedTab } from '@trueblocks/ui';
import {
  GetTab,
  ListDatabaseArchives,
  ListDatabaseRecords,
  SetTab,
} from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';
import { storage } from '../../wailsjs/go/models';

type DatabaseFileRow = {
  archive: storage.DatabaseArchiveManifest;
  file: storage.DatabaseFileManifest;
};

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Databases() {
  const [archives, setArchives] = useState<storage.DatabaseArchiveManifest[]>([]);
  const [selectedName, setSelectedName] = useState('');
  const [records, setRecords] = useState<dalle.DatabaseRecordsResult | null>(null);
  const [error, setError] = useState('');
  const { activeTab, setActiveTab } = usePersistedTab({
    key: 'databases',
    defaultTab: 'list',
    loadTab: GetTab,
    saveTab: SetTab,
    tabs: ['list', 'detail'],
    cycleViewId: 'databases',
  });

  const files = archives.flatMap((archive) => archive.files.map((file) => ({ archive, file })));
  const selectedFile = files.find(({ file }) => file.name === selectedName)?.file;

  const fileColumns: Column<DatabaseFileRow>[] = useMemo(
    () => [
      {
        key: 'name',
        label: 'Name',
        width: '30%',
        render: ({ file }) => file.name,
        sortValue: ({ file }) => file.name.toLowerCase(),
        scrollOnSelect: true,
      },
      {
        key: 'rows',
        label: 'Rows',
        width: '12%',
        render: ({ file }) => file.rows.toLocaleString(),
        sortValue: ({ file }) => file.rows,
      },
      {
        key: 'columns',
        label: 'Columns',
        width: '43%',
        render: ({ file }) => file.columns.join(', '),
        sortValue: ({ file }) => file.columns.length,
      },
      {
        key: 'version',
        label: 'Version',
        width: '15%',
        render: ({ archive }) => archive.version,
        sortValue: ({ archive }) => archive.version,
      },
    ],
    [],
  );

  const recordColumns: Column<storage.DatabaseRecord>[] = useMemo(() => {
    const valueColumns = (selectedFile?.columns ?? []).map((column, index) => ({
      key: `value-${index}`,
      label: column,
      render: (record: storage.DatabaseRecord) => record.values[index] ?? '',
      sortValue: (record: storage.DatabaseRecord) => record.values[index] ?? '',
    }));
    return [
      {
        key: 'key',
        label: 'Key',
        width: '20%',
        render: (record: storage.DatabaseRecord) => record.key,
        sortValue: (record: storage.DatabaseRecord) => record.key,
        scrollOnSelect: true,
      },
      ...valueColumns,
    ];
  }, [selectedFile]);

  const loadRecords = useCallback((name: string) => {
    setSelectedName(name);
    setError('');
    ListDatabaseRecords(name, 200)
      .then(setRecords)
      .catch((err: unknown) => setError(messageFromError(err)));
  }, []);

  const load = useCallback(
    (preferred = '') => {
      setError('');
      ListDatabaseArchives()
        .then((items) => {
          const next = items ?? [];
          setArchives(next);
          const nextFiles = next.flatMap((archive) => archive.files.map((file) => file.name));
          const nextSelected = nextFiles.includes(preferred) ? preferred : nextFiles[0] || '';
          if (nextSelected) loadRecords(nextSelected);
        })
        .catch((err: unknown) => setError(messageFromError(err)));
    },
    [loadRecords],
  );

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    const handleRefresh = (event: Event) => {
      if ((event as CustomEvent).detail !== 'databases') return;
      load(selectedName);
    };

    window.addEventListener('view:refresh', handleRefresh);
    return () => window.removeEventListener('view:refresh', handleRefresh);
  }, [load, selectedName]);

  const selectDatabase = (row: DatabaseFileRow) => {
    loadRecords(row.file.name);
    setActiveTab('detail');
  };

  const searchFiles = (row: DatabaseFileRow, search: string) => {
    const query = search.toLowerCase();
    return [row.file.name, row.archive.version, row.file.path, ...row.file.columns].some((value) =>
      value.toLowerCase().includes(query),
    );
  };

  const searchRecords = (record: storage.DatabaseRecord, search: string) => {
    const query = search.toLowerCase();
    return [record.key, ...record.values].some((value) => value.toLowerCase().includes(query));
  };

  return (
    <Stack gap="md">
      <Group justify="space-between" align="end">
        <Stack gap={2}>
          <Title order={2}>Databases</Title>
          <Text c="dimmed" size="sm">
            {files.length} tables
          </Text>
        </Stack>
      </Group>
      {error && <Text c="red">{error}</Text>}
      {archives.map((archive) => (
        <Text key={archive.version} size="sm" c="dimmed">
          {archive.version} · {archive.archiveHash}
        </Text>
      ))}
      <Tabs value={activeTab} onChange={(value) => value && setActiveTab(value)}>
        <Tabs.List>
          <Tabs.Tab value="list">List</Tabs.Tab>
          <Tabs.Tab value="detail" disabled={!selectedName}>
            Detail
          </Tabs.Tab>
        </Tabs.List>
        <Tabs.Panel value="list" pt="md">
          <DataTable<DatabaseFileRow>
            tableName="dalle-databases"
            data={files}
            columns={fileColumns}
            getRowKey={({ archive, file }) => `${archive.version}-${file.path}`}
            onRowClick={selectDatabase}
            searchFn={searchFiles}
          />
        </Tabs.Panel>
        <Tabs.Panel value="detail" pt="md">
          <Stack gap="xs">
            <Text fw={700}>{selectedName || 'Records'}</Text>
            <Text size="sm" c="dimmed">
              {selectedFile?.rows.toLocaleString() ?? 0} rows · showing first{' '}
              {(records?.records ?? []).length.toLocaleString()}
            </Text>
            <DataTable<storage.DatabaseRecord>
              tableName={`dalle-database-${selectedName || 'records'}`}
              data={records?.records ?? []}
              columns={recordColumns}
              getRowKey={(record) => record.key}
              searchFn={searchRecords}
            />
          </Stack>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
