import { useCallback, useMemo, useSyncExternalStore } from 'react';

import { dalle } from '@models';

export const getItemKey = (item: dalle.DalleDress | null): string => {
  if (!item) return '';
  return `${item.original}:${item.series || ''}`;
};

interface GalleryState {
  selectedKey: string | null;
  orig: string | null;
  series: string | null;
  galleryItems: dalle.DalleDress[];
  groupedBySeries: Record<string, dalle.DalleDress[]>;
  groupedByAddress: Record<string, dalle.DalleDress[]>;
}

const initial: GalleryState = {
  selectedKey: null,
  orig: null,
  series: null,
  galleryItems: [],
  groupedBySeries: {},
  groupedByAddress: {},
};

class GalleryStore {
  private state: GalleryState = { ...initial };
  private listeners = new Set<() => void>();
  private setState(updates: Partial<GalleryState>) {
    this.state = { ...this.state, ...updates };
    this.listeners.forEach((l) => l());
  }

  subscribe = (listener: () => void) => {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  };

  getSnapshot = (): GalleryState => this.state;

  setSelection(key: string | null) {
    if (!key) {
      this.setState({ selectedKey: null, orig: null, series: null });
      return;
    }
    const found = this.state.galleryItems.find((i) => getItemKey(i) === key);
    if (found) {
      this.setState({
        selectedKey: key,
        orig: found.original,
        series: found.series,
      });
    }
  }

  clear() {
    this.setState({ ...initial });
  }

  ingest(items: dalle.DalleDress[] | null) {
    const list = items ? [...items] : [];
    const groupedBySeries: Record<string, dalle.DalleDress[]> = {};
    const groupedByAddress: Record<string, dalle.DalleDress[]> = {};
    for (const it of list) {
      const sKey = it.series || '';
      if (!groupedBySeries[sKey]) groupedBySeries[sKey] = [];
      groupedBySeries[sKey].push(it);
      const aKey = it.original || '';
      if (!groupedByAddress[aKey]) groupedByAddress[aKey] = [];
      groupedByAddress[aKey].push(it);
    }
    this.setState({ galleryItems: list, groupedBySeries, groupedByAddress });
  }
}

const store = new GalleryStore();

export const useGalleryStore = () => {
  const sel = useSyncExternalStore(store.subscribe, store.getSnapshot);
  const getSelectionKey = useCallback(() => sel.selectedKey, [sel.selectedKey]);
  const getSelectedItem = useCallback(() => {
    const k = sel.selectedKey;
    if (!k) return null;
    return sel.galleryItems.find((i) => getItemKey(i) === k) || null;
  }, [sel.selectedKey, sel.galleryItems]);
  const setSelection = useCallback(
    (key: string | null) => store.setSelection(key),
    [],
  );
  const clearSelection = useCallback(() => store.clear(), []);
  const ingestItems = useCallback(
    (items: dalle.DalleDress[] | null) => store.ingest(items),
    [],
  );
  const useDerived = (sortMode: 'series' | 'address') => {
    const { groupedBySeries, groupedByAddress } = sel;
    const groupNames = useMemo(() => {
      if (sortMode === 'series')
        return Object.keys(groupedBySeries).sort((a, b) => a.localeCompare(b));
      return Object.keys(groupedByAddress).sort((a, b) => a.localeCompare(b));
    }, [sortMode, groupedBySeries, groupedByAddress]);
    const groupedItems = useMemo(() => {
      if (sortMode === 'series') return groupedBySeries;
      return groupedByAddress;
    }, [sortMode, groupedBySeries, groupedByAddress]);
    return { groupNames, groupedItems };
  };
  const handleKey = useCallback(
    (
      e: React.KeyboardEvent<HTMLDivElement>,
      items: dalle.DalleDress[],
      columns?: number,
      onDoubleClick?: (item: dalle.DalleDress) => void,
      groupNames?: Array<string>,
      groupedItems?: Record<string, dalle.DalleDress[]>,
    ) => {
      if (!items.length) return;
      const selectedKey = getSelectionKey();
      if (!selectedKey) return;
      let nextIdx = items.findIndex((g) => getItemKey(g) === selectedKey);
      // Find current group and position
      let groupIdx = 0,
        itemIdxInGroup = 0;
      if (groupNames && groupedItems) {
        let found = false;
        for (let g = 0; g < groupNames.length; g++) {
          const groupKey = groupNames[g] ?? '';
          const group: Array<dalle.DalleDress> = groupedItems[groupKey] || [];
          const idx = group.findIndex(
            (i: dalle.DalleDress) => getItemKey(i) === selectedKey,
          );
          if (idx !== -1) {
            groupIdx = g;
            itemIdxInGroup = idx;
            found = true;
            break;
          }
        }
        if (!found) {
          groupIdx = 0;
          itemIdxInGroup = 0;
        }
      }
      if (e.key === 'ArrowRight') {
        nextIdx = (nextIdx + 1 + items.length) % items.length;
        e.preventDefault();
      } else if (e.key === 'ArrowLeft') {
        nextIdx = (nextIdx - 1 + items.length) % items.length;
        e.preventDefault();
      } else if (e.key === 'Home') {
        nextIdx = 0;
        e.preventDefault();
      } else if (e.key === 'End') {
        nextIdx = items.length - 1;
        e.preventDefault();
      } else if (e.key === 'Enter' && onDoubleClick) {
        const item = items.find((g) => getItemKey(g) === selectedKey);
        if (item) onDoubleClick(item);
        return;
      } else if (
        columns &&
        (e.key === 'ArrowDown' || e.key === 'ArrowUp') &&
        groupNames &&
        groupedItems
      ) {
        // Grid navigation for ragged rows using groupedItems
        const groupKey = groupNames[groupIdx] ?? '';
        const group: Array<dalle.DalleDress> = groupedItems[groupKey] || [];
        const row = Math.floor(itemIdxInGroup / columns);
        const col = itemIdxInGroup % columns;
        let targetRow = e.key === 'ArrowDown' ? row + 1 : row - 1;
        let targetGroupIdx = groupIdx;
        let targetItemIdxInGroup = null;
        let totalRows = Math.ceil(group.length / columns);
        if (targetRow >= 0 && targetRow < totalRows) {
          const start = targetRow * columns;
          const end = Math.min(start + columns, group.length);
          const rowLength = end - start;
          const clampedCol = Math.min(col, rowLength - 1);
          targetItemIdxInGroup = start + clampedCol;
        } else if (targetRow < 0 && groupIdx > 0) {
          // Move to previous group, last row
          targetGroupIdx = groupIdx - 1;
          const prevGroupKey = groupNames[targetGroupIdx] ?? '';
          const prevGroup: Array<dalle.DalleDress> =
            groupedItems[prevGroupKey] || [];
          const prevTotalRows = Math.ceil(prevGroup.length / columns);
          const start = (prevTotalRows - 1) * columns;
          const end = prevGroup.length;
          const rowLength = end - start;
          const clampedCol = Math.min(col, rowLength - 1);
          targetItemIdxInGroup = start + clampedCol;
        } else if (targetRow >= totalRows && groupIdx < groupNames.length - 1) {
          // Move to next group, first row
          targetGroupIdx = groupIdx + 1;
          const nextGroupKey = groupNames[targetGroupIdx] ?? '';
          const nextGroup: Array<dalle.DalleDress> =
            groupedItems[nextGroupKey] || [];
          const start = 0;
          const end = Math.min(columns, nextGroup.length);
          const rowLength = end - start;
          const clampedCol = Math.min(col, rowLength - 1);
          targetItemIdxInGroup = start + clampedCol;
        }
        if (targetItemIdxInGroup !== null) {
          const targetGroupKey = groupNames[targetGroupIdx] ?? '';
          const targetGroup: Array<dalle.DalleDress> =
            groupedItems[targetGroupKey] || [];
          const targetItem = targetGroup[targetItemIdxInGroup] ?? null;
          if (targetItem) {
            setSelection(getItemKey(targetItem));
            e.preventDefault();
          }
        }
      } else {
        return;
      }
      const next = items[nextIdx];
      if (
        next &&
        (e.key === 'ArrowRight' ||
          e.key === 'ArrowLeft' ||
          e.key === 'Home' ||
          e.key === 'End')
      ) {
        setSelection(getItemKey(next));
      }
    },
    [getSelectionKey, setSelection],
  );
  return {
    orig: sel.orig,
    series: sel.series,
    galleryItems: sel.galleryItems,
    getSelectionKey,
    getSelectedItem,
    setSelection,
    clearSelection,
    ingestItems,
    useDerived,
    handleKey,
  };
};
