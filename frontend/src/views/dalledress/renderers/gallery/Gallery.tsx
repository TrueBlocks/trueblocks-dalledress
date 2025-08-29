import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import type { DataFacet } from '@hooks';
import { Center, Container, Title } from '@mantine/core';
import { dalledress, project, types } from '@models';
import {
  GalleryControls,
  SeriesGallery,
} from 'src/views/dalledress/components';
import { setPendingGeneratorSelection } from 'src/views/dalledress/renderers/generator';

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

  const scrollSelectedIntoView = useCallback((relPath: string | null) => {
    if (!relPath || !scrollRef.current) return;
    const el = scrollRef.current.querySelector(`[data-relpath="${relPath}"]`);
    if (el && 'scrollIntoView' in el)
      (el as HTMLElement).scrollIntoView({ block: 'nearest' });
  }, []);

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
