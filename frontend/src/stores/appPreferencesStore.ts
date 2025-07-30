import {
  ConvertToAddress,
  GetActiveProject,
  GetAppPreferences,
  SetActiveAddress,
  SetActiveContract,
  SetAppPreferences,
  SetLastFacet,
  SetLastView,
} from '@app';
import { preferences, types } from '@models';
import { Log, getAddressString } from '@utils';
import { SetActiveChain } from 'wailsjs/go/project/Project';

// State interface for simplified preferences (no caching of project data)
interface AppPreferencesState {
  // UI state (from global preferences)
  lastTheme: string;
  lastLanguage: string;
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  lastFacetMap: Record<string, types.DataFacet>;

  // Project state
  lastProject: string;
  activeChain: string;
  activeAddress: string;
  activeContract: string;
  lastView: string;

  // Debug state
  debugMode: boolean;

  // Loading state
  loading: boolean;
}

// Initial state
const initialState: AppPreferencesState = {
  lastProject: '',
  activeChain: '',
  activeAddress: '0xf503017d7baf7fbc0fff7492b751025c6a78179b',
  activeContract: '0x52df6e4d9989e7cf4739d687c765e75323a1b14c',
  lastTheme: 'dark',
  lastLanguage: 'en',
  lastView: '/',
  menuCollapsed: true,
  helpCollapsed: true,
  lastFacetMap: {},
  debugMode: false,
  loading: true,
};

class AppPreferencesStore {
  private state: AppPreferencesState = { ...initialState };
  private listeners = new Set<() => void>();
  private isTestMode = false;

  // Set test mode to skip backend calls
  setTestMode = (testMode: boolean): void => {
    this.isTestMode = testMode;
    if (testMode) {
      // In test mode, set loading to false and use defaults
      this.setState({ loading: false });
    }
  };

  // Get current state (required by useSyncExternalStore)
  getState = (): AppPreferencesState => {
    return this.state;
  };

  // Subscribe to state changes (required by useSyncExternalStore)
  subscribe = (listener: () => void): (() => void) => {
    this.listeners.add(listener);

    // Return unsubscribe function
    return () => {
      this.listeners.delete(listener);
    };
  };

  // Notify all listeners of state changes
  private notify = (): void => {
    this.listeners.forEach((listener) => listener());
  };

  // Internal method to update state and notify listeners
  private setState = (updates: Partial<AppPreferencesState>): void => {
    this.state = { ...this.state, ...updates };
    this.notify();
  };

  // Helper function to update backend preferences
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

  // Initialize the store by loading state from backend
  initialize = async (): Promise<void> => {
    // Skip initialization in test mode
    if (this.isTestMode) {
      return;
    }

    try {
      this.setState({ loading: true });
      let retries = 0;
      const maxRetries = 10;
      let prefs = null;
      while (retries < maxRetries) {
        try {
          prefs = await GetAppPreferences();
          break; // Success, exit retry loop
        } catch (error) {
          retries++;
          Log(
            `App preferences load attempt ${retries}/${maxRetries} failed: ${JSON.stringify(error)}`,
          );
          if (retries >= maxRetries) {
            Log(
              'WARN: Failed to load app preferences after max retries, using defaults',
            );
            // Use safe defaults if backend is not ready
            this.setState({
              lastTheme: 'dark',
              lastLanguage: 'en',
              menuCollapsed: false,
              helpCollapsed: false,
              debugMode: false,
              loading: false,
              lastView: '/',
              lastFacetMap: {},
              lastProject: '',
              activeChain: '',
              activeAddress: '0xf503017d7baf7fbc0fff7492b751025c6a78179b',
            });
            return;
          }
          // Wait before retry
          await new Promise((resolve) => setTimeout(resolve, 200));
        }
      }

      if (prefs) {
        let lastView = '/';
        let activeChain = '';
        let activeAddress = '0xf503017d7baf7fbc0fff7492b751025c6a78179b';
        let activeContract = '0x52df6e4d9989e7cf4739d687c765e75323a1b14c';
        let lastFacetMap: Record<string, types.DataFacet> = {};

        try {
          const activeProject = await GetActiveProject();
          lastView = activeProject.lastView || '/';
          activeChain = activeProject.activeChain || '';
          activeAddress =
            getAddressString(activeProject.activeAddress) ||
            '0xf503017d7baf7fbc0fff7492b751025c6a78179b';
          activeContract =
            activeProject.activeContract ||
            '0x52df6e4d9989e7cf4739d687c765e75323a1b14c';
          lastFacetMap = (activeProject.lastFacetMap || {}) as Record<
            string,
            types.DataFacet
          >;
        } catch (error) {
          Log(
            'WARN: Failed to get project properties from active project, using defaults: ' +
              String(error),
          );
        }

        this.setState({
          lastProject: prefs.lastProject || '',
          activeChain: activeChain,
          activeAddress: activeAddress,
          activeContract: activeContract,
          lastTheme: prefs.lastTheme || 'dark',
          lastLanguage: prefs.lastLanguage || 'en',
          lastView: lastView,
          menuCollapsed: prefs.menuCollapsed || false,
          helpCollapsed: prefs.helpCollapsed || false,
          lastFacetMap: lastFacetMap,
          debugMode:
            (prefs as preferences.AppPreferences & { debugMode?: boolean })
              .debugMode || false,
          loading: false,
        });
      }
    } catch (error) {
      Log('ERROR: Failed to load app preferences: ' + String(error));
      this.setState({ loading: false });
    }
  };

  switchProject = async (project: string): Promise<void> => {
    await this.updatePreferences({ lastProject: project });
    this.setState({ lastProject: project });
  };

  toggleTheme = async (): Promise<void> => {
    const newTheme = this.state.lastTheme === 'dark' ? 'light' : 'dark';
    await this.updatePreferences({ lastTheme: newTheme });
    this.setState({ lastTheme: newTheme });
  };

  changeLanguage = async (language: string): Promise<void> => {
    await this.updatePreferences({ lastLanguage: language });
    this.setState({ lastLanguage: language });
  };

  setMenuCollapsed = async (collapsed: boolean): Promise<void> => {
    await this.updatePreferences({ menuCollapsed: collapsed });
    this.setState({ menuCollapsed: collapsed });
  };

  setHelpCollapsed = async (collapsed: boolean): Promise<void> => {
    await this.updatePreferences({ helpCollapsed: collapsed });
    this.setState({ helpCollapsed: collapsed });
  };

  toggleDebugMode = async (): Promise<void> => {
    const newDebugMode = !this.state.debugMode;
    await this.updatePreferences({ debugMode: newDebugMode });
    this.setState({ debugMode: newDebugMode });
  };

  // Action methods that update both local state and backend
  setActiveAddress = async (address: string): Promise<void> => {
    const convertedAddress = await ConvertToAddress(address);
    if (convertedAddress && typeof convertedAddress === 'object') {
      // await this.updatePreferences({ activeAddress: address });
      await SetActiveAddress(convertedAddress);
      this.setState({ activeAddress: address });
    } else {
      throw new Error('Invalid address format');
    }
  };

  setActiveChain = async (chain: string): Promise<void> => {
    // await this.updatePreferences({ activeChain: chain });
    await SetActiveChain(chain);
    this.setState({ activeChain: chain });
  };

  setActiveContract = async (contract: string): Promise<void> => {
    // await this.updatePreferences({ activeContract: contract });
    await SetActiveContract(contract);
    this.setState({ activeContract: contract });
  };

  setLastFacet = async (view: string, facet: string): Promise<void> => {
    await SetLastFacet(view, facet);
    this.setState({
      lastFacetMap: {
        ...this.state.lastFacetMap,
        [view]: facet as types.DataFacet,
      },
    });
  };

  setLastView = async (view: string): Promise<void> => {
    // await this.updatePreferences({ lastView: view });
    await SetLastView(view);
    this.setState({ lastView: view });
  };

  // Computed getters
  get isDarkMode(): boolean {
    return this.state.lastTheme === 'dark';
  }

  get hasActiveProject(): boolean {
    return Boolean(this.state.lastProject);
  }

  get canExport(): boolean {
    return Boolean(this.state.lastProject && this.state.activeAddress);
  }
}

// Create and export a singleton instance
export const appPreferencesStore = new AppPreferencesStore();

// Initialize the store when the module loads (only when not in test environment)
if (
  typeof window !== 'undefined' &&
  typeof import.meta.env.VITEST === 'undefined'
) {
  appPreferencesStore.initialize();
} else {
  // In test environment, enable test mode
  appPreferencesStore.setTestMode(true);
}
