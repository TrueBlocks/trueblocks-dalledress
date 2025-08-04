import { useSyncExternalStore } from 'react';

import {
  ConvertToAddress,
  GetActiveProjectData,
  SetActiveAddress,
  SetActiveChain,
  SetActiveContract,
  SetLastFacet,
  SetLastView,
} from '@app';
import { types } from '@models';
import { Log, addressToHex } from '@utils';

export interface UseActiveProjectReturn {
  loading: boolean;
  lastProject: string;
  activeChain: string;
  activeAddress: string;
  activeContract: string;
  lastView: string;
  lastFacetMap: Record<string, types.DataFacet>;
  getLastFacet: (view: string) => string;
  setActiveAddress: (address: string) => Promise<void>;
  setActiveChain: (chain: string) => Promise<void>;
  setActiveContract: (contract: string) => Promise<void>;
  setLastView: (view: string) => Promise<void>;
  setLastFacet: (view: string, facet: types.DataFacet) => Promise<void>;
  switchProject: (project: string) => Promise<void>;
  hasActiveProject: boolean;
  canExport: boolean;
  effectiveAddress: string;
  effectiveChain: string;
}

interface ProjectState {
  loading: boolean;
  lastProject: string;
  activeChain: string;
  activeAddress: string;
  activeContract: string;
  lastView: string;
  lastFacetMap: Record<string, types.DataFacet>;
}

const initialProjectState: ProjectState = {
  loading: false,
  lastProject: '',
  activeChain: '',
  activeAddress: '',
  activeContract: '',
  lastView: '/',
  lastFacetMap: {},
};

class ProjectStore {
  private state: ProjectState = { ...initialProjectState };
  private listeners = new Set<() => void>();
  private cachedSnapshot: UseActiveProjectReturn | null = null;

  private getLastFacet = (view: string): string => {
    const vR = view.replace(/^\/+/, '');
    return this.state.lastFacetMap[vR] || '';
  };

  getSnapshot = (): UseActiveProjectReturn => {
    if (!this.cachedSnapshot) {
      this.cachedSnapshot = {
        loading: this.state.loading,
        lastProject: this.state.lastProject,
        activeChain: this.state.activeChain,
        activeAddress: this.state.activeAddress,
        activeContract: this.state.activeContract,
        lastView: this.state.lastView,
        lastFacetMap: this.state.lastFacetMap,
        getLastFacet: this.getLastFacet,
        setActiveAddress: this.setActiveAddress,
        setActiveChain: this.setActiveChain,
        setActiveContract: this.setActiveContract,
        setLastView: this.setLastView,
        setLastFacet: this.setLastFacet,
        switchProject: this.switchProject,
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

  private setState = (updates: Partial<ProjectState>): void => {
    this.state = { ...this.state, ...updates };
    this.cachedSnapshot = null;
    this.notify();
  };

  initialize = async (): Promise<void> => {
    try {
      this.setState({ loading: true });

      const projectData = await GetActiveProjectData();

      this.setState({
        activeChain: projectData.activeChain || '',
        activeAddress: projectData.activeAddress || '',
        activeContract: projectData.activeContract || '',
        lastView: projectData.lastView || '/',
        lastFacetMap: Object.fromEntries(
          Object.entries(projectData.lastFacetMap || {}).map(([key, value]) => [
            key,
            value as types.DataFacet,
          ]),
        ),
        loading: false,
      });
    } catch (error) {
      Log('ERROR: Failed to load project data: ' + String(error));
      this.setState({ loading: false });
    }
  };

  switchProject = async (project: string): Promise<void> => {
    this.setState({ lastProject: project });
  };

  setActiveAddress = async (address: string): Promise<void> => {
    try {
      const result = await ConvertToAddress(address);
      if (result && typeof result === 'object') {
        const hexAddress = addressToHex(result);
        await SetActiveAddress(hexAddress);
        this.setState({ activeAddress: hexAddress });
      } else {
        const errorMsg = `Invalid address - ConvertToAddress returned: ${JSON.stringify(result)}`;
        Log(`ERROR: ${errorMsg}`);
        throw new Error(errorMsg);
      }
    } catch (error) {
      Log(`ERROR: Failed to set active address: ${String(error)}`);
      throw error;
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

const projectStore = new ProjectStore();

if (
  typeof window !== 'undefined' &&
  typeof import.meta.env.VITEST === 'undefined'
) {
  setTimeout(() => {
    projectStore.initialize();
  }, 0);
}

export const useActiveProject = (): UseActiveProjectReturn => {
  return useSyncExternalStore(projectStore.subscribe, projectStore.getSnapshot);
};
