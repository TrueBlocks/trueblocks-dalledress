import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import type { DataFacet } from '@hooks';
import { Center, Container, Title } from '@mantine/core';
import { dalle, dalledress, project, types } from '@models';

import { GalleryControls, SeriesGallery } from '../../components';
import { useDalleDressSelection } from '../../store';

export function Gallery({
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

  const galleryItems: dalle.DalleDress[] = useMemo(
    () => (pageData?.gallery ? [...pageData.gallery] : []),
    [pageData?.gallery],
  );

  const grouped = useMemo(() => {
    const item: Record<string, dalle.DalleDress[]> = {};
    galleryItems.forEach((it) => {
      const key =
        controls.sortMode === 'series' ? it.series || '' : it.original || '';
      if (!item[key]) item[key] = [];
      item[key].push(it);
    });
    return item;
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

  const handleItemClick = useCallback((item: dalle.DalleDress) => {
    setSelected(item.annotatedPath);
  }, []);

  const { setDressSelection } = useDalleDressSelection();
  const handleItemDoubleClick = useCallback(
    (item: dalle.DalleDress) => {
      setSelected(item.annotatedPath);
      setDressSelection(item.original, item.series, item.annotatedPath);
      if (setActiveFacet)
        setActiveFacet(types.DataFacet.GENERATOR as DataFacet);
    },
    [setActiveFacet, setDressSelection],
  );

  const scrollSelectedIntoView = useCallback((annotatedPath: string | null) => {
    if (!annotatedPath || !scrollRef.current) return;
    const el = scrollRef.current.querySelector(
      `[data-relpath="${annotatedPath}"]`,
    );
    if (el && 'scrollIntoView' in el)
      (el as HTMLElement).scrollIntoView({ block: 'nearest' });
  }, []);

  useEffect(() => {
    scrollSelectedIntoView(selected);
  }, [selected, scrollSelectedIntoView]);

  const handleKey = useCallback(
    (e: React.KeyboardEvent<HTMLDivElement>) => {
      if (!flattened.length) return;
      const idx = flattened.findIndex((f) => f.annotatedPath === selected);
      if (e.key === 'ArrowDown' || e.key === 'ArrowRight') {
        const next =
          flattened[(idx + 1 + flattened.length) % flattened.length] || null;
        if (next) setSelected(next.annotatedPath);
        e.preventDefault();
      } else if (e.key === 'ArrowUp' || e.key === 'ArrowLeft') {
        const next =
          flattened[(idx - 1 + flattened.length) % flattened.length] || null;
        if (next) setSelected(next.annotatedPath);
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
