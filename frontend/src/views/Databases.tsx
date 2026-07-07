import { useEffect, useState } from 'react';
import { Button, Group, Stack, Table, Text, Title } from '@mantine/core';
import { ListDatabaseArchives } from '../../wailsjs/go/app/App';
import { storage } from '../../wailsjs/go/models';

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Databases() {
  const [archives, setArchives] = useState<storage.DatabaseArchiveManifest[]>([]);
  const [error, setError] = useState('');

  const load = () => {
    setError('');
    ListDatabaseArchives()
      .then((items) => setArchives(items ?? []))
      .catch((err: unknown) => setError(messageFromError(err)));
  };

  useEffect(load, []);

  const files = archives.flatMap((archive) => archive.files.map((file) => ({ archive, file })));

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>Databases</Title>
        <Button onClick={load}>Refresh</Button>
      </Group>
      {error && <Text c="red">{error}</Text>}
      {archives.map((archive) => (
        <Text key={archive.version} size="sm" c="dimmed">
          {archive.version} · {archive.archiveHash}
        </Text>
      ))}
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
            <Table.Tr key={`${archive.version}-${file.path}`}>
              <Table.Td>{file.name}</Table.Td>
              <Table.Td>{file.rows}</Table.Td>
              <Table.Td>{file.columns.join(', ')}</Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Stack>
  );
}