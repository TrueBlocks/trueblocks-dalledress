import { useEffect, useState } from 'react';
import { Button, Group, Stack, Table, Text, TextInput, Title } from '@mantine/core';
import { ListImages } from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Images() {
  const [series, setSeries] = useState('');
  const [images, setImages] = useState<dalle.ImageMetadataRecord[]>([]);
  const [error, setError] = useState('');

  const load = () => {
    setError('');
    ListImages(series)
      .then((items) => setImages(items ?? []))
      .catch((err: unknown) => setError(messageFromError(err)));
  };

  useEffect(() => {
    setError('');
    ListImages('')
      .then((items) => setImages(items ?? []))
      .catch((err: unknown) => setError(messageFromError(err)));
  }, []);

  return (
    <Stack>
      <Title order={2}>Images</Title>
      <Group align="end">
        <TextInput
          label="Series filter"
          value={series}
          onChange={(event) => setSeries(event.currentTarget.value)}
        />
        <Button onClick={load}>Refresh</Button>
      </Group>
      {error && <Text c="red">{error}</Text>}
      <Table>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Seed</Table.Th>
            <Table.Th>Series</Table.Th>
            <Table.Th>Input</Table.Th>
            <Table.Th>Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {images.map((record) => (
            <Table.Tr key={record.metadata.imageId || record.path}>
              <Table.Td>{record.metadata.seed}</Table.Td>
              <Table.Td>{record.metadata.series?.name}</Table.Td>
              <Table.Td>{record.metadata.input}</Table.Td>
              <Table.Td>{record.metadata.status?.completed ? 'complete' : 'pending'}</Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Stack>
  );
}
