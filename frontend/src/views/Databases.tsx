import { useCallback, useEffect, useState } from 'react';
import {
  Group,
  Paper,
  ScrollArea,
  SimpleGrid,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { ListDatabaseArchives, ListDatabaseRecords } from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';
import { storage } from '../../wailsjs/go/models';

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Databases() {
  const [archives, setArchives] = useState<storage.DatabaseArchiveManifest[]>([]);
  const [selectedName, setSelectedName] = useState('');
  const [records, setRecords] = useState<dalle.DatabaseRecordsResult | null>(null);
  const [filter, setFilter] = useState('');
  const [error, setError] = useState('');

  const files = archives.flatMap((archive) => archive.files.map((file) => ({ archive, file })));
  const selectedFile = files.find(({ file }) => file.name === selectedName)?.file;
  const visibleRecords = (records?.records ?? []).filter((record) => {
    const query = filter.trim().toLowerCase();
    if (!query) return true;
    return [record.key, ...record.values].some((value) => value.toLowerCase().includes(query));
  });

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
      <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="md">
        <Paper withBorder p="md">
          <ScrollArea h="calc(100vh - 250px)">
            <Table>
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>Name</Table.Th>
                  <Table.Th>Rows</Table.Th>
                  <Table.Th>Columns</Table.Th>
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody>
                {files.map(({ archive, file }) => (
                  <Table.Tr
                    key={`${archive.version}-${file.path}`}
                    onClick={() => loadRecords(file.name)}
                    style={{
                      cursor: 'pointer',
                      background:
                        file.name === selectedName ? 'var(--mantine-color-blue-light)' : undefined,
                    }}
                  >
                    <Table.Td>{file.name}</Table.Td>
                    <Table.Td>{file.rows}</Table.Td>
                    <Table.Td>{file.columns.join(', ')}</Table.Td>
                  </Table.Tr>
                ))}
              </Table.Tbody>
            </Table>
          </ScrollArea>
        </Paper>

        <Paper withBorder p="md">
          <Stack gap="md">
            <Group justify="space-between" align="end">
              <Stack gap={2}>
                <Title order={3}>{selectedName || 'Records'}</Title>
                <Text size="sm" c="dimmed">
                  {selectedFile?.rows ?? 0} rows · showing {visibleRecords.length}
                </Text>
              </Stack>
              <TextInput
                label="Filter records"
                value={filter}
                onChange={(event) => setFilter(event.currentTarget.value)}
              />
            </Group>
            <ScrollArea h="calc(100vh - 330px)">
              <Table>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th>Key</Table.Th>
                    {(selectedFile?.columns ?? []).map((column) => (
                      <Table.Th key={column}>{column}</Table.Th>
                    ))}
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                  {visibleRecords.map((record, index) => (
                    <Table.Tr key={`${record.key}-${index}`}>
                      <Table.Td>{record.key}</Table.Td>
                      {record.values.map((value, valueIndex) => (
                        <Table.Td key={`${record.key}-${valueIndex}`}>{value}</Table.Td>
                      ))}
                    </Table.Tr>
                  ))}
                </Table.Tbody>
              </Table>
            </ScrollArea>
          </Stack>
        </Paper>
      </SimpleGrid>
    </Stack>
  );
}
