import {
  ReactElement,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';

import { GetDalleDressCurrent } from '@app';
import { useActiveProject } from '@hooks';
import type { DataFacet } from '@hooks';
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
import { dalle, dalledress, project, types } from '@models';
import { Log } from '@utils';

import { DalleDressCard } from './DalleDressCard';
import { GalleryControls } from './GalleryControls';
import { SeriesGallery } from './SeriesGallery';

let pendingGeneratorSelection: {
  address?: string;
  series?: string | null;
} | null = null;
export const setPendingGeneratorSelection = (
  address: string,
  series: string | null,
) => {
  pendingGeneratorSelection = { address, series };
};

function GeneratorRenderer({
  pageData,
}: {
  pageData: dalledress.DalleDressPage | null;
}) {
  const { activeAddress, setActiveAddress } = useActiveProject();
  const [selectedSeries, setSelectedSeries] = useState<string | null>('empty');
  const [selectedThumb, setSelectedThumb] = useState<string | null>(null);
  const [current, setCurrent] = useState<dalle.DalleDress | null>(
    pageData?.currentDress || null,
  );

  useEffect(() => {
    setCurrent(pageData?.currentDress || null);
  }, [pageData?.currentDress]);

  const galleryItems = useMemo(
    () => (pageData?.gallery ? [...pageData.gallery] : []),
    [pageData?.gallery],
  );

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
    galleryItems.forEach((g) => set.add(g.address));
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
        (g) => g.address === activeAddress && g.series,
      );
      if (first?.series) setSelectedSeries(first.series);
      if (!selectedThumb && first) setSelectedThumb(first.relPath);
    }
  }, [selectedSeries, selectedThumb, galleryItems, activeAddress]);

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

  const handleThumbSelect = useCallback((item: dalledress.GalleryItem) => {
    setSelectedThumb(item.relPath);
    if (item.series) setSelectedSeries(item.series);
    Log('generator:thumb:select:' + item.relPath);
  }, []);

  const handleThumbDouble = useCallback(
    (item: dalledress.GalleryItem) => {
      Log('generator:thumb:dbl:' + item.relPath);
      handleThumbSelect(item);
    },
    [handleThumbSelect],
  );

  useEffect(() => {
    if (!selectedSeries) return;
    const match = galleryItems.find(
      (g) =>
        g.series === selectedSeries &&
        (!activeAddress || g.address === activeAddress),
    );
    if (match) setSelectedThumb(match.relPath);
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
    if (!pendingGeneratorSelection) return;
    const { address, series } = pendingGeneratorSelection;
    let changed = false;
    if (address && address !== activeAddress) {
      setActiveAddress(address);
      changed = true;
      Log('generator:pref:address:' + address);
    }
    if (series && series !== selectedSeries) {
      setSelectedSeries(series);
      changed = true;
      Log('generator:pref:series:' + series);
    }
    if (changed) pendingGeneratorSelection = null;
  }, [activeAddress, selectedSeries, setActiveAddress]);

  const attributes = useMemo(() => current?.attributes || [], [current]);
  const selectedGalleryItem = useMemo(() => {
    if (!selectedThumb) return null;
    return galleryItems.find((g) => g.relPath === selectedThumb) || null;
  }, [selectedThumb, galleryItems]);

  const displayImageUrl = useMemo(() => {
    if (current?.imageUrl) return current.imageUrl;
    if (selectedGalleryItem?.url) return selectedGalleryItem.url;
    const first = galleryItems.find((g) => g.address === activeAddress);
    if (first?.url) return first.url;
    if (galleryItems.length && galleryItems[0])
      return galleryItems[0].url || '';
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
                <Title order={6}>Enhanced Prompt</Title>
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
              >
                {galleryItems.map((g) => (
                  <div key={g.relPath} style={{ width: 72, flex: '0 0 auto' }}>
                    <DalleDressCard
                      item={g}
                      onClick={handleThumbSelect}
                      onDoubleClick={handleThumbDouble}
                      selected={g.relPath === selectedThumb}
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
              {[
                'Connect',
                'Claim',
                'Mint',
                'Burn',
                'Trade',
                'Eject',
                'Merch',
              ].map((label) => (
                <Button
                  key={label}
                  variant="default"
                  size="xs"
                  fullWidth
                  onClick={() => Log('generator:button:' + label.toLowerCase())}
                >
                  {label}
                </Button>
              ))}
            </div>
          </div>
        </div>
      </Stack>
    </Container>
  );
}
function GalleryRenderer({
  pageData,
  viewStateKey,
  setActiveFacet,
}: {
  pageData: dalledress.DalleDressPage | null;
  viewStateKey?: project.ViewStateKey;
  setActiveFacet?: (f: DataFacet) => void;
}) {
  const keyScopeRef = useRef<HTMLDivElement | null>(null);
  const scrollRef = useRef<HTMLDivElement | null>(null);
  const [controls, setControls] = useState<{
    sortMode: 'series' | 'address';
    columns: number;
  }>({ sortMode: 'series', columns: 6 });

  const galleryItems = useMemo(
    () => (pageData?.gallery ? [...pageData.gallery] : []),
    [pageData?.gallery],
  );

  const grouped = useMemo(() => {
    const out: Record<string, dalledress.GalleryItem[]> = {};
    galleryItems.forEach((it) => {
      const key =
        controls.sortMode === 'series' ? it.series || '' : it.address || '';
      if (!out[key]) out[key] = [];
      out[key].push(it);
    });
    return out;
  }, [galleryItems, controls.sortMode]);

  const seriesNames = useMemo(
    () => Object.keys(grouped).sort((a, b) => a.localeCompare(b)),
    [grouped],
  );

  const flattened = useMemo(
    () => seriesNames.flatMap((s) => grouped[s] || []),
    [seriesNames, grouped],
  );
  const [selected, setSelected] = useState<string | null>(null);

  const scrollSelectedIntoView = useCallback((relPath: string | null) => {
    if (!relPath || !scrollRef.current) return;
    const el = scrollRef.current.querySelector(`[data-relpath="${relPath}"]`);
    if (el && 'scrollIntoView' in el)
      (el as HTMLElement).scrollIntoView({ block: 'nearest' });
  }, []);

  const handleItemClick = useCallback((item: dalledress.GalleryItem) => {
    setSelected(item.relPath);
  }, []);

  const handleItemDoubleClick = useCallback(
    (item: dalledress.GalleryItem) => {
      setSelected(item.relPath);
      setPendingGeneratorSelection(item.address, item.series);
      if (setActiveFacet)
        setActiveFacet(types.DataFacet.GENERATOR as DataFacet);
    },
    [setActiveFacet],
  );

  useEffect(() => {
    scrollSelectedIntoView(selected);
  }, [selected, scrollSelectedIntoView]);

  const handleKey = useCallback(
    (e: React.KeyboardEvent<HTMLDivElement>) => {
      if (!flattened.length) return;
      const idx = flattened.findIndex((f) => f.relPath === selected);
      if (e.key === 'ArrowDown' || e.key === 'ArrowRight') {
        const next =
          flattened[(idx + 1 + flattened.length) % flattened.length] || null;
        if (next) setSelected(next.relPath);
        e.preventDefault();
      } else if (e.key === 'ArrowUp' || e.key === 'ArrowLeft') {
        const next =
          flattened[(idx - 1 + flattened.length) % flattened.length] || null;
        if (next) setSelected(next.relPath);
        e.preventDefault();
      } else if (e.key === 'Enter') {
        if (idx >= 0 && flattened[idx]) handleItemDoubleClick(flattened[idx]);
      }
    },
    [flattened, selected, handleItemDoubleClick],
  );

  return (
    <Container
      size="xl"
      py="md"
      ref={keyScopeRef}
      tabIndex={0}
      onKeyDown={handleKey}
      onMouseDown={() => keyScopeRef.current && keyScopeRef.current.focus()}
      style={{ outline: 'none' }}
    >
      <Title order={4} mb="sm">
        Preview Gallery
      </Title>
      {viewStateKey && (
        <GalleryControls
          viewStateKey={viewStateKey}
          value={{ sortMode: controls.sortMode, columns: controls.columns }}
          onChange={(v) => setControls(v)}
        />
      )}
      {!galleryItems.length && (
        <Center
          style={{
            opacity: 0.6,
            fontSize: 12,
            fontFamily: 'monospace',
          }}
        >
          No images found
        </Center>
      )}
      <div
        ref={scrollRef}
        style={{
          maxHeight: 'calc(100vh - 260px)',
          overflowY: 'auto',
          paddingRight: 4,
        }}
      >
        {seriesNames.map((series) => (
          <SeriesGallery
            key={series || 'unknown'}
            series={series}
            items={grouped[series] || []}
            columns={controls.columns}
            onItemClick={handleItemClick}
            onItemDoubleClick={handleItemDoubleClick}
            selectedRelPath={selected}
          />
        ))}
      </div>
    </Container>
  );
}

export function renderers(
  pageData: dalledress.DalleDressPage | null,
  viewStateKey?: project.ViewStateKey,
  setActiveFacet?: (f: DataFacet) => void,
) {
  return {
    [types.DataFacet.GENERATOR]: () => (
      <GeneratorRenderer pageData={pageData} />
    ),
    [types.DataFacet.GALLERY]: () => (
      <GalleryRenderer
        pageData={pageData}
        viewStateKey={viewStateKey}
        setActiveFacet={setActiveFacet}
      />
    ),
  } as Record<types.DataFacet, () => ReactElement>;
}
