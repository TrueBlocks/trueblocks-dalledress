// ADD_ROUTE
import { types } from '@models';

export const ABIS_ROUTE = '/abis';
export const ABIS_DEFAULT_LIST = types.ListKind.DOWNLOADED;

export const ACTION_MESSAGES = {
  RELOAD_STATUS: 'Reloaded ABI data. Fetching fresh data...',
  DELETE_SUCCESS: (address: string) =>
    `Address ${address} was deleted successfully`,
  DELETE_FAILURE: (address: string, error: string) =>
    `Failed to delete address ${address}: ${error}`,
} as const;

// ADD_ROUTE
