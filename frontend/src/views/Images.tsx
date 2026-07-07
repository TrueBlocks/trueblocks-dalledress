import { useEffect, useState } from 'react';
import {
  Badge,
  Button,
  Group,
  Image,
  Paper,
  ScrollArea,
  SegmentedControl,
  SimpleGrid,
  Stack,
  Table,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { GetImageArtifactDataURL, ListImages } from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';

type ArtifactKind = 'annotated' | 'generated';

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

function recordKey(record: dalle.ImageMetadataRecord): string {
  return record.metadata.imageId || record.metadata.seed || record.path;
}

function displayTitle(record: dalle.ImageMetadataRecord): string {
  return record.metadata.prompts?.titlePrompt || record.metadata.input || record.metadata.seed;
}

function statusLabel(record: dalle.ImageMetadataRecord): string {
  if (record.metadata.artifacts?.annotated) return 'annotated';
  if (record.metadata.artifacts?.generated) return 'generated';
  if (record.metadata.status?.completed) return 'metadata';
  return 'pending';
}

export function Images() {
  const [series, setSeries] = useState('');
  const [images, setImages] = useState<dalle.ImageMetadataRecord[]>([]);
  const [selectedId, setSelectedId] = useState('');
  const [artifact, setArtifact] = useState<ArtifactKind>('annotated');
  const [artifactURL, setArtifactURL] = useState('');
  const [error, setError] = useState('');

  const selected = images.find((record) => recordKey(record) === selectedId) ?? images[0];
  const selectedRecordId = selected ? recordKey(selected) : '';

  const load = () => {
    setError('');
    ListImages(series)
      .then((items) => {
        const next = items ?? [];
        setImages(next);
        setSelectedId((current) => {
          if (next.some((record) => recordKey(record) === current)) return current;
          return next[0] ? recordKey(next[0]) : '';
        });
      })
      .catch((err: unknown) => setError(messageFromError(err)));
  };

  useEffect(() => {
    setError('');
    ListImages('')
      .then((items) => {
        const next = items ?? [];
        setImages(next);
        setSelectedId(next[0] ? recordKey(next[0]) : '');
      })
      .catch((err: unknown) => setError(messageFromError(err)));
  }, []);

  useEffect(() => {
    if (!selectedRecordId) {
      setArtifactURL('');
      return;
    }
    const preferredArtifact =
      artifact === 'annotated' && !selected?.metadata.artifacts?.annotated ? 'generated' : artifact;
    GetImageArtifactDataURL(selectedRecordId, preferredArtifact)
      .then(setArtifactURL)
      .catch(() => setArtifactURL(''));
  }, [artifact, selected, selectedRecordId]);

  return (
    <Stack gap="md">
      <Group justify="space-between" align="end">
        <Stack gap={2}>
          <Title order={2}>Images</Title>
          <Text c="dimmed" size="sm">
            {images.length} records
          </Text>
        </Stack>
        <Group align="end">
          <TextInput
            label="Series filter"
            value={series}
            onChange={(event) => setSeries(event.currentTarget.value)}
          />
          <Button onClick={load}>Refresh</Button>
        </Group>
      </Group>

      {error && <Text c="red">{error}</Text>}

      <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="md">
        <ScrollArea h="calc(100vh - 220px)">
          <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="sm">
            {images.map((record) => {
              const key = recordKey(record);
              const isSelected = key === selectedRecordId;
              return (
                <Paper
                  key={key}
                  withBorder
                  p="sm"
                  style={{
                    cursor: 'pointer',
                    borderColor: isSelected ? 'var(--mantine-color-blue-6)' : undefined,
                  }}
                  onClick={() => setSelectedId(key)}
                >
                  <Stack gap="xs">
                    <Group justify="space-between" align="start">
                      <Text fw={700} lineClamp={2}>
                        {displayTitle(record)}
                      </Text>
                      <Badge variant="light">{statusLabel(record)}</Badge>
                    </Group>
                    <Text size="xs" c="dimmed" lineClamp={1}>
                      {record.metadata.series?.name} · {record.metadata.seed}
                    </Text>
                    <Text size="sm" lineClamp={2}>
                      {record.metadata.input}
                    </Text>
                  </Stack>
                </Paper>
              );
            })}
          </SimpleGrid>
        </ScrollArea>

        <Paper withBorder p="md">
          {selected ? (
            <Stack gap="md">
              <Group justify="space-between" align="start">
                <Stack gap={2}>
                  <Title order={3}>{displayTitle(selected)}</Title>
                  <Text size="sm" c="dimmed">
                    {selected.metadata.series?.name} · {selected.metadata.seed}
                  </Text>
                </Stack>
                <SegmentedControl
                  value={artifact}
                  onChange={(value) => setArtifact(value as ArtifactKind)}
                  data={[
                    { value: 'annotated', label: 'Annotated' },
                    { value: 'generated', label: 'Generated' },
                  ]}
                />
              </Group>

              {artifactURL ? (
                <Image src={artifactURL} radius="sm" fit="contain" mah="52vh" />
              ) : (
                <Paper withBorder p="xl">
                  <Text c="dimmed">No {artifact} artifact is available for this image.</Text>
                </Paper>
              )}

              <SimpleGrid cols={{ base: 1, md: 2 }} spacing="sm">
                <Paper withBorder p="sm">
                  <Text fw={700} size="sm">
                    Prompt
                  </Text>
                  <Text size="sm">
                    {selected.metadata.prompts?.prompt || 'No prompt recorded.'}
                  </Text>
                </Paper>
                <Paper withBorder p="sm">
                  <Text fw={700} size="sm">
                    Artifacts
                  </Text>
                  <Text size="xs" c="dimmed">
                    Generated: {selected.metadata.artifacts?.generated || 'none'}
                  </Text>
                  <Text size="xs" c="dimmed">
                    Annotated: {selected.metadata.artifacts?.annotated || 'none'}
                  </Text>
                  <Text size="xs" c="dimmed">
                    Metadata: {selected.path}
                  </Text>
                </Paper>
              </SimpleGrid>

              <Table>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th>Attribute</Table.Th>
                    <Table.Th>Database</Table.Th>
                    <Table.Th>Record</Table.Th>
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                  {(selected.metadata.selectedRecords ?? []).map((record) => (
                    <Table.Tr key={`${record.attribute}-${record.database}-${record.rowIndex}`}>
                      <Table.Td>{record.attribute}</Table.Td>
                      <Table.Td>{record.database}</Table.Td>
                      <Table.Td>{record.record}</Table.Td>
                    </Table.Tr>
                  ))}
                </Table.Tbody>
              </Table>
            </Stack>
          ) : (
            <Text c="dimmed">No images found.</Text>
          )}
        </Paper>
      </SimpleGrid>
    </Stack>
  );
}
