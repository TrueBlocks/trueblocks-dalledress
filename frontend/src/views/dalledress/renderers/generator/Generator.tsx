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
import { useGalleryStore } from '../../store';

export function Generator({
  pageData,
}: {
  pageData: dalledress.DalleDressPage | null;
}) {
  const { activeAddress, setActiveAddress } = useActiveProject();
  const [selectedSeries, setSelectedSeries] = useState<string | null>('empty');
  const [selected, setSelected] = useState<string | null>(null);
  const [current, setCurrent] = useState<dalle.DalleDress | null>(
    pageData?.currentDress || null,
  );
  const { speaking, audioUrl, audioRef, speak } = useSpeakPrompt({
    activeAddress: activeAddress || null,
    selectedSeries,
    hasEnhancedPrompt: !!current?.enhancedPrompt,
  });
  const thumbRowRef = useRef<HTMLDivElement | null>(null);
  const icons = useIconSets();
  const SpeakIcon = icons.Speak;

  useEffect(() => {
    setCurrent(pageData?.currentDress || null);
  }, [pageData?.currentDress]);

  const {
    orig,
    series,
    getSelection,
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
      if (!selected && first) setSelected(first.annotatedPath);
    }
  }, [selectedSeries, selected, galleryItems, activeAddress]);

  const handleAddressChange = useCallback(
    (value: string | null) => {
      if (!value) return;
      Log('generator:address:' + value);
      setActiveAddress(value);
    },
    [setActiveAddress],
  );

  const handleSeriesChange = useCallback((value: string | null) => {
    setSelectedSeries(value);
    if (value) Log('generator:series:' + value);
  }, []);

  const handleGenerate = useCallback(() => {
    if (!activeAddress || !selectedSeries) return;
    Log('generator:generate');
  }, [activeAddress, selectedSeries]);

  const handleThumbSelect = useCallback((item: dalle.DalleDress) => {
    setSelected(item.annotatedPath);
    if (item.series) setSelectedSeries(item.series);
    Log('generator:thumb:select:' + item.annotatedPath);
  }, []);

  const handleThumbDouble = useCallback(
    (item: dalle.DalleDress) => {
      Log('generator:thumb:dbl:' + item.annotatedPath);
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
    if (match) setSelected(match.annotatedPath);
  }, [selectedSeries, activeAddress, galleryItems]);

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
      Log('generator:pref:address:' + orig);
    }
    if (series && series !== selectedSeries) {
      setSelectedSeries(series);
      changed = true;
      Log('generator:pref:series:' + series);
    }
    setSelected(getSelection());
    if (changed) clearSelection();
  }, [
    orig,
    series,
    activeAddress,
    selectedSeries,
    getSelection,
    setActiveAddress,
    setSelectedSeries,
    clearSelection,
  ]);

  useEffect(() => {
    if (!selected || !thumbRowRef.current) return;
    const el = thumbRowRef.current.querySelector(
      `[data-relpath="${selected}"]`,
    );
    if (el && 'scrollIntoView' in el)
      (el as HTMLElement).scrollIntoView({
        inline: 'nearest',
        block: 'nearest',
      });
  }, [selected]);

  const thumbItems = useMemo(() => {
    const effectiveOrig = orig || null;
    if (!effectiveOrig) return galleryItems;
    return galleryItems.filter((g) => g.original === effectiveOrig);
  }, [galleryItems, orig]);

  useEffect(() => {
    if (!selected) return;
    if (!thumbItems.find((g) => g.annotatedPath === selected)) {
      const first = thumbItems[0];
      if (first) setSelected(first.annotatedPath);
    }
  }, [thumbItems, selected]);

  const handleThumbKey = useCallback(
    (e: React.KeyboardEvent<HTMLDivElement>) => {
      if (!thumbItems.length) return;
      const idx = selected
        ? thumbItems.findIndex((g) => g.annotatedPath === selected)
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
        setSelected(next.annotatedPath);
        if (next.series) setSelectedSeries(next.series);
      }
    },
    [thumbItems, selected, setSelectedSeries],
  );

  const attributes = useMemo(() => current?.attributes || [], [current]);
  const selectedGalleryItem = useMemo(() => {
    if (!selected) return null;
    return galleryItems.find((g) => g.annotatedPath === selected) || null;
  }, [selected, galleryItems]);

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
                {thumbItems.map((g) => (
                  <div
                    key={g.annotatedPath}
                    data-relpath={g.annotatedPath}
                    style={{ width: 72, flex: '0 0 auto' }}
                  >
                    <DalleDressCard
                      item={g}
                      onClick={handleThumbSelect}
                      onDoubleClick={handleThumbDouble}
                      selected={g.annotatedPath === selected}
                    />
                  </div>
                ))}
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
