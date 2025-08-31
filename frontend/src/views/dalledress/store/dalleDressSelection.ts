import { useSyncExternalStore } from 'react';

interface SelectionState {
  orig: string | null;
  series: string | null;
  path: string | null;
}

const initial: SelectionState = {
  orig: null,
  series: null,
  path: null,
};

class SelectionStore {
  private state: SelectionState = { ...initial };
  private listeners = new Set<() => void>();
  private setState(updates: Partial<SelectionState>) {
    this.state = { ...this.state, ...updates };
    this.listeners.forEach((l) => l());
  }

  subscribe = (listener: () => void) => {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  };

  getSnapshot = (): SelectionState => this.state;

  setOriginal(orig: string | null) {
    this.setState({ orig });
  }

  setSeries(series: string | null) {
    this.setState({ series });
  }

  setPath(path: string | null) {
    this.setState({ path });
  }

  setAll(orig: string | null, series: string | null, path: string | null) {
    this.setState({ orig, series, path });
  }

  clear() {
    this.setState({ ...initial });
  }
}

const store = new SelectionStore();

export const useDalleDressSelection = () => {
  const sel = useSyncExternalStore(store.subscribe, store.getSnapshot);
  return {
    ...sel,
    setDressOriginal: (orig: string | null) => store.setOriginal(orig),
    setDressSeries: (series: string | null) => store.setSeries(series),
    setDressPath: (path: string | null) => store.setPath(path),
    setDressSelection: (
      orig: string | null,
      series: string | null,
      path: string | null,
    ) => store.setAll(orig, series, path),
    clearDressSelection: () => store.clear(),
  };
};

// test-only helper
export const __resetDressSelectionForTests = () => store.clear();
