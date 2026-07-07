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
  Tooltip,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { IconExternalLink, IconFolderOpen, IconPhoto, IconX } from '@tabler/icons-react';
import { DetailHeader, usePersistedTab } from '@trueblocks/ui';
import { StatusBar, StatusLevel } from '../components/StatusBar';
import { DeleteImage, GetImageArtifactDataURL, GetTab, ListImages } from '../../wailsjs/go/app/App';
import { ExportImage, OpenImageArtifact, RevealImageArtifact } from '../../wailsjs/go/app/App';
import { RegenerateImage, SetTab } from '../../wailsjs/go/app/App';
import { dalle } from '../../wailsjs/go/models';

type ArtifactKind = 'annotated' | 'generated';

type ImagesProps = {
  selectedImageId?: string;
};

type StatusState = {
  visible: boolean;
  level: StatusLevel;
  message: string;
  meta?: string;
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
  const [status, setStatus] = useState<StatusState>({
    visible: false,
    level: 'progress',
    message: '',
  });
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
  const selectedTitle = selected ? displayTitle(selected) : '';
  const selectedIndex = selected
    ? images.findIndex((record) => recordKey(record) === selectedRecordId)
    : -1;
  const hasPrevious = selectedIndex > 0;
  const hasNext = selectedIndex >= 0 && selectedIndex < images.length - 1;

  const selectByOffset = useCallback(
    (offset: number) => {
      if (selectedIndex < 0) return;
      const next = images[selectedIndex + offset];
      if (next) setSelectedId(recordKey(next));
    },
    [images, selectedIndex],
  );

  const returnToGallery = useCallback(() => {
    setActiveTab('gallery');
  }, [setActiveTab]);

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
    const label = action === 'open' ? 'Opening image' : 'Showing image in Finder';
    setError('');
    setActionMessage('');
    setStatus({ visible: true, level: 'progress', message: label, meta: displayTitle(selected) });
    const operation = action === 'open' ? OpenImageArtifact : RevealImageArtifact;
    operation(selectedRecordId, artifact)
      .then(() => {
        const message = action === 'open' ? 'Opened image.' : 'Opened in Finder.';
        setActionMessage(message);
        setStatus({ visible: true, level: 'success', message, meta: displayTitle(selected) });
        notifications.show({ title: 'Image', message, color: 'green', autoClose: 2000 });
      })
      .catch((err: unknown) => {
        const message = messageFromError(err);
        setError(message);
        setStatus({ visible: true, level: 'error', message });
        notifications.show({
          title: 'Image action failed',
          message,
          color: 'red',
          autoClose: 5000,
        });
      });
  };

  const deleteImage = (id: string) => {
    if (!window.confirm('Delete this image and its metadata?')) return;
    const deleting = images.find((record) => recordKey(record) === id);
    const title = deleting ? displayTitle(deleting) : id;
    setError('');
    setActionMessage('');
    setImageActionId(id);
    setStatus({ visible: true, level: 'progress', message: 'Deleting image', meta: title });
    DeleteImage(id)
      .then(() => {
        const message = 'Deleted image and metadata.';
        setActionMessage(message);
        setStatus({ visible: true, level: 'success', message, meta: title });
        notifications.show({
          title: 'Image deleted',
          message: title,
          color: 'green',
          autoClose: 2500,
        });
        return load(series, '');
      })
      .catch((err: unknown) => {
        const message = messageFromError(err);
        setError(message);
        setStatus({ visible: true, level: 'error', message, meta: title });
        notifications.show({ title: 'Delete failed', message, color: 'red', autoClose: 5000 });
      })
      .finally(() => setImageActionId(''));
  };

  const refreshSelected = useCallback(() => {
    if (!selectedRecordId) return;
    setError('');
    setActionMessage('');
    setImageActionId(selectedRecordId);
    setStatus({
      visible: true,
      level: 'progress',
      message: 'Regenerating image',
      meta: selectedTitle,
    });
    RegenerateImage(selectedRecordId)
      .then((next) => {
        const message = 'Regenerated image from scratch.';
        setActionMessage(message);
        setStatus({
          visible: true,
          level: 'success',
          message,
          meta: next.metadata?.prompts?.titlePrompt || next.seed,
        });
        notifications.show({
          title: 'Image regenerated',
          message: next.metadata?.prompts?.titlePrompt || next.seed,
          color: 'green',
          autoClose: 2500,
        });
        return load(series, next.metadata?.imageId || next.seed);
      })
      .catch((err: unknown) => {
        const message = messageFromError(err);
        setError(message);
        setStatus({ visible: true, level: 'error', message, meta: selectedTitle });
        notifications.show({
          title: 'Regeneration failed',
          message,
          color: 'red',
          autoClose: 5000,
        });
      })
      .finally(() => setImageActionId(''));
  }, [load, selectedRecordId, selectedTitle, series]);

  const exportText = () => {
    if (!selectedRecordId) return;
    setError('');
    setActionMessage('');
    setStatus({
      visible: true,
      level: 'progress',
      message: 'Exporting text artifacts',
      meta: displayTitle(selected),
    });
    ExportImage(selectedRecordId, dalle.ExportImageOptions.createFrom({}))
      .then((result) => {
        const message = `Exported text files to ${result.dir}.`;
        setActionMessage(message);
        setStatus({ visible: true, level: 'success', message });
        notifications.show({
          title: 'Image text exported',
          message: result.dir,
          color: 'green',
          autoClose: 3000,
        });
      })
      .catch((err: unknown) => {
        const message = messageFromError(err);
        setError(message);
        setStatus({ visible: true, level: 'error', message });
        notifications.show({ title: 'Export failed', message, color: 'red', autoClose: 5000 });
      });
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
      if (activeTab !== 'detail') return;
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
  }, [activeTab, selectByOffset]);

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
          <Text fw={700} size="xl">
            Images
          </Text>
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
                <DetailHeader
                  hasPrev={hasPrevious}
                  hasNext={hasNext}
                  onPrev={() => selectByOffset(-1)}
                  onNext={() => selectByOffset(1)}
                  onBack={returnToGallery}
                  currentIndex={selectedIndex}
                  totalCount={images.length}
                  icon={<IconPhoto size={24} />}
                  title={<Text fw={700}>{displayTitle(selected)}</Text>}
                  subtitle={
                    <Text size="sm" c="dimmed">
                      {selected.metadata.series?.name} · {selected.metadata.seed}
                    </Text>
                  }
                  actionsRight={
                    <Group gap="xs" wrap="nowrap">
                      <Tooltip label="Open image">
                        <ActionIcon
                          variant="light"
                          disabled={!artifactURL}
                          onClick={() => runArtifactAction('open')}
                          aria-label="Open image"
                        >
                          <IconExternalLink size={18} />
                        </ActionIcon>
                      </Tooltip>
                      <Tooltip label="Show in Finder">
                        <ActionIcon
                          variant="light"
                          disabled={!artifactURL}
                          onClick={() => runArtifactAction('reveal')}
                          aria-label="Show in Finder"
                        >
                          <IconFolderOpen size={18} />
                        </ActionIcon>
                      </Tooltip>
                      <Button variant="light" onClick={exportText}>
                        Export Text
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
                  }
                />

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
      <StatusBar
        visible={status.visible}
        level={status.level}
        message={status.message}
        meta={status.meta}
      />
    </Stack>
  );
}
