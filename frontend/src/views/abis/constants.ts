import { types } from '@models';

// Default list kind when component first loads
export const DEFAULT_LIST_KIND = types.ListKind.DOWNLOADED;

// State setter map to eliminate switch statement repetition
export const LOADER_SETTERS = {
  [types.ListKind.DOWNLOADED]: 'setIsDownloadedLoaded',
  [types.ListKind.KNOWN]: 'setIsKnownLoaded',
  [types.ListKind.FUNCTIONS]: 'setIsFuncsLoaded',
  [types.ListKind.EVENTS]: 'setIsEventsLoaded',
} as const;

// Data setter map for fetchData switch statement
export const DATA_SETTERS = {
  [types.ListKind.DOWNLOADED]: 'setDownloadedAbis',
  [types.ListKind.KNOWN]: 'setKnownAbis',
  [types.ListKind.FUNCTIONS]: 'setFunctions',
  [types.ListKind.EVENTS]: 'setEvents',
} as const;

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
