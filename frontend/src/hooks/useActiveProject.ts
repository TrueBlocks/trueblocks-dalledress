import { useSyncExternalStore } from 'react';

import { types } from '@models';
import { appPreferencesStore } from '@stores';

export interface UseActiveProjectReturn {
  // State
  lastTheme: string;
  lastLanguage: string;
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  debugMode: boolean;
  loading: boolean;

  // Project-specific state (from active project - computed via backend calls)
  lastProject: string;
  activeChain: string;
  activeAddress: string;
  activeContract: string;
  lastView: string;
  lastFacetMap: Record<string, types.DataFacet>;

  // Actions that call backend directly
  setActiveAddress: (address: string) => Promise<void>;
  switchProject: (project: string) => Promise<void>;
  toggleTheme: () => Promise<void>;
  changeLanguage: (language: string) => Promise<void>;
  setMenuCollapsed: (collapsed: boolean) => Promise<void>;
  setHelpCollapsed: (collapsed: boolean) => Promise<void>;
  setLastView: (view: string) => Promise<void>;
  setLastFacet: (view: string, facet: types.DataFacet) => Promise<void>;
  setActiveChain: (chain: string) => Promise<void>;
  setActiveContract: (contract: string) => Promise<void>;
  isDarkMode: boolean;
  toggleDebugMode: () => Promise<void>;

  // Computed values (sync calls to store)
  hasActiveProject: boolean;
  canExport: boolean;
}

export const useActiveProject = (): UseActiveProjectReturn => {
  const state = useSyncExternalStore(
    appPreferencesStore.subscribe,
    appPreferencesStore.getState,
  );

  return {
    // State from preferences store
    lastTheme: state.lastTheme,
    lastLanguage: state.lastLanguage,
    menuCollapsed: state.menuCollapsed,
    helpCollapsed: state.helpCollapsed,
    debugMode: state.debugMode,
    loading: state.loading,

    // Project-specific state (from active project)
    activeAddress: state.activeAddress,
    lastProject: state.lastProject,
    activeChain: state.activeChain,
    activeContract: state.activeContract,
    lastView: state.lastView,
    lastFacetMap: state.lastFacetMap,

    // Actions - enhanced project-aware methods
    setActiveAddress: appPreferencesStore.setActiveAddress,
    setActiveChain: appPreferencesStore.setActiveChain,
    setActiveContract: appPreferencesStore.setActiveContract,
    switchProject: appPreferencesStore.switchProject,
    toggleTheme: appPreferencesStore.toggleTheme,
    changeLanguage: appPreferencesStore.changeLanguage,
    setMenuCollapsed: appPreferencesStore.setMenuCollapsed,
    setHelpCollapsed: appPreferencesStore.setHelpCollapsed,
    setLastView: appPreferencesStore.setLastView,
    setLastFacet: appPreferencesStore.setLastFacet,
    toggleDebugMode: appPreferencesStore.toggleDebugMode,

    // Computed values - these are sync getters from the store
    isDarkMode: appPreferencesStore.isDarkMode,
    hasActiveProject: appPreferencesStore.hasActiveProject,
    canExport: appPreferencesStore.canExport,
  };
};
