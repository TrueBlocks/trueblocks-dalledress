import { useSyncExternalStore } from 'react';

import { UseActiveProjectReturn, appPreferencesStore } from '@stores';

export type { UseActiveProjectReturn };

export const useActiveProject = (): UseActiveProjectReturn => {
  return useSyncExternalStore(
    appPreferencesStore.subscribe,
    appPreferencesStore.getSnapshot,
  );
};
