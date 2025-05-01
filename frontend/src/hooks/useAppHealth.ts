import { useEffect } from 'react';

import { useAppContext } from '@contexts';
import { checkAndNavigateToWizard } from '@utils';

export const useAppHealth = () => {
  const { ready, isWizard, navigate } = useAppContext();

  useEffect(() => {
    if (!ready) return;

    const interval = setInterval(() => {
      checkAndNavigateToWizard(navigate, isWizard);
    }, 1500);

    return () => clearInterval(interval);
  }, [ready, isWizard, navigate]);
};
