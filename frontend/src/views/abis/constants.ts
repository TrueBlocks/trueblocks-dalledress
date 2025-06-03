import { types } from '@models';

// Default list kind when component first loads
export const DEFAULT_LIST_KIND = types.ListKind.DOWNLOADED;

// Default tab route
export const ABIS_ROUTE = '/abis';

// Action messages for UI feedback
export const ACTION_MESSAGES = {
  RELOAD_STATUS: 'Reloaded ABI data. Fetching fresh data...',
  DELETE_SUCCESS: (address: string) =>
    `Address ${address} was deleted successfully`,
  DELETE_FAILURE: (address: string, error: string) =>
    `Failed to delete address ${address}: ${error}`,
} as const;
