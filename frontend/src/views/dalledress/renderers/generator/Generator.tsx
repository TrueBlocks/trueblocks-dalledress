import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { GetDalleDressCurrent } from '@app';
import { useActiveProject, useIconSets } from '@hooks';
import {
  Button,
  Center,
  Container,
  Group,
  Image,
  ScrollArea,
  Select,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { dalle, dalledress, types } from '@models';
import { Log } from '@utils';

import { DalleDressCard } from '../../components';
import { useSpeakPrompt } from '../../hooks/useSpeakPrompt';
import { getItemKey, useGalleryStore } from '../../store';

export function Generator({
  pageData,
}: {
  pageData: dalledress.DalleDressPage | null;
}) {
  const { activeAddress, setActiveAddress } = useActiveProject();
  const [selectedSeries, setSelectedSeries] = useState<string | null>('empty');
  const [current, setCurrent] = useState<dalle.DalleDress | null>(null);
  const { speaking, audioUrl, audioRef, speak } = useSpeakPrompt({
    activeAddress: activeAddress || null,
    selectedSeries,
    hasEnhancedPrompt: !!current?.enhancedPrompt,
  });
  const thumbRowRef = useRef<HTMLDivElement | null>(null);
  const icons = useIconSets();
  const SpeakIcon = icons.Speak;

  const {
    orig,
    series,
    getSelectionKey,
    setSelection,
    clearSelection,
    ingestItems,
    galleryItems,
  } = useGalleryStore();

  useEffect(() => {
    ingestItems(pageData?.dalledress || []);
  }, [pageData?.dalledress, ingestItems]);

  const seriesOptions = useMemo(() => {
    const set = new Set<string>();
    set.add('empty');
    galleryItems.forEach((g) => {
      if (g.series) set.add(g.series);
    });
    const arr = Array.from(set).sort((a, b) => a.localeCompare(b));
    if (arr[0] !== 'empty') {
      const rest = arr.filter((s) => s !== 'empty');
      return ['empty', ...rest];
    }
    return arr;
  }, [galleryItems]);

  const addressOptions = useMemo(() => {
    const set = new Set<string>();
    if (activeAddress) set.add(activeAddress);
    galleryItems.forEach((g) => set.add(g.original));
    return Array.from(set).sort((a, b) => a.localeCompare(b));
  }, [galleryItems, activeAddress]);

  useEffect(() => {
    if (!activeAddress && addressOptions.length) {
      const first = addressOptions[0];
      if (first) setActiveAddress(first);
    }
  }, [activeAddress, addressOptions, setActiveAddress]);

  useEffect(() => {
    if (!selectedSeries) {
      const first = galleryItems.find(
        (g) => g.original === activeAddress && g.series,
      );
      if (first?.series) setSelectedSeries(first.series);
      if (!selectedSeries && first) setSelectedSeries(first.series);
    }
  }, [selectedSeries, galleryItems, activeAddress]);

  const handleAddressChange = useCallback(
    (value: string | null) => {
      if (!value) return;
      setActiveAddress(value);
    },
    [setActiveAddress],
  );

  const handleSeriesChange = useCallback((value: string | null) => {
    setSelectedSeries(value);
  }, []);

  const handleGenerate = useCallback(() => {
    if (!activeAddress || !selectedSeries) return;
  }, [activeAddress, selectedSeries]);

  const handleThumbSelect = useCallback(
    (item: dalle.DalleDress) => {
      const key = getItemKey(item);
      setSelection(key);
      if (item.series) setSelectedSeries(item.series);
    },
    [setSelection],
  );

  const handleThumbDouble = useCallback(
    (item: dalle.DalleDress) => {
      handleThumbSelect(item);
    },
    [handleThumbSelect],
  );

  useEffect(() => {
    if (!selectedSeries) return;
    const match = galleryItems.find(
      (g) =>
        g.series === selectedSeries &&
        (!activeAddress || g.original === activeAddress),
    );
    if (match) setSelection(getItemKey(match));
  }, [selectedSeries, activeAddress, galleryItems, setSelection]);

  useEffect(() => {
    if (!activeAddress || !selectedSeries) return;
    GetDalleDressCurrent(
      {
        collection: 'dalledress',
        dataFacet: types.DataFacet.GENERATOR,
        address: activeAddress,
      },
      selectedSeries,
    ).then((dd) => setCurrent(dd as unknown as dalle.DalleDress));
  }, [activeAddress, selectedSeries]);

  useEffect(() => {
    if (!orig && !series) return;
    let changed = false;
    if (orig && orig !== activeAddress) {
      setActiveAddress(orig);
      changed = true;
    }
    if (series && series !== selectedSeries) {
      setSelectedSeries(series);
      changed = true;
    }
    if (changed) clearSelection();
  }, [
    orig,
    series,
    activeAddress,
    selectedSeries,
    setActiveAddress,
    setSelectedSeries,
    clearSelection,
  ]);

  useEffect(() => {
    const selectedKey = getSelectionKey();
    if (!selectedKey || !thumbRowRef.current) return;
    const el = thumbRowRef.current.querySelector(`[data-key="${selectedKey}"]`);
    if (el && 'scrollIntoView' in el)
      (el as HTMLElement).scrollIntoView({
        inline: 'nearest',
        block: 'nearest',
      });
  }, [getSelectionKey]);

  const thumbItems = useMemo(() => {
    const effectiveOrig = orig || activeAddress || null;
    if (!effectiveOrig) return galleryItems;
    return galleryItems.filter((g) => g.original === effectiveOrig);
  }, [galleryItems, orig, activeAddress]);

  useEffect(() => {
    const selectedKey = getSelectionKey();
    if (!selectedKey) return;
    if (!thumbItems.find((g) => getItemKey(g) === selectedKey)) {
      const first = thumbItems[0];
      if (first) setSelection(getItemKey(first));
    }
  }, [thumbItems, getSelectionKey, setSelection]);

  const handleThumbKey = useCallback(
    (e: React.KeyboardEvent<HTMLDivElement>) => {
      if (!thumbItems.length) return;
      const selectedKey = getSelectionKey();
      const idx = selectedKey
        ? thumbItems.findIndex((g) => getItemKey(g) === selectedKey)
        : -1;
      let nextIdx = idx;
      if (e.key === 'ArrowRight') {
        nextIdx = (idx + 1 + thumbItems.length) % thumbItems.length;
        e.preventDefault();
      } else if (e.key === 'ArrowLeft') {
        nextIdx = (idx - 1 + thumbItems.length) % thumbItems.length;
        e.preventDefault();
      } else if (e.key === 'Home') {
        nextIdx = 0;
        e.preventDefault();
      } else if (e.key === 'End') {
        nextIdx = thumbItems.length - 1;
        e.preventDefault();
      } else {
        return;
      }
      const next = thumbItems[nextIdx];
      if (next) {
        setSelection(getItemKey(next));
        if (next.series) setSelectedSeries(next.series);
      }
    },
    [thumbItems, getSelectionKey, setSelection, setSelectedSeries],
  );

  const attributes = useMemo(() => current?.attributes || [], [current]);
  const selectedGalleryItem = useMemo(() => {
    const selectedKey = getSelectionKey();
    if (!selectedKey) return null;
    return galleryItems.find((g) => getItemKey(g) === selectedKey) || null;
  }, [getSelectionKey, galleryItems]);

  const displayImageUrl = useMemo(() => {
    if (current?.imageUrl) return current.imageUrl;
    if (selectedGalleryItem?.imageUrl) return selectedGalleryItem.imageUrl;
    const first = galleryItems.find((g) => g.original === activeAddress);
    if (first?.imageUrl) return first.imageUrl;
    if (galleryItems.length && galleryItems[0])
      return galleryItems[0].imageUrl || '';
    return '';
  }, [current?.imageUrl, selectedGalleryItem, galleryItems, activeAddress]);

  return (
    <Container size="xl" py="md">
      <Stack gap="sm">
        <Group align="flex-end" gap="sm" wrap="nowrap">
          <Select
            label="Address"
            placeholder="Select address"
            searchable
            value={activeAddress || null}
            data={addressOptions.map((a) => ({ value: a, label: a }))}
            onChange={handleAddressChange}
            w={380}
            size="xs"
          />
          <Select
            label="Series"
            placeholder="Series"
            value={selectedSeries}
            data={seriesOptions.map((s) => ({
              value: s,
              label: s === 'empty' ? '<empty>' : s,
            }))}
            onChange={handleSeriesChange}
            w={220}
            size="xs"
            disabled={!seriesOptions.length}
          />
          <Button
            size="xs"
            variant="default"
            onClick={handleGenerate}
            disabled={!activeAddress || !selectedSeries}
            style={{ alignSelf: 'flex-end' }}
          >
            Generate
          </Button>
        </Group>
        <div
          style={{ display: 'flex', alignItems: 'flex-start', width: '100%' }}
        >
          <div style={{ flex: '0 0 55%', maxWidth: '55%' }}>
            <Title order={6}>Image</Title>
            <div
              style={{
                border: '1px solid #444',
                background: '#222',
                marginTop: 4,
                width: '100%',
                position: 'relative',
                overflow: 'hidden',
              }}
            >
              {displayImageUrl ? (
                <Image
                  alt="Generated"
                  fit="contain"
                  radius="sm"
                  src={displayImageUrl}
                  style={{
                    display: 'block',
                    width: '100%',
                    height: 'auto',
                    objectFit: 'contain',
                  }}
                  onError={() => Log('generator:image:error')}
                />
              ) : (
                <Center h={160}>
                  <Text size="xs" c="dimmed">
                    No image
                  </Text>
                </Center>
              )}
            </div>
          </div>
          <div
            style={{
              flex: '0 0 35%',
              maxWidth: '35%',
              display: 'flex',
              flexDirection: 'column',
              gap: 12,
              paddingLeft: 12,
            }}
          >
            <div>
              <Title order={6}>Attributes</Title>
              <ScrollArea
                h={140}
                scrollbarSize={4}
                type="auto"
                offsetScrollbars
                style={{
                  marginTop: 4,
                  border: '1px solid #444',
                  background: '#222',
                  padding: 6,
                }}
              >
                <Stack
                  gap={2}
                  style={{ fontFamily: 'monospace', fontSize: 11 }}
                >
                  {attributes.length === 0 && (
                    <Text size="xs" c="dimmed">
                      No attributes
                    </Text>
                  )}
                  {attributes.map((a, i) => (
                    <Text key={i} size="xs">
                      {a.name}:{' '}
                      {a.value || a.selector || a.number || a.count || ''}
                    </Text>
                  ))}
                </Stack>
              </ScrollArea>
            </div>
            <div>
              <Title order={6}>Prompt</Title>
              <ScrollArea
                h={160}
                scrollbarSize={4}
                type="auto"
                offsetScrollbars
                style={{
                  marginTop: 4,
                  border: '1px solid #444',
                  background: '#222',
                  padding: 6,
                }}
              >
                <Text
                  size="xs"
                  style={{
                    fontFamily: 'monospace',
                    fontSize: 12,
                    whiteSpace: 'pre-wrap',
                  }}
                >
                  {current?.prompt || ''}
                </Text>
              </ScrollArea>
            </div>
            {!!current?.enhancedPrompt && (
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                  <Title order={6} style={{ flexGrow: 1 }}>
                    Enhanced Prompt
                  </Title>
                  <Button
                    variant="default"
                    size="xs"
                    loading={speaking}
                    onClick={speak}
                    leftSection={<SpeakIcon size={12} />}
                  >
                    Speak
                  </Button>
                </div>
                {audioUrl && (
                  <audio
                    ref={audioRef}
                    src={audioUrl}
                    controls
                    style={{ width: '100%', marginTop: 6 }}
                    onPlay={() => Log('readtome:play:start')}
                    onEnded={() => Log('readtome:play:end')}
                    onError={() => Log('readtome:play:error')}
                    autoPlay
                  />
                )}
                <ScrollArea
                  h={160}
                  scrollbarSize={4}
                  type="auto"
                  offsetScrollbars
                  style={{
                    marginTop: 4,
                    border: '1px solid #444',
                    background: '#222',
                    padding: 6,
                  }}
                >
                  <Text
                    size="xs"
                    style={{
                      fontFamily: 'monospace',
                      fontSize: 12,
                      whiteSpace: 'pre-wrap',
                    }}
                  >
                    {current?.enhancedPrompt || ''}
                  </Text>
                </ScrollArea>
              </div>
            )}
            <div>
              <Title order={6}>Thumbnails</Title>
              <div
                style={{
                  display: 'flex',
                  gap: 4,
                  overflowX: 'auto',
                  padding: '4px 2px',
                  marginTop: 4,
                }}
                ref={thumbRowRef}
                tabIndex={0}
                onKeyDown={handleThumbKey}
              >
                {thumbItems.map((g) => {
                  const itemKey = getItemKey(g);
                  return (
                    <div
                      key={itemKey}
                      data-key={itemKey}
                      style={{ width: 72, flex: '0 0 auto' }}
                    >
                      <DalleDressCard
                        item={g}
                        onClick={handleThumbSelect}
                        onDoubleClick={handleThumbDouble}
                        selected={itemKey === getSelectionKey()}
                      />
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
          <div
            style={{
              flex: '0 0 10%',
              maxWidth: '10%',
              display: 'flex',
              flexDirection: 'column',
              paddingLeft: 12,
            }}
          >
            <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
              {['Claim', 'Mint', 'Burn', 'Trade', 'Eject', 'Merch'].map(
                (label) => (
                  <Button
                    key={label}
                    variant="default"
                    size="xs"
                    fullWidth
                    onClick={() =>
                      Log('generator:button:' + label.toLowerCase())
                    }
                  >
                    {label}
                  </Button>
                ),
              )}
            </div>
          </div>
        </div>
      </Stack>
    </Container>
  );
}
