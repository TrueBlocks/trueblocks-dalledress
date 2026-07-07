import { useCallback, useEffect, useState } from 'react';
import {
  ActionIcon,
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
  Tooltip,
} from '@mantine/core';
import { IconRefresh, IconX } from '@tabler/icons-react';
import { DeleteImage, GetImageArtifactDataURL, ListImages } from '../../wailsjs/go/app/App';
import { ExportImage, OpenImageArtifact, RevealImageArtifact } from '../../wailsjs/go/app/App';
import { RegenerateImage } from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';

type ArtifactKind = 'annotated' | 'generated';

type ImagesProps = {
  selectedImageId?: string;
};

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

function selectRecordId(
  records: dalle.ImageMetadataRecord[],
  current: string,
  preferred: string,
): string {
  if (preferred && records.some((record) => recordKey(record) === preferred)) return preferred;
  if (current && records.some((record) => recordKey(record) === current)) return current;
  return records[0] ? recordKey(records[0]) : '';
}

export function Images({ selectedImageId = '' }: ImagesProps) {
  const [series, setSeries] = useState('');
  const [images, setImages] = useState<dalle.ImageMetadataRecord[]>([]);
  const [selectedId, setSelectedId] = useState('');
  const [artifact, setArtifact] = useState<ArtifactKind>('annotated');
  const [artifactURL, setArtifactURL] = useState('');
  const [error, setError] = useState('');
  const [actionMessage, setActionMessage] = useState('');
  const [imageActionId, setImageActionId] = useState('');

  const selected = images.find((record) => recordKey(record) === selectedId) ?? images[0];
  const selectedRecordId = selected ? recordKey(selected) : '';
  const selectedIndex = selected
    ? images.findIndex((record) => recordKey(record) === selectedRecordId)
    : -1;

  const selectByOffset = useCallback(
    (offset: number) => {
      if (selectedIndex < 0) return;
      const next = images[selectedIndex + offset];
      if (next) setSelectedId(recordKey(next));
    },
    [images, selectedIndex],
  );

  const load = useCallback(
    (seriesFilter = series, preferredId = selectedImageId) => {
      setError('');
      return ListImages(seriesFilter)
        .then((items) => {
          const next = items ?? [];
          setImages(next);
          setSelectedId((current) => selectRecordId(next, current, preferredId));
        })
        .catch((err: unknown) => setError(messageFromError(err)));
    },
    [selectedImageId, series],
  );

  const runArtifactAction = (action: 'open' | 'reveal') => {
    if (!selectedRecordId) return;
    setError('');
    setActionMessage('');
    const operation = action === 'open' ? OpenImageArtifact : RevealImageArtifact;
    operation(selectedRecordId, artifact)
      .then(() => setActionMessage(action === 'open' ? 'Opened image.' : 'Opened in Finder.'))
      .catch((err: unknown) => setError(messageFromError(err)));
  };

  const deleteImage = (id: string) => {
    if (!window.confirm('Delete this image and its metadata?')) return;
    setError('');
    setActionMessage('');
    setImageActionId(id);
    DeleteImage(id)
      .then(() => {
        setActionMessage('Deleted image and metadata.');
        return load(series, '');
      })
      .catch((err: unknown) => setError(messageFromError(err)))
      .finally(() => setImageActionId(''));
  };

  const refreshSelected = () => {
    if (!selectedRecordId) return;
    setError('');
    setActionMessage('');
    setImageActionId(selectedRecordId);
    RegenerateImage(selectedRecordId)
      .then((next) => {
        setActionMessage('Regenerated image from scratch.');
        return load(series, next.metadata?.imageId || next.seed);
      })
      .catch((err: unknown) => setError(messageFromError(err)))
      .finally(() => setImageActionId(''));
  };

  const exportText = () => {
    if (!selectedRecordId) return;
    setError('');
    setActionMessage('');
    ExportImage(selectedRecordId, dalle.ExportImageOptions.createFrom({}))
      .then((result) => setActionMessage(`Exported text files to ${result.dir}.`))
      .catch((err: unknown) => setError(messageFromError(err)));
  };

  useEffect(() => {
    if (selectedImageId) setSeries('');
    load(selectedImageId ? '' : series, selectedImageId);
  }, [load, selectedImageId, series]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.metaKey || event.ctrlKey || event.altKey) return;
      if (event.target instanceof HTMLElement) {
        const editableTags = ['INPUT', 'TEXTAREA', 'SELECT', 'BUTTON'];
        if (editableTags.includes(event.target.tagName)) return;
      }
      if (event.key === 'ArrowLeft') {
        event.preventDefault();
        selectByOffset(-1);
      }
      if (event.key === 'ArrowRight') {
        event.preventDefault();
        selectByOffset(1);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [selectByOffset]);

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
          <Button variant="light" onClick={() => load()} disabled={Boolean(imageActionId)}>
            Reload
          </Button>
          <Button
            leftSection={<IconRefresh size={16} />}
            onClick={refreshSelected}
            loading={Boolean(imageActionId && imageActionId === selectedRecordId)}
            disabled={!selectedRecordId}
          >
            Refresh
          </Button>
        </Group>
      </Group>

      {error && <Text c="red">{error}</Text>}
      {actionMessage && (
        <Text c="dimmed" size="sm">
          {actionMessage}
        </Text>
      )}

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
                      <Group gap="xs" wrap="nowrap">
                        <Badge variant="light">{statusLabel(record)}</Badge>
                        <Tooltip label="Delete image and metadata">
                          <ActionIcon
                            aria-label="Delete image and metadata"
                            variant="subtle"
                            color="red"
                            loading={imageActionId === key}
                            onClick={(event) => {
                              event.stopPropagation();
                              deleteImage(key);
                            }}
                          >
                            <IconX size={14} />
                          </ActionIcon>
                        </Tooltip>
                      </Group>
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
                  <Text size="xs" c="dimmed">
                    {selectedIndex + 1} of {images.length}
                  </Text>
                </Stack>
                <Group>
                  <Button
                    variant="light"
                    disabled={!artifactURL}
                    onClick={() => runArtifactAction('open')}
                  >
                    Open
                  </Button>
                  <Button
                    variant="light"
                    disabled={!artifactURL}
                    onClick={() => runArtifactAction('reveal')}
                  >
                    Show in Finder
                  </Button>
                  <Button variant="light" onClick={exportText}>
                    Export Text
                  </Button>
                  <Button
                    variant="light"
                    disabled={selectedIndex <= 0}
                    onClick={() => selectByOffset(-1)}
                  >
                    Previous
                  </Button>
                  <Button
                    variant="light"
                    disabled={selectedIndex < 0 || selectedIndex >= images.length - 1}
                    onClick={() => selectByOffset(1)}
                  >
                    Next
                  </Button>
                  <SegmentedControl
                    value={artifact}
                    onChange={(value) => setArtifact(value as ArtifactKind)}
                    data={[
                      { value: 'annotated', label: 'Annotated' },
                      { value: 'generated', label: 'Generated' },
                    ]}
                  />
                </Group>
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
