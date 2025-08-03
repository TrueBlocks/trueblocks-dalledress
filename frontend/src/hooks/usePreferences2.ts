import { useSyncExternalStore } from 'react';

import {
  GetAppPreferences,
  SetAppPreferences,
  SetDebugMode,
  SetLanguage,
  SetTheme,
} from '@app';
import { preferences } from '@models';
import { Log } from '@utils';

export interface UsePreferencesReturn2 {
  lastTheme: string;
  lastLanguage: string;
  debugMode: boolean;
  loading: boolean;
  toggleTheme: () => Promise<void>;
  changeLanguage: (language: string) => Promise<void>;
  toggleDebugMode: () => Promise<void>;
  isDarkMode: boolean;
}

interface PreferencesState {
  lastTheme: string;
  lastLanguage: string;
  debugMode: boolean;
  loading: boolean;
}

const initialPreferencesState: PreferencesState = {
  lastTheme: 'dark',
  lastLanguage: 'en',
  debugMode: false,
  loading: false,
};

class PreferencesStore {
  private state: PreferencesState = { ...initialPreferencesState };
  private listeners = new Set<() => void>();
  private isTestMode = false;
  private cachedSnapshot: UsePreferencesReturn2 | null = null;

  setTestMode = (testMode: boolean): void => {
    this.isTestMode = testMode;
    if (testMode) {
      this.setState({ loading: false });
    }
  };

  getSnapshot = (): UsePreferencesReturn2 => {
    if (!this.cachedSnapshot) {
      this.cachedSnapshot = {
        lastTheme: this.state.lastTheme,
        lastLanguage: this.state.lastLanguage,
        debugMode: this.state.debugMode,
        loading: this.state.loading,
        toggleTheme: this.toggleTheme,
        changeLanguage: this.changeLanguage,
        toggleDebugMode: this.toggleDebugMode,
        isDarkMode: this.isDarkMode,
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

  private setState = (updates: Partial<PreferencesState>): void => {
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
      Log('ERROR: Failed to update preferences: ' + String(error));
      throw error;
    }
  };

  initialize = async (): Promise<void> => {
    if (this.isTestMode) return;

    try {
      this.setState({ loading: true });

      const prefs = await GetAppPreferences();

      this.setState({
        lastTheme: prefs.lastTheme || 'dark',
        lastLanguage: prefs.lastLanguage || 'en',
        debugMode: prefs.debugMode || false,
        loading: false,
      });
    } catch (error) {
      Log('ERROR: Failed to load preferences: ' + String(error));
      this.setState({ loading: false });
    }
  };

  toggleTheme = async (): Promise<void> => {
    const newTheme = this.state.lastTheme === 'dark' ? 'light' : 'dark';
    await SetTheme(newTheme);
    await this.updatePreferences({ lastTheme: newTheme });
    this.setState({ lastTheme: newTheme });
  };

  changeLanguage = async (language: string): Promise<void> => {
    await SetLanguage(language);
    await this.updatePreferences({ lastLanguage: language });
    this.setState({ lastLanguage: language });
  };

  toggleDebugMode = async (): Promise<void> => {
    const newDebugMode = !this.state.debugMode;
    await SetDebugMode(newDebugMode);
    await this.updatePreferences({ debugMode: newDebugMode });
    this.setState({ debugMode: newDebugMode });
  };

  get isDarkMode(): boolean {
    return this.state.lastTheme === 'dark';
  }
}

const preferencesStore = new PreferencesStore();

if (
  typeof window !== 'undefined' &&
  typeof import.meta.env.VITEST === 'undefined'
) {
  setTimeout(() => {
    preferencesStore.initialize();
  }, 0);
} else {
  preferencesStore.setTestMode(true);
}

export const usePreferences2 = (): UsePreferencesReturn2 => {
  return useSyncExternalStore(
    preferencesStore.subscribe,
    preferencesStore.getSnapshot,
  );
};
