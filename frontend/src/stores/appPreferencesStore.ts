import {
  ConvertToAddress,
  GetActiveProjectData,
  GetAppPreferences,
  SetActiveAddress,
  SetActiveChain,
  SetActiveContract,
  SetAppPreferences,
  SetLastFacet,
  SetLastView,
} from '@app';
import { preferences, types } from '@models';
import { Log } from '@utils';

export interface UseActiveProjectReturn {
  lastTheme: string;
  lastLanguage: string;
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  debugMode: boolean;
  showDetailPanel: boolean;
  loading: boolean;
  lastProject: string;
  activeChain: string;
  activeAddress: string;
  activeContract: string;
  lastView: string;
  lastFacetMap: Record<string, types.DataFacet>;
  setActiveAddress: (address: string) => Promise<void>;
  setActiveChain: (chain: string) => Promise<void>;
  setActiveContract: (contract: string) => Promise<void>;
  setLastView: (view: string) => Promise<void>;
  setLastFacet: (view: string, facet: types.DataFacet) => Promise<void>;
  switchProject: (project: string) => Promise<void>;
  toggleTheme: () => Promise<void>;
  changeLanguage: (language: string) => Promise<void>;
  setMenuCollapsed: (collapsed: boolean) => Promise<void>;
  setHelpCollapsed: (collapsed: boolean) => Promise<void>;
  toggleDebugMode: () => Promise<void>;
  setShowDetailPanel: (enabled: boolean) => Promise<void>;
  isDarkMode: boolean;
  hasActiveProject: boolean;
  canExport: boolean;
  effectiveAddress: string;
  effectiveChain: string;
}

interface AppPreferencesState {
  lastTheme: string;
  lastLanguage: string;
  menuCollapsed: boolean;
  helpCollapsed: boolean;
  debugMode: boolean;
  showDetailPanel: boolean;
  lastProject: string;
  loading: boolean;
  activeChain: string;
  activeAddress: string;
  activeContract: string;
  lastView: string;
  lastFacetMap: Record<string, types.DataFacet>;
  canExport: boolean;
}

// Initial state - values will be populated from active project context
const initialState: AppPreferencesState = {
  canExport: false,
  lastFacetMap: {},
  activeAddress: '', // Will be set from active project
  activeContract: '',
  activeChain: '', // Will be set from active project
  lastTheme: 'light',
  lastLanguage: 'en',
  menuCollapsed: false,
  helpCollapsed: false,
  debugMode: false,
  showDetailPanel: false,
  lastView: '/',
  lastProject: '',
  loading: false,
};

class AppPreferencesStore {
  private state: AppPreferencesState = { ...initialState };
  private listeners = new Set<() => void>();
  private isTestMode = false;
  private cachedSnapshot: UseActiveProjectReturn | null = null;

  setTestMode = (testMode: boolean): void => {
    this.isTestMode = testMode;
    if (testMode) {
      this.setState({ loading: false });
    }
  };

  getSnapshot = (): UseActiveProjectReturn => {
    if (!this.cachedSnapshot) {
      this.cachedSnapshot = {
        lastTheme: this.state.lastTheme,
        lastLanguage: this.state.lastLanguage,
        menuCollapsed: this.state.menuCollapsed,
        helpCollapsed: this.state.helpCollapsed,
        debugMode: this.state.debugMode,
        showDetailPanel: this.state.showDetailPanel,
        loading: this.state.loading,
        lastProject: this.state.lastProject,
        activeChain: this.state.activeChain,
        activeAddress: this.state.activeAddress,
        activeContract: this.state.activeContract,
        lastView: this.state.lastView,
        lastFacetMap: this.state.lastFacetMap,
        setActiveAddress: this.setActiveAddress,
        setActiveChain: this.setActiveChain,
        setActiveContract: this.setActiveContract,
        setLastView: this.setLastView,
        setLastFacet: this.setLastFacet,
        switchProject: this.switchProject,
        toggleTheme: this.toggleTheme,
        changeLanguage: this.changeLanguage,
        setMenuCollapsed: this.setMenuCollapsed,
        setHelpCollapsed: this.setHelpCollapsed,
        toggleDebugMode: this.toggleDebugMode,
        setShowDetailPanel: this.setShowDetailPanel,
        isDarkMode: this.isDarkMode,
        hasActiveProject: this.hasActiveProject,
        canExport: this.canExport,
        effectiveAddress: this.state.activeAddress,
        effectiveChain: this.state.activeChain,
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

  private setState = (updates: Partial<AppPreferencesState>): void => {
    this.state = { ...this.state, ...updates };
    this.cachedSnapshot = null; // Invalidate cache when state changes
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

      // Use single call to get all project data and app preferences
      const [prefs, projectData] = await Promise.all([
        GetAppPreferences(),
        GetActiveProjectData(),
      ]);

      this.setState({
        lastProject: prefs.lastProject || '',
        lastTheme: prefs.lastTheme || 'dark',
        lastLanguage: prefs.lastLanguage || 'en',
        menuCollapsed: prefs.menuCollapsed || false,
        helpCollapsed: prefs.helpCollapsed || false,
        debugMode:
          (prefs as preferences.AppPreferences & { debugMode?: boolean })
            .debugMode || false,
        showDetailPanel: prefs.showDetailPanel || false,
        // Use project data directly - no hardcoded fallbacks
        activeChain: projectData.activeChain || '',
        activeAddress: projectData.activeAddress || '',
        activeContract: projectData.activeContract || '',
        lastView: projectData.lastView || '/',
        // Convert string map to DataFacet map
        lastFacetMap: Object.fromEntries(
          Object.entries(projectData.lastFacetMap || {}).map(([key, value]) => [
            key,
            value as types.DataFacet,
          ]),
        ),
        loading: false,
      });
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

  setShowDetailPanel = async (enabled: boolean): Promise<void> => {
    await this.updatePreferences({ showDetailPanel: enabled });
    this.setState({ showDetailPanel: enabled });
  };

  setActiveAddress = async (address: string): Promise<void> => {
    const convertedAddress = await ConvertToAddress(address);
    if (convertedAddress && typeof convertedAddress === 'object') {
      await SetActiveAddress(convertedAddress);
      this.setState({ activeAddress: address });
    } else {
      throw new Error('Invalid address format');
    }
  };

  setActiveChain = async (chain: string): Promise<void> => {
    await SetActiveChain(chain);
    this.setState({ activeChain: chain });
  };

  setActiveContract = async (contract: string): Promise<void> => {
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
    await SetLastView(view);
    this.setState({ lastView: view });
  };

  get isDarkMode(): boolean {
    return this.state.lastTheme === 'dark';
  }

  get hasActiveProject(): boolean {
    return Boolean(
      this.state.lastProject &&
        this.state.activeAddress &&
        this.state.activeChain,
    );
  }

  get canExport(): boolean {
    return Boolean(
      this.state.lastProject &&
        this.state.activeAddress &&
        this.state.activeChain,
    );
  }
}

// Create and export a singleton instance
export const appPreferencesStore = new AppPreferencesStore();

// Initialize the store when the module loads (only when not in test environment)
// Use setTimeout to ensure all modules are loaded first
if (
  typeof window !== 'undefined' &&
  typeof import.meta.env.VITEST === 'undefined'
) {
  setTimeout(() => {
    appPreferencesStore.initialize();
  }, 0);
} else {
  // In test environment, enable test mode
  appPreferencesStore.setTestMode(true);
}
