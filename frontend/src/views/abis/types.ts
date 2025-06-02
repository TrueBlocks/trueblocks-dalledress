import { types } from '@models';

//--------------------------------------------------------------------
export type IndexedAbi = types.Abi & {
  [key: string]: unknown;
};

export type IndexedFunction = types.Function & {
  [key: string]: unknown;
};

export type AbiRow = (types.Abi | types.Function) & {
  [key: string]: unknown;
};

//--------------------------------------------------------------------
export interface TableConfigProps {
  downloaded: IndexedAbi[];
  known: IndexedAbi[];
  functions: IndexedFunction[];
  events: IndexedFunction[];
  isDownloadedLoaded: boolean;
  isKnownLoaded: boolean;
  isFuncsLoaded: boolean;
  isEventsLoaded: boolean;
  processingAddresses: Set<string>;
  setSelectedAddress: (address: string) => void;
  setLocation: (path: string) => void;
  handleAction: (address: string) => void;
}
