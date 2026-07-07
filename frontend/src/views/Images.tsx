import { useCallback, useEffect, useState } from 'react';
import {
  ActionIcon,
  Box,
  Button,
  Group,
  Image,
  Paper,
  ScrollArea,
  SegmentedControl,
  SimpleGrid,
  Stack,
  Table,
  Tabs,
  Text,
  TextInput,
  Title,
  Tooltip,
} from '@mantine/core';
import { IconX } from '@tabler/icons-react';
import { usePersistedTab } from '@trueblocks/ui';
import { DeleteImage, GetImageArtifactDataURL, GetTab, ListImages } from '../../wailsjs/go/app/App';
import { ExportImage, OpenImageArtifact, RevealImageArtifact } from '../../wailsjs/go/app/App';
import { RegenerateImage, SetTab } from '../../wailsjs/go/app/App';
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

function thumbnailArtifact(record: dalle.ImageMetadataRecord): ArtifactKind | '' {
  if (record.metadata.artifacts?.annotated) return 'annotated';
  if (record.metadata.artifacts?.generated) return 'generated';
  return '';
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
  const [thumbnailURLs, setThumbnailURLs] = useState<Record<string, string>>({});
  const [error, setError] = useState('');
  const [actionMessage, setActionMessage] = useState('');
  const [imageActionId, setImageActionId] = useState('');
  const { activeTab, setActiveTab } = usePersistedTab({
    key: 'images',
    defaultTab: 'gallery',
    loadTab: GetTab,
    saveTab: SetTab,
    tabs: ['gallery', 'detail'],
    cycleViewId: 'images',
  });

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

  const refreshSelected = useCallback(() => {
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
  }, [load, selectedRecordId, series]);

  const exportText = () => {
    if (!selectedRecordId) return;
    setError('');
    setActionMessage('');
    ExportImage(selectedRecordId, dalle.ExportImageOptions.createFrom({}))
      .then((result) => setActionMessage(`Exported text files to ${result.dir}.`))
      .catch((err: unknown) => setError(messageFromError(err)));
  };

  useEffect(() => {
    if (selectedImageId) {
      setSeries('');
      setActiveTab('detail');
    }
    load(selectedImageId ? '' : series, selectedImageId);
  }, [load, selectedImageId, series, setActiveTab]);

  useEffect(() => {
    const handleRefresh = (event: Event) => {
      if ((event as CustomEvent).detail !== 'images') return;
      if (activeTab === 'detail' && selectedRecordId) {
        refreshSelected();
        return;
      }
      load(series, '');
    };

    window.addEventListener('view:refresh', handleRefresh);
    return () => window.removeEventListener('view:refresh', handleRefresh);
  }, [activeTab, load, refreshSelected, selectedRecordId, series]);

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

  useEffect(() => {
    let cancelled = false;
    setThumbnailURLs({});

    Promise.all(
      images.map(async (record) => {
        const key = recordKey(record);
        const thumbnail = thumbnailArtifact(record);
        if (!thumbnail) return [key, ''] as const;
        try {
          return [key, await GetImageArtifactDataURL(key, thumbnail)] as const;
        } catch {
          return [key, ''] as const;
        }
      }),
    ).then((entries) => {
      if (cancelled) return;
      setThumbnailURLs(
        Object.fromEntries(entries.filter((entry) => Boolean(entry[1]))) as Record<string, string>,
      );
    });

    return () => {
      cancelled = true;
    };
  }, [images]);

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
        </Group>
      </Group>

      {error && <Text c="red">{error}</Text>}
      {actionMessage && (
        <Text c="dimmed" size="sm">
          {actionMessage}
        </Text>
      )}

      <Tabs value={activeTab} onChange={(value) => value && setActiveTab(value)}>
        <Tabs.List>
          <Tabs.Tab value="gallery">Gallery</Tabs.Tab>
          <Tabs.Tab value="detail" disabled={!selected}>
            Detail
          </Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="gallery" pt="md">
          <ScrollArea h="calc(100vh - 230px)">
            <SimpleGrid cols={{ base: 2, sm: 3, md: 4, xl: 5 }} spacing="xs">
              {images.map((record) => {
                const key = recordKey(record);
                const isSelected = key === selectedRecordId;
                const thumbnailURL = thumbnailURLs[key];
                return (
                  <Paper
                    key={key}
                    withBorder
                    p={4}
                    style={{
                      cursor: 'pointer',
                      borderColor: isSelected ? 'var(--mantine-color-blue-6)' : undefined,
                    }}
                    onClick={() => {
                      setSelectedId(key);
                      setActiveTab('detail');
                    }}
                  >
                    <Tooltip label={displayTitle(record)}>
                      <Box
                        pos="relative"
                        style={{
                          aspectRatio: '1 / 1',
                          overflow: 'hidden',
                          borderRadius: 'var(--mantine-radius-sm)',
                        }}
                      >
                        {thumbnailURL ? (
                          <Image src={thumbnailURL} h="100%" w="100%" fit="cover" />
                        ) : (
                          <Paper h="100%" bg="gray.0" radius="sm">
                            <Stack h="100%" align="center" justify="center" gap={2}>
                              <Text size="xs" c="dimmed">
                                No image
                              </Text>
                            </Stack>
                          </Paper>
                        )}
                        <ActionIcon
                          aria-label="Delete image and metadata"
                          variant="filled"
                          color="red"
                          size="sm"
                          loading={imageActionId === key}
                          style={{ position: 'absolute', top: 6, right: 6 }}
                          onClick={(event) => {
                            event.stopPropagation();
                            deleteImage(key);
                          }}
                        >
                          <IconX size={13} />
                        </ActionIcon>
                      </Box>
                    </Tooltip>
                  </Paper>
                );
              })}
            </SimpleGrid>
          </ScrollArea>
        </Tabs.Panel>

        <Tabs.Panel value="detail" pt="md">
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
                  <Image src={artifactURL} radius="sm" fit="contain" mah="60vh" />
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
              <Text c="dimmed">No image selected.</Text>
            )}
          </Paper>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
