// MONITORS_ROUTE
import { types } from '@models';

export const MONITORS_ROUTE = '/monitors';
export const MONITORS_DEFAULT_LIST = types.ListKind.MONITORS;

export const ACTION_MESSAGES = {
  RELOAD_STATUS: 'Reloaded monitor data. Fetching fresh data...',
  DELETE_SUCCESS: (address: string) =>
    `Monitor ${address} was deleted successfully`,
  UNDELETE_SUCCESS: (address: string) =>
    `Monitor ${address} was undeleted successfully`,
  REMOVE_SUCCESS: (address: string) =>
    `Monitor ${address} was removed successfully`,
  CLEAN_SUCCESS: (count: number) => `Cleaned ${count} monitor(s) successfully`,
  DELETE_FAILURE: (address: string, error: string) =>
    `Failed to delete monitor ${address}: ${error}`,
  UNDELETE_FAILURE: (address: string, error: string) =>
    `Failed to undelete monitor ${address}: ${error}`,
  REMOVE_FAILURE: (address: string, error: string) =>
    `Failed to remove monitor ${address}: ${error}`,
  CLEAN_FAILURE: (error: string) => `Failed to clean monitors: ${error}`,
} as const;

// MONITORS_ROUTE
