import { useSyncExternalStore } from 'react';

import {
  GetAppPreferences,
  SetAppPreferences,
  SetHelpCollapsed,
  SetMenuCollapsed,
} from '@app';
import { preferences } from '@models';
import { Log } from '@utils';

export interface UseUIStateReturn2 {
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  showDetailPanel: boolean;
  loading: boolean;
  setMenuCollapsed: (collapsed: boolean) => Promise<void>;
  setHelpCollapsed: (collapsed: boolean) => Promise<void>;
  setShowDetailPanel: (enabled: boolean) => Promise<void>;
}

interface UIState {
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  showDetailPanel: boolean;
  loading: boolean;
}

const initialUIState: UIState = {
  menuCollapsed: false,
  helpCollapsed: false,
  showDetailPanel: false,
  loading: false,
};

class UIStateStore {
  private state: UIState = { ...initialUIState };
  private listeners = new Set<() => void>();
  private isTestMode = false;
  private cachedSnapshot: UseUIStateReturn2 | null = null;

  setTestMode = (testMode: boolean): void => {
    this.isTestMode = testMode;
    if (testMode) {
      this.setState({ loading: false });
    }
  };

  getSnapshot = (): UseUIStateReturn2 => {
    if (!this.cachedSnapshot) {
      this.cachedSnapshot = {
        menuCollapsed: this.state.menuCollapsed,
        helpCollapsed: this.state.helpCollapsed,
        showDetailPanel: this.state.showDetailPanel,
        loading: this.state.loading,
        setMenuCollapsed: this.setMenuCollapsed,
        setHelpCollapsed: this.setHelpCollapsed,
        setShowDetailPanel: this.setShowDetailPanel,
      };
    }
    return this.cachedSnapshot;
  };

  subscribe = (listener: () => void): (() => void) => {
    this.listeners.add(listener);
    return () => {
      this.listeners.delete(listener);
    };
  };

  private notify = (): void => {
    this.listeners.forEach((listener) => listener());
  };

  private setState = (updates: Partial<UIState>): void => {
    this.state = { ...this.state, ...updates };
    this.cachedSnapshot = null;
    this.notify();
  };

  private updatePreferences = async (
    updates: Partial<preferences.AppPreferences>,
  ): Promise<void> => {
    if (this.isTestMode) return;

    try {
      const currentPrefs = await GetAppPreferences();
      const updatedPrefs = preferences.AppPreferences.createFrom({
        ...currentPrefs,
        ...updates,
      });
      await SetAppPreferences(updatedPrefs);
    } catch (error) {
      Log('ERROR: Failed to update UI state preferences: ' + String(error));
      throw error;
    }
  };

  initialize = async (): Promise<void> => {
    if (this.isTestMode) return;

    try {
      this.setState({ loading: true });

      const prefs = await GetAppPreferences();

      this.setState({
        menuCollapsed: prefs.menuCollapsed || false,
        helpCollapsed: prefs.helpCollapsed || false,
        showDetailPanel: prefs.showDetailPanel || false,
        loading: false,
      });
    } catch (error) {
      Log('ERROR: Failed to load UI state: ' + String(error));
      this.setState({ loading: false });
    }
  };

  setMenuCollapsed = async (collapsed: boolean): Promise<void> => {
    await SetMenuCollapsed(collapsed);
    await this.updatePreferences({ menuCollapsed: collapsed });
    this.setState({ menuCollapsed: collapsed });
  };

  setHelpCollapsed = async (collapsed: boolean): Promise<void> => {
    await SetHelpCollapsed(collapsed);
    await this.updatePreferences({ helpCollapsed: collapsed });
    this.setState({ helpCollapsed: collapsed });
  };

  setShowDetailPanel = async (enabled: boolean): Promise<void> => {
    await this.updatePreferences({ showDetailPanel: enabled });
    this.setState({ showDetailPanel: enabled });
  };
}

const uiStateStore = new UIStateStore();

if (
  typeof window !== 'undefined' &&
  typeof import.meta.env.VITEST === 'undefined'
) {
  setTimeout(() => {
    uiStateStore.initialize();
  }, 0);
} else {
  uiStateStore.setTestMode(true);
}

export const useUIState = (): UseUIStateReturn2 => {
  return useSyncExternalStore(uiStateStore.subscribe, uiStateStore.getSnapshot);
};
