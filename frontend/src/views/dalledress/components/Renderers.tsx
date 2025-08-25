import {
  ReactElement,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';

import { useActiveProject } from '@hooks';
import {
  Button,
  Center,
  Container,
  Grid,
  Group,
  Image,
  ScrollArea,
  Select,
  Stack,
  Text,
  Textarea,
  Title,
} from '@mantine/core';
import { dalledress, project, types } from '@models';
import { Log } from '@utils';

import { GalleryControls } from './GalleryControls';
import { GeneratorThumb } from './GeneratorThumb';
import { SeriesGallery } from './SeriesGallery';

function GeneratorRenderer({
  pageData,
}: {
  pageData: dalledress.DalleDressPage | null;
}) {
  const { activeAddress, setActiveAddress } = useActiveProject();
  const [selectedSeries, setSelectedSeries] = useState<string | null>(null);
  const [selectedThumb, setSelectedThumb] = useState<string | null>(null);

  const current = pageData?.currentDress || null;
  const galleryItems = useMemo(
    () => (pageData?.gallery ? [...pageData.gallery] : []),
    [pageData?.gallery],
  );

  const seriesOptions = useMemo(() => {
    const set = new Set<string>();
    galleryItems.forEach((g) => {
      if (g.series) set.add(g.series);
    });
    return Array.from(set).sort((a, b) => a.localeCompare(b));
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
      const firstForAddress = galleryItems.find(
        (g) => g.address === activeAddress && g.series,
      );
      if (firstForAddress?.series) setSelectedSeries(firstForAddress.series);
      if (!selectedThumb && firstForAddress)
        setSelectedThumb(firstForAddress.relPath);
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

  const attributes = useMemo(() => current?.attributes || [], [current]);

  const selectedGalleryItem = useMemo(() => {
    if (!selectedThumb) return null;
    return galleryItems.find((g) => g.relPath === selectedThumb) || null;
  }, [selectedThumb, galleryItems]);

  const displayImageUrl = useMemo(() => {
    if (current?.imageUrl) return current.imageUrl;
    if (selectedGalleryItem?.url) return selectedGalleryItem.url;
    const firstForAddress = galleryItems.find(
      (g) => g.address === activeAddress,
    );
    if (firstForAddress?.url) return firstForAddress.url;
    if (galleryItems?.length) return galleryItems[0]?.url || '';
    return '';
  }, [current?.imageUrl, selectedGalleryItem, galleryItems, activeAddress]);

  return (
    <Container size="xl" py="md">
      <Grid gutter="md">
        <Grid.Col span={10}>
          <Stack gap="xs">
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
                data={seriesOptions.map((s) => ({ value: s, label: s }))}
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
            <Stack gap="md">
              <div>
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
              <div>
                <Title order={6}>Attributes</Title>
                <ScrollArea
                  h={120}
                  scrollbarSize={4}
                  type="auto"
                  offsetScrollbars
                  style={{ marginTop: 4 }}
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
                <Textarea
                  value={current?.prompt || ''}
                  placeholder="Prompt"
                  autosize
                  minRows={4}
                  styles={{
                    input: { fontFamily: 'monospace', fontSize: 12 },
                  }}
                  readOnly
                  style={{ marginTop: 4 }}
                />
              </div>
              {!!current?.enhancedPrompt && (
                <div>
                  <Title order={6}>Enhanced Prompt</Title>
                  <Textarea
                    value={current?.enhancedPrompt || ''}
                    placeholder="Enhanced Prompt"
                    autosize
                    minRows={4}
                    styles={{
                      input: { fontFamily: 'monospace', fontSize: 12 },
                    }}
                    readOnly
                    style={{ marginTop: 4 }}
                  />
                </div>
              )}
              <div>
                <Title order={6}>Thumbnails</Title>
                <Group
                  gap={4}
                  wrap="nowrap"
                  style={{ overflowX: 'auto', marginTop: 4 }}
                >
                  {galleryItems.map((g) => (
                    <GeneratorThumb
                      key={g.relPath}
                      item={g}
                      onSelect={handleThumbSelect}
                      selected={g.relPath === selectedThumb}
                    />
                  ))}
                </Group>
              </div>
            </Stack>
          </Stack>
        </Grid.Col>
        <Grid.Col span={2}>
          <Stack gap="xs" style={{ width: 140, marginLeft: 'auto' }}>
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
          </Stack>
        </Grid.Col>
      </Grid>
    </Container>
  );
}

interface GalleryRendererProps {
  pageData: dalledress.DalleDressPage | null;
  viewStateKey?: project.ViewStateKey;
}

function GalleryRenderer({ pageData, viewStateKey }: GalleryRendererProps) {
  const [controls, setControls] = useState({
    sortMode: 'series' as 'series' | 'address',
    columns: 6,
  });
  const [selected, setSelected] = useState<string | null>(null);
  const keyScopeRef = useRef<HTMLDivElement | null>(null);
  const scrollRef = useRef<HTMLDivElement | null>(null);
  const preferredColRef = useRef<number | null>(null);

  const galleryItems = useMemo(
    () => (pageData?.gallery ? [...pageData.gallery] : []),
    [pageData?.gallery],
  );

  const { grouped, seriesNames } = useMemo(() => {
    const g = galleryItems.reduce<Record<string, dalledress.GalleryItem[]>>(
      (acc, item) => {
        const key = item.series || '';
        if (!acc[key]) acc[key] = [];
        acc[key].push(item);
        return acc;
      },
      {},
    );
    const names = Object.keys(g).sort((a, b) => a.localeCompare(b));
    names.forEach((s) => {
      const arr = g[s];
      if (!arr) return;
      if (controls.sortMode === 'address') {
        arr.sort((a, b) => {
          const ad = a.address.localeCompare(b.address);
          if (ad !== 0) return ad;
          return a.fileName.localeCompare(b.fileName);
        });
      } else {
        arr.sort((a, b) => {
          if (a.index >= 0 && b.index >= 0 && a.index !== b.index)
            return a.index - b.index;
          return a.fileName.localeCompare(b.fileName);
        });
      }
    });
    return { grouped: g, seriesNames: names };
  }, [galleryItems, controls.sortMode]);

  const rowMatrix = useMemo(() => {
    const cols = controls.columns || 1;
    const rowsAcc: dalledress.GalleryItem[][] = [];
    seriesNames.forEach((s) => {
      const its = grouped[s] || [];
      for (let i = 0; i < its.length; i += cols) {
        rowsAcc.push(its.slice(i, i + cols));
      }
    });
    return rowsAcc;
  }, [grouped, seriesNames, controls.columns]);

  const positionMap = useMemo(() => {
    const map = new Map<string, { row: number; col: number }>();
    rowMatrix.forEach((row, r) => {
      row.forEach((item, c) => map.set(item.relPath, { row: r, col: c }));
    });
    return map;
  }, [rowMatrix]);

  useEffect(() => {
    if (!rowMatrix.length) {
      setSelected(null);
      preferredColRef.current = null;
      return;
    }
    if (!selected) {
      const firstRow = rowMatrix[0];
      const first = firstRow && firstRow[0];
      if (first) {
        setSelected(first.relPath);
        if (preferredColRef.current == null) preferredColRef.current = 0;
      }
    }
  }, [rowMatrix, selected]);

  const moveHorizontal = useCallback(
    (dir: 1 | -1) => {
      if (!selected) {
        if (rowMatrix[0] && rowMatrix[0][0])
          setSelected(rowMatrix[0][0].relPath);
        return;
      }
      const pos = positionMap.get(selected);
      if (!pos) return;
      let { row, col } = pos;
      col += dir;
      let rowArr = rowMatrix[row] || [];
      if (col < 0) {
        row = (row - 1 + rowMatrix.length) % rowMatrix.length;
        rowArr = rowMatrix[row] || [];
        col = rowArr.length - 1;
      } else if (col >= rowArr.length) {
        row = (row + 1) % rowMatrix.length;
        rowArr = rowMatrix[row] || [];
        col = 0;
      }
      const target = rowArr[col];
      if (target) {
        preferredColRef.current = col;
        setSelected(target.relPath);
      }
    },
    [selected, rowMatrix, positionMap],
  );

  const moveVertical = useCallback(
    (dir: 1 | -1) => {
      if (!rowMatrix.length) return;
      if (!selected) {
        if (rowMatrix[0] && rowMatrix[0][0]) {
          setSelected(rowMatrix[0][0].relPath);
          if (preferredColRef.current == null) preferredColRef.current = 0;
        }
        return;
      }
      const pos = positionMap.get(selected);
      if (!pos) return;
      const { row, col } = pos;
      if (preferredColRef.current == null) preferredColRef.current = col;
      let targetRow = row + dir;
      targetRow =
        ((targetRow % rowMatrix.length) + rowMatrix.length) % rowMatrix.length;
      const targetRowItems = rowMatrix[targetRow];
      if (!targetRowItems || !targetRowItems.length) return;
      const desiredCol = preferredColRef.current;
      const actualCol = Math.min(desiredCol, targetRowItems.length - 1);
      const targetItem = targetRowItems[actualCol];
      if (targetItem) setSelected(targetItem.relPath);
    },
    [selected, rowMatrix, positionMap],
  );

  const handleKey = useCallback(
    (e: KeyboardEvent) => {
      if (
        ![
          'ArrowLeft',
          'ArrowRight',
          'ArrowUp',
          'ArrowDown',
          'Enter',
          'Home',
          'End',
        ].includes(e.key)
      )
        return;
      if (
        e.key.startsWith('Arrow') ||
        e.key === 'Enter' ||
        e.key === 'Home' ||
        e.key === 'End'
      ) {
        e.preventDefault();
      }
      switch (e.key) {
        case 'ArrowLeft':
          moveHorizontal(-1);
          break;
        case 'ArrowRight':
          moveHorizontal(1);
          break;
        case 'ArrowUp':
          moveVertical(-1);
          break;
        case 'ArrowDown':
          moveVertical(1);
          break;
        case 'Enter': {
          if (selected) {
            const pos = positionMap.get(selected);
            if (pos) Log('Enter pressed on ' + selected);
          }
          break;
        }
        case 'Home': {
          if (rowMatrix[0] && rowMatrix[0][0]) {
            preferredColRef.current = 0;
            setSelected(rowMatrix[0][0].relPath);
          }
          break;
        }
        case 'End': {
          const lastRow = rowMatrix[rowMatrix.length - 1];
          if (lastRow && lastRow.length) {
            const last = lastRow[lastRow.length - 1];
            if (last) {
              preferredColRef.current = lastRow.length - 1;
              setSelected(last.relPath);
            }
          }
          break;
        }
      }
    },
    [moveHorizontal, moveVertical, selected, rowMatrix, positionMap],
  );

  useEffect(() => {
    const el = keyScopeRef.current;
    if (!el) return;
    el.addEventListener('keydown', handleKey);
    return () => {
      el.removeEventListener('keydown', handleKey);
    };
  }, [handleKey]);

  useEffect(() => {
    if (keyScopeRef.current) keyScopeRef.current.focus();
  }, []);

  const handleItemClick = useCallback(
    (item: dalledress.GalleryItem) => {
      Log('gallery:card:' + item.relPath);
      setSelected(item.relPath);
      const pos = positionMap.get(item.relPath);
      if (pos) preferredColRef.current = pos.col;
    },
    [positionMap],
  );

  const hasAutoFocusedRef = useRef(false);
  useEffect(() => {
    if (!selected) return;
    if (hasAutoFocusedRef.current) return;
    const el = keyScopeRef.current;
    if (!el) return;
    requestAnimationFrame(() => {
      if (document.activeElement !== el) el.focus({ preventScroll: true });
      hasAutoFocusedRef.current = true;
    });
  }, [selected]);

  useEffect(() => {
    if (!selected) return;
    const container = scrollRef.current || keyScopeRef.current;
    if (!container) return;
    const cssEsc: (s: string) => string = ((): ((s: string) => string) => {
      const w = window as unknown as {
        CSS?: { escape?: (s: string) => string };
      };
      if (w.CSS && typeof w.CSS.escape === 'function') return w.CSS.escape;
      return (s: string) => s.replace(/"/g, '\\"');
    })();
    const run = () => {
      const el = container.querySelector(
        `[data-relpath="${cssEsc(selected)}"]`,
      ) as HTMLElement | null;
      if (!el) return;
      const cRect = container.getBoundingClientRect();
      const r = el.getBoundingClientRect();
      const outY = r.top < cRect.top || r.bottom > cRect.bottom;
      const outX = r.left < cRect.left || r.right > cRect.right;
      if (outY || outX) {
        const topDelta = r.top - cRect.top;
        const desiredScrollTop = container.scrollTop + topDelta - 12;
        container.scrollTo({
          top: desiredScrollTop,
          behavior: 'instant' as ScrollBehavior,
        });
        el.scrollIntoView({ block: 'nearest', inline: 'nearest' });
      }
    };
    requestAnimationFrame(() => run());
  }, [selected]);

  return (
    <Container
      size="xl"
      py="md"
      ref={keyScopeRef}
      tabIndex={0}
      onMouseDown={() => keyScopeRef.current && keyScopeRef.current.focus()}
      style={{ outline: 'none' }}
    >
      <Title order={4} mb="sm">
        Preview Gallery
      </Title>
      {viewStateKey && (
        <GalleryControls
          viewStateKey={viewStateKey}
          value={controls}
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
) {
  return {
    [types.DataFacet.GENERATOR]: () => (
      <GeneratorRenderer pageData={pageData} />
    ),
    [types.DataFacet.GALLERY]: () => (
      <GalleryRenderer pageData={pageData} viewStateKey={viewStateKey} />
    ),
  } as Record<types.DataFacet, () => ReactElement>;
}
