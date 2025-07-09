import { useSyncExternalStore } from 'react';

import { types } from '@models';
import { appPreferencesStore } from '@stores';

export interface UseActiveProjectReturn {
  // State
  lastProject: string;
  lastChain: string;
  lastAddress: string;
  lastTheme: string;
  lastLanguage: string;
  lastView: string;
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  lastTab: Record<string, types.DataFacet>;
  loading: boolean;

  // Actions
  setActiveAddress: (address: string) => Promise<void>;
  setActiveChain: (chain: string) => Promise<void>;
  switchProject: (project: string) => Promise<void>;
  toggleTheme: () => Promise<void>;
  changeLanguage: (language: string) => Promise<void>;
  setMenuCollapsed: (collapsed: boolean) => Promise<void>;
  setHelpCollapsed: (collapsed: boolean) => Promise<void>;
  setLastTab: (route: string, tab: types.DataFacet) => Promise<void>;
  setLastView: (view: string) => Promise<void>;
  isDarkMode: boolean;
  toggleDarkMode: () => Promise<void>;

  // Computed values
  hasActiveProject: boolean;
  canExport: boolean;
}

export const useActiveProject = (): UseActiveProjectReturn => {
  const state = useSyncExternalStore(
    appPreferencesStore.subscribe,
    appPreferencesStore.getState,
  );

  return {
    // State
    lastProject: state.lastProject,
    lastChain: state.lastChain,
    lastAddress: state.lastAddress,
    lastTheme: state.lastTheme,
    lastLanguage: state.lastLanguage,
    lastView: state.lastView,
    menuCollapsed: state.menuCollapsed,
    helpCollapsed: state.helpCollapsed,
    lastTab: state.lastTab,
    loading: state.loading,

    // Actions - these are bound methods from the store
    setActiveAddress: appPreferencesStore.setActiveAddress,
    setActiveChain: appPreferencesStore.setActiveChain,
    switchProject: appPreferencesStore.switchProject,
    toggleTheme: appPreferencesStore.toggleTheme,
    changeLanguage: appPreferencesStore.changeLanguage,
    setMenuCollapsed: appPreferencesStore.setMenuCollapsed,
    setHelpCollapsed: appPreferencesStore.setHelpCollapsed,
    setLastTab: appPreferencesStore.setLastTab,
    setLastView: appPreferencesStore.setLastView,
    toggleDarkMode: appPreferencesStore.toggleDarkMode,

    // Computed values - these are getters from the store
    isDarkMode: appPreferencesStore.isDarkMode,
    hasActiveProject: appPreferencesStore.hasActiveProject,
    canExport: appPreferencesStore.canExport,
  };
};
