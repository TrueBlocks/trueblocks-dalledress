import { useCallback, useSyncExternalStore } from 'react';

interface GalleryState {
  orig: string | null;
  series: string | null;
}

const initial: GalleryState = {
  orig: null,
  series: null,
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

  setOriginal(orig: string | null) {
    this.setState({ orig });
  }

  setSeries(series: string | null) {
    this.setState({ series });
  }

  setAll(orig: string | null, series: string | null) {
    this.setState({ orig, series });
  }

  clear() {
    this.setState({ ...initial });
  }
}

const store = new GalleryStore();

export const useGalleryStore = () => {
  const sel = useSyncExternalStore(store.subscribe, store.getSnapshot);
  const getPath = useCallback(() => {
    if (!sel.series || !sel.orig) return null;
    return `./${sel.series}/annotated/${sel.orig}.png`;
  }, [sel.series, sel.orig]);
  const setDressOriginal = useCallback(
    (orig: string | null) => store.setOriginal(orig),
    [],
  );
  const setDressSeries = useCallback(
    (series: string | null) => store.setSeries(series),
    [],
  );
  const setDressSelection = useCallback(
    (orig: string | null, series: string | null) => store.setAll(orig, series),
    [],
  );
  const clearDressSelection = useCallback(() => store.clear(), []);
  return {
    ...sel,
    getPath,
    setDressOriginal,
    setDressSeries,
    setDressSelection,
    clearDressSelection,
  };
};
