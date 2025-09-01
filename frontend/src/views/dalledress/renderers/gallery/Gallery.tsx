import { useCallback, useEffect, useRef, useState } from 'react';

import type { DataFacet } from '@hooks';
import { Center, Container, Title } from '@mantine/core';
import { dalle, dalledress, project, types } from '@models';

import { GalleryControls, GalleryGrouping } from '../../components';
import { getItemKey, useGalleryStore } from '../../store';

export type GalleryProps = {
  pageData: dalledress.DalleDressPage | null;
  viewStateKey?: project.ViewStateKey;
  setActiveFacet?: (f: DataFacet) => void;
};

export function Gallery({
  pageData,
  viewStateKey,
  setActiveFacet,
}: GalleryProps) {
  const [controls, setControls] = useState<{
    sortMode: 'series' | 'address';
    columns: number;
  }>({ sortMode: 'series', columns: 6 });
  const scrollRef = useRef<HTMLDivElement | null>(null);
  const keyScopeRef = useRef<HTMLDivElement | null>(null);

  const {
    getSelectionKey,
    setSelection,
    ingestItems,
    galleryItems,
    useDerived,
  } = useGalleryStore();
  const { groupNames, groupedItems, flattenedItems } = useDerived(
    controls.sortMode,
  );

  useEffect(() => {
    ingestItems(pageData?.dalledress || []);
  }, [pageData?.dalledress, ingestItems]);

  useEffect(() => {
    if (!getSelectionKey() && flattenedItems.length > 0) {
      const firstItem = flattenedItems[0];
      if (firstItem) {
        setSelection(getItemKey(firstItem));
      }
    }
  }, [flattenedItems, getSelectionKey, setSelection]);

  // --------------------------------------
  const scrollSelectedIntoView = useCallback((selected: string | null) => {
    if (!selected || !scrollRef.current) {
      return;
    }
    const el = scrollRef.current.querySelector(`[data-key="${selected}"]`);
    if (el && 'scrollIntoView' in el) {
      (el as HTMLElement).scrollIntoView({ block: 'nearest' });
    }
  }, []);

  useEffect(() => {
    scrollSelectedIntoView(getSelectionKey());
  }, [getSelectionKey, scrollSelectedIntoView, flattenedItems]);

  // --------------------------------------
  const handleItemClick = useCallback(
    (item: dalle.DalleDress) => {
      setSelection(getItemKey(item));
    },
    [setSelection],
  );

  const handleItemDoubleClick = useCallback(
    (item: dalle.DalleDress) => {
      setSelection(getItemKey(item));
      if (setActiveFacet)
        setActiveFacet(types.DataFacet.GENERATOR as DataFacet);
    },
    [setActiveFacet, setSelection],
  );

  const handleKey = useCallback(
    (e: React.KeyboardEvent<HTMLDivElement>) => {
      if (!flattenedItems.length) return;
      const selectedKey = getSelectionKey();
      const idx = flattenedItems.findIndex(
        (f) => getItemKey(f) === selectedKey,
      );
      if (e.key === 'ArrowDown' || e.key === 'ArrowRight') {
        const next =
          flattenedItems[
            (idx + 1 + flattenedItems.length) % flattenedItems.length
          ] || null;
        if (next) setSelection(getItemKey(next));
        e.preventDefault();
      } else if (e.key === 'ArrowUp' || e.key === 'ArrowLeft') {
        const next =
          flattenedItems[
            (idx - 1 + flattenedItems.length) % flattenedItems.length
          ] || null;
        if (next) setSelection(getItemKey(next));
        e.preventDefault();
      } else if (e.key === 'Enter') {
        if (idx >= 0 && flattenedItems[idx])
          handleItemDoubleClick(flattenedItems[idx]);
      }
    },
    [flattenedItems, getSelectionKey, handleItemDoubleClick, setSelection],
  );

  // --------------------------------------
  return (
    <Container
      size="xl"
      py="md"
      ref={keyScopeRef}
      tabIndex={0}
      onKeyDown={handleKey}
      onMouseDown={() => keyScopeRef.current && keyScopeRef.current.focus()}
      style={{ outline: 'none', width: '100%' }}
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
        {groupNames.map((series) => (
          <GalleryGrouping
            key={series || 'unknown'}
            series={series}
            items={groupedItems[series] || []}
            columns={controls.columns}
            onItemClick={handleItemClick}
            onItemDoubleClick={handleItemDoubleClick}
            selected={getSelectionKey()}
          />
        ))}
      </div>
    </Container>
  );
}
