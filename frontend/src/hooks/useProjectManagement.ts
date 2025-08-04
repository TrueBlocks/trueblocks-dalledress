import { useSyncExternalStore } from 'react';

import {
  ClearActiveProject,
  CloseProject,
  GetAppPreferences,
  GetOpenProjects,
  NewProject,
  OpenProjectFile,
  SetAppPreferences,
  SwitchToProject,
} from '@app';
import { preferences } from '@models';
import { Log } from '@utils';

export interface UseProjectManagementReturn {
  lastProject: string;
  projects: Record<string, unknown>[];
  loading: boolean;
  switchProject: (projectId: string) => Promise<void>;
  newProject: (name: string, address: string) => Promise<void>;
  openProjectFile: (path: string) => Promise<void>;
  closeProject: (projectId: string) => Promise<void>;
  clearActiveProject: () => Promise<void>;
  refreshProjects: () => Promise<void>;
}

interface ProjectManagementState {
  lastProject: string;
  projects: Record<string, unknown>[];
  loading: boolean;
}

const initialProjectManagementState: ProjectManagementState = {
  lastProject: '',
  projects: [],
  loading: false,
};

class ProjectManagementStore {
  private state: ProjectManagementState = { ...initialProjectManagementState };
  private listeners = new Set<() => void>();
  private cachedSnapshot: UseProjectManagementReturn | null = null;

  getSnapshot = (): UseProjectManagementReturn => {
    if (!this.cachedSnapshot) {
      this.cachedSnapshot = {
        lastProject: this.state.lastProject,
        projects: this.state.projects,
        loading: this.state.loading,
        switchProject: this.switchProject,
        newProject: this.newProject,
        openProjectFile: this.openProjectFile,
        closeProject: this.closeProject,
        clearActiveProject: this.clearActiveProject,
        refreshProjects: this.refreshProjects,
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

  private setState = (updates: Partial<ProjectManagementState>): void => {
    this.state = { ...this.state, ...updates };
    this.cachedSnapshot = null;
    this.notify();
  };

  private updatePreferences = async (
    updates: Partial<preferences.AppPreferences>,
  ): Promise<void> => {
    try {
      const currentPrefs = await GetAppPreferences();
      const updatedPrefs = preferences.AppPreferences.createFrom({
        ...currentPrefs,
        ...updates,
      });
      await SetAppPreferences(updatedPrefs);
    } catch (error) {
      Log(
        'ERROR: Failed to update project management preferences: ' +
          String(error),
      );
      throw error;
    }
  };

  initialize = async (): Promise<void> => {
    try {
      this.setState({ loading: true });

      const [prefs, projects] = await Promise.all([
        GetAppPreferences(),
        GetOpenProjects(),
      ]);

      this.setState({
        lastProject: prefs.lastProject || '',
        projects: projects || [],
        loading: false,
      });
    } catch (error) {
      Log('ERROR: Failed to load project management data: ' + String(error));
      this.setState({ loading: false });
    }
  };

  refreshProjects = async (): Promise<void> => {
    try {
      const projects = await GetOpenProjects();
      this.setState({ projects: projects || [] });
    } catch (error) {
      Log('ERROR: Failed to refresh projects: ' + String(error));
      throw error;
    }
  };

  switchProject = async (projectId: string): Promise<void> => {
    await SwitchToProject(projectId);
    await this.updatePreferences({ lastProject: projectId });
    this.setState({ lastProject: projectId });
    await this.refreshProjects();
  };

  newProject = async (name: string, address: string): Promise<void> => {
    await NewProject(name, address);
    await this.refreshProjects();
  };

  openProjectFile = async (path: string): Promise<void> => {
    await OpenProjectFile(path);
    await this.refreshProjects();
  };

  closeProject = async (projectId: string): Promise<void> => {
    await CloseProject(projectId);
    await this.refreshProjects();
  };

  clearActiveProject = async (): Promise<void> => {
    await ClearActiveProject();
    await this.updatePreferences({ lastProject: '' });
    this.setState({ lastProject: '' });
    await this.refreshProjects();
  };
}

const projectManagementStore = new ProjectManagementStore();

if (
  typeof window !== 'undefined' &&
  typeof import.meta.env.VITEST === 'undefined'
) {
  setTimeout(() => {
    projectManagementStore.initialize();
  }, 0);
}

export const useProjectManagement = (): UseProjectManagementReturn => {
  return useSyncExternalStore(
    projectManagementStore.subscribe,
    projectManagementStore.getSnapshot,
  );
};
