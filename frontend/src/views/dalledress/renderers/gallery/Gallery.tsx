import { useCallback, useEffect, useRef, useState } from 'react';

import type { DataFacet } from '@hooks';
import { Center, Container, Title } from '@mantine/core';
import { dalle, dalledress, project, types } from '@models';

import { GalleryControls, GalleryGrouping } from '../../components';
import { useScrollSelectedIntoView } from '../../hooks/useScrollSelectedIntoView';
import { getItemKey, useGalleryStore } from '../../store';

export type GalleryProps = {
  pageData: dalledress.DalleDressPage | null;
  viewStateKey?: project.ViewStateKey;
  setActiveFacet?: (f: DataFacet) => void;
};

export const Gallery = ({
  pageData,
  viewStateKey,
  setActiveFacet,
}: GalleryProps) => {
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
    keyScopeRef.current?.focus({ preventScroll: true });
  }, []);

  useEffect(() => {
    if (!getSelectionKey() && flattenedItems.length > 0) {
      const firstItem = flattenedItems[0];
      if (firstItem) {
        setSelection(getItemKey(firstItem));
      }
    }
  }, [flattenedItems, getSelectionKey, setSelection]);

  // --------------------------------------
  const selectedKey = getSelectionKey();
  useScrollSelectedIntoView(scrollRef, selectedKey, { block: 'nearest' });
  useEffect(() => {
    keyScopeRef.current?.focus({ preventScroll: true });
  }, [selectedKey, flattenedItems]);

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
      if (!selectedKey) return;

      const { columns } = controls;
      let nextItem: dalle.DalleDress | null = null;

      // Define handlers for different keys
      switch (e.key) {
        case 'ArrowRight': {
          const idx = flattenedItems.findIndex(
            (f) => getItemKey(f) === selectedKey,
          );
          if (idx !== -1) {
            nextItem =
              flattenedItems[
                (idx + 1 + flattenedItems.length) % flattenedItems.length
              ] || null;
          }
          break;
        }
        case 'ArrowLeft': {
          const idx = flattenedItems.findIndex(
            (f) => getItemKey(f) === selectedKey,
          );
          if (idx !== -1) {
            nextItem =
              flattenedItems[
                (idx - 1 + flattenedItems.length) % flattenedItems.length
              ] || null;
          }
          break;
        }
        case 'ArrowDown':
        case 'ArrowUp': {
          let currentGroupIndex = -1;
          let currentItemIndexInGroup = -1;

          // Find the current item's "coordinates"
          for (let i = 0; i < groupNames.length; i++) {
            const groupName = groupNames[i];
            if (!groupName) continue;
            const group = groupedItems[groupName] || [];
            const itemIndex = group.findIndex(
              (item: dalle.DalleDress) => getItemKey(item) === selectedKey,
            );
            if (itemIndex !== -1) {
              currentGroupIndex = i;
              currentItemIndexInGroup = itemIndex;
              break;
            }
          }

          if (currentGroupIndex === -1) return;

          const currentColumn = currentItemIndexInGroup % columns;
          const currentGroupName = groupNames[currentGroupIndex];
          const currentGroup = currentGroupName
            ? groupedItems[currentGroupName] || []
            : [];

          if (e.key === 'ArrowDown') {
            const nextItemIndexInGroup = currentItemIndexInGroup + columns;
            if (nextItemIndexInGroup < currentGroup.length) {
              nextItem = currentGroup[nextItemIndexInGroup] || null;
            } else {
              const nextGroupIndex =
                (currentGroupIndex + 1) % groupNames.length;
              const nextGroupName = groupNames[nextGroupIndex];
              const nextGroup = nextGroupName
                ? groupedItems[nextGroupName] || []
                : [];
              if (nextGroup.length > 0) {
                const targetIndex = Math.min(
                  currentColumn,
                  nextGroup.length - 1,
                );
                nextItem = nextGroup[targetIndex] || null;
              }
            }
          } else {
            // ArrowUp
            const prevItemIndexInGroup = currentItemIndexInGroup - columns;
            if (prevItemIndexInGroup >= 0) {
              nextItem = currentGroup[prevItemIndexInGroup] || null;
            } else {
              const prevGroupIndex =
                (currentGroupIndex - 1 + groupNames.length) % groupNames.length;
              const prevGroupName = groupNames[prevGroupIndex];
              const prevGroup = prevGroupName
                ? groupedItems[prevGroupName] || []
                : [];
              if (prevGroup.length > 0) {
                const lastRowStartIndex =
                  Math.floor((prevGroup.length - 1) / columns) * columns;
                const targetIndex = Math.min(
                  lastRowStartIndex + currentColumn,
                  prevGroup.length - 1,
                );
                nextItem = prevGroup[targetIndex] || null;
              }
            }
          }
          break;
        }
        case 'Home':
          nextItem = flattenedItems[0] || null;
          break;
        case 'End':
          nextItem = flattenedItems[flattenedItems.length - 1] || null;
          break;
        case 'Enter': {
          const item = flattenedItems.find(
            (f) => getItemKey(f) === selectedKey,
          );
          if (item) {
            handleItemDoubleClick(item);
          }
          return; // No selection change, just action
        }
        default:
          return; // Not a key we handle
      }

      if (nextItem) {
        e.preventDefault();
        setSelection(getItemKey(nextItem));
      }
    },
    [
      flattenedItems,
      getSelectionKey,
      handleItemDoubleClick,
      setSelection,
      controls,
      groupNames,
      groupedItems,
    ],
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
            selected={selectedKey}
          />
        ))}
      </div>
    </Container>
  );
};
