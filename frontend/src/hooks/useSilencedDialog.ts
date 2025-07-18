import { useCallback, useEffect, useState } from 'react';

import { GetAppPreferences } from '@app';

/**
 * Hook to check if a dialog has been silenced by the user
 * @param dialogKey - The key identifier for the dialog
 * @returns boolean indicating if the dialog is silenced
 */
export const useSilencedDialog = (dialogKey: string) => {
  const [isSilenced, setIsSilenced] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const checkSilenced = useCallback(async () => {
    try {
      setIsLoading(true);
      const prefs = await GetAppPreferences();
      const silenced = prefs.silencedDialogs?.[dialogKey] === true;
      setIsSilenced(silenced);
    } catch (error) {
      console.error('Failed to check silenced dialog:', error);
      setIsSilenced(false); // Default to not silenced if we can't check
    } finally {
      setIsLoading(false);
    }
  }, [dialogKey]);

  useEffect(() => {
    checkSilenced();
  }, [checkSilenced]);

  return { isSilenced, isLoading, recheckSilenced: checkSilenced };
};
