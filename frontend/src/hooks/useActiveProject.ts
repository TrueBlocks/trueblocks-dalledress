import { useCallback, useEffect, useState } from 'react';

import * as App from '@app';
import { preferences, types } from '@models';

interface ProjectState {
  // Core project state
  lastProject: string;
  lastChain: string;
  lastAddress: string;

  // UI state
  lastTheme: string;
  lastLanguage: string;
  lastView: string;
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  lastTab: Record<string, types.ListKind>;

  // Loading state
  loading: boolean;
}

interface ProjectActions {
  // Address/Chain actions
  setActiveAddress: (address: string) => Promise<void>;
  setActiveChain: (chain: string) => Promise<void>;

  // Project actions
  switchProject: (project: string) => Promise<void>;

  // UI actions
  toggleTheme: () => Promise<void>;
  changeLanguage: (language: string) => Promise<void>;
  setMenuCollapsed: (collapsed: boolean) => Promise<void>;
  setHelpCollapsed: (collapsed: boolean) => Promise<void>;
  setLastTab: (route: string, tab: types.ListKind) => Promise<void>;
  setLastView: (view: string) => Promise<void>;

  // Computed theme action
  isDarkMode: boolean;
  toggleDarkMode: () => Promise<void>;
}

interface ComputedValues {
  // Computed project state
  hasActiveProject: boolean;
  canExport: boolean;
  effectiveAddress: string;
  effectiveChain: string;
}

export interface UseActiveProjectReturn
  extends ProjectState,
    ProjectActions,
    ComputedValues {}

export const useActiveProject = (): UseActiveProjectReturn => {
  const [state, setState] = useState<ProjectState>({
    lastProject: '',
    lastChain: '',
    lastAddress: '',
    lastTheme: 'dark',
    lastLanguage: 'en',
    lastView: '/',
    menuCollapsed: true,
    helpCollapsed: true,
    lastTab: {},
    loading: true,
  });

  // Load initial state from backend
  useEffect(() => {
    const loadState = async () => {
      try {
        setState((prev) => ({ ...prev, loading: true }));
        const prefs = await App.GetAppPreferences();

        setState({
          lastProject: prefs.lastProject || '',
          lastChain: prefs.lastChain || '',
          lastAddress: prefs.lastAddress || '',
          lastTheme: prefs.lastTheme || 'dark',
          lastLanguage: prefs.lastLanguage || 'en',
          lastView: prefs.lastView || '/',
          menuCollapsed: prefs.menuCollapsed || false,
          helpCollapsed: prefs.helpCollapsed || false,
          lastTab: (prefs.lastTab || {}) as Record<string, types.ListKind>,
          loading: false,
        });
      } catch (error) {
        console.error('Failed to load app preferences:', error);
        setState((prev) => ({ ...prev, loading: false }));
      }
    };

    loadState();
  }, []);

  // Helper function to update preferences
  const updatePreferences = useCallback(
    async (updates: Partial<preferences.AppPreferences>) => {
      try {
        const currentPrefs = await App.GetAppPreferences();
        // Create a new AppPreferences instance instead of directly spreading
        const updatedPrefs = preferences.AppPreferences.createFrom({
          ...currentPrefs,
          ...updates,
        });
        await App.SetAppPreferences(updatedPrefs);
      } catch (error) {
        console.error('Failed to update preferences:', error);
        throw error;
      }
    },
    [],
  );

  // Address/Chain actions
  const setActiveAddress = useCallback(
    async (address: string) => {
      await updatePreferences({ lastAddress: address });
      setState((prev) => ({ ...prev, lastAddress: address }));
    },
    [updatePreferences],
  );

  const setActiveChain = useCallback(
    async (chain: string) => {
      await updatePreferences({ lastChain: chain });
      setState((prev) => ({ ...prev, lastChain: chain }));
    },
    [updatePreferences],
  );

  // Project actions
  const switchProject = useCallback(
    async (project: string) => {
      await updatePreferences({ lastProject: project });
      setState((prev) => ({ ...prev, lastProject: project }));
    },
    [updatePreferences],
  );

  // UI actions
  const toggleTheme = useCallback(async () => {
    const newTheme = state.lastTheme === 'dark' ? 'light' : 'dark';
    await updatePreferences({ lastTheme: newTheme });
    setState((prev) => ({ ...prev, lastTheme: newTheme }));
  }, [state.lastTheme, updatePreferences]);

  const changeLanguage = useCallback(
    async (language: string) => {
      await updatePreferences({ lastLanguage: language });
      setState((prev) => ({ ...prev, lastLanguage: language }));
    },
    [updatePreferences],
  );

  const setMenuCollapsed = useCallback(
    async (collapsed: boolean) => {
      await updatePreferences({ menuCollapsed: collapsed });
      setState((prev) => ({ ...prev, menuCollapsed: collapsed }));
    },
    [updatePreferences],
  );

  const setHelpCollapsed = useCallback(
    async (collapsed: boolean) => {
      await updatePreferences({ helpCollapsed: collapsed });
      setState((prev) => ({ ...prev, helpCollapsed: collapsed }));
    },
    [updatePreferences],
  );

  const setLastTab = useCallback(async (route: string, tab: types.ListKind) => {
    // Use the backend SetLastTab method for this
    await App.SetLastTab(route, tab);
    setState((prev) => ({
      ...prev,
      lastTab: { ...prev.lastTab, [route]: tab },
    }));
  }, []);

  const setLastView = useCallback(
    async (view: string) => {
      await updatePreferences({ lastView: view });
      setState((prev) => ({ ...prev, lastView: view }));
    },
    [updatePreferences],
  );

  // Theme computed values and actions
  const isDarkMode = state.lastTheme === 'dark';
  const toggleDarkMode = useCallback(async () => {
    const newTheme = isDarkMode ? 'light' : 'dark';
    await updatePreferences({ lastTheme: newTheme });
    setState((prev) => ({ ...prev, lastTheme: newTheme }));
  }, [isDarkMode, updatePreferences]);

  // Computed values
  const hasActiveProject = Boolean(state.lastProject);
  const canExport = Boolean(state.lastProject && state.lastAddress);
  const effectiveAddress = state.lastAddress || '';
  const effectiveChain = state.lastChain || 'mainnet';

  return {
    // State
    ...state,

    // Actions
    setActiveAddress,
    setActiveChain,
    switchProject,
    toggleTheme,
    changeLanguage,
    setMenuCollapsed,
    setHelpCollapsed,
    setLastTab,
    setLastView,
    isDarkMode,
    toggleDarkMode,

    // Computed values
    hasActiveProject,
    canExport,
    effectiveAddress,
    effectiveChain,
  };
};
