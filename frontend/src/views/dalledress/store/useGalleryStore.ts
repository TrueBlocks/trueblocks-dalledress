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
    const flattenedItems = useMemo(
      () => groupNames.flatMap((s) => groupedItems[s] || []),
      [groupNames, groupedItems],
    );
    return { groupNames, groupedItems, flattenedItems };
  };
  const handleKey = useCallback(
    (
      e: React.KeyboardEvent<HTMLDivElement>,
      items: dalle.DalleDress[],
      columns?: number,
      onDoubleClick?: (item: dalle.DalleDress) => void,
    ) => {
      if (!items.length) return;
      const selectedKey = getSelectionKey();
      if (!selectedKey) return;
      let nextIdx = items.findIndex((g) => getItemKey(g) === selectedKey);
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
      } else if (columns && (e.key === 'ArrowDown' || e.key === 'ArrowUp')) {
        // Grid navigation for Gallery
        const idx = nextIdx;
        if (e.key === 'ArrowDown') {
          nextIdx = Math.min(idx + columns, items.length - 1);
          e.preventDefault();
        } else if (e.key === 'ArrowUp') {
          nextIdx = Math.max(idx - columns, 0);
          e.preventDefault();
        }
      } else {
        return;
      }
      const next = items[nextIdx];
      if (next) setSelection(getItemKey(next));
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
