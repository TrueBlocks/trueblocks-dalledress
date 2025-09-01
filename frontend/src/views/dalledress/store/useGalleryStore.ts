import { useCallback, useMemo, useSyncExternalStore } from 'react';

import { dalle } from '@models';

interface GalleryState {
  orig: string | null;
  series: string | null;
  galleryItems: dalle.DalleDress[];
  groupedBySeries: Record<string, dalle.DalleDress[]>;
  groupedByAddress: Record<string, dalle.DalleDress[]>;
}

const initial: GalleryState = {
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

  setAll(orig: string | null, series: string | null) {
    this.setState({ orig, series });
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
  const getSelection = useCallback(() => {
    if (!sel.series || !sel.orig) return null;
    return `./${sel.series}/annotated/${sel.orig}.png`;
  }, [sel.series, sel.orig]);
  const setSelection = useCallback(
    (orig: string | null, series: string | null) => store.setAll(orig, series),
    [],
  );
  const clearSelection = useCallback(() => store.clear(), []);
  const ingestItems = useCallback(
    (items: dalle.DalleDress[] | null) => store.ingest(items),
    [],
  );
  const useDerived = (sortMode: 'series' | 'address') => {
    const groupNames = useMemo(() => {
      if (sortMode === 'series')
        return Object.keys(sel.groupedBySeries).sort((a, b) =>
          a.localeCompare(b),
        );
      return Object.keys(sel.groupedByAddress).sort((a, b) =>
        a.localeCompare(b),
      );
    }, [sortMode]);
    const groupedItems = useMemo(() => {
      if (sortMode === 'series') return sel.groupedBySeries;
      return sel.groupedByAddress;
    }, [sortMode]);
    const flattenedItems = useMemo(
      () => groupNames.flatMap((s) => groupedItems[s] || []),
      [groupNames, groupedItems],
    );
    return { groupNames, groupedItems, flattenedItems };
  };
  return {
    orig: sel.orig,
    series: sel.series,
    galleryItems: sel.galleryItems,
    getSelection,
    setSelection,
    clearSelection,
    ingestItems,
    useDerived,
  };
};
