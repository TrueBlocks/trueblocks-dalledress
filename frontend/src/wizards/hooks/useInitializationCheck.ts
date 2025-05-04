import { useEffect } from 'react';

import { GetUserInfoStatus, Logger } from '@app';
import { useAppContext } from '@contexts';
import { checkAndNavigateToWizard } from '@utils';

import { WizardState } from '..';

export const useInitializationCheck = (
  state: WizardState,
  loadInitialData: () => Promise<void>,
) => {
  const { navigate, isWizard } = useAppContext();

  useEffect(() => {
    loadInitialData();
  }, [loadInitialData]);

  useEffect(() => {
    if (state.ui.initialLoading) return;

    const checkWizardState = async () => {
      try {
        await checkAndNavigateToWizard(navigate, isWizard);
      } catch (err) {
        Logger('Failed to check wizard state: ' + err);
        if (!isWizard) {
          try {
            navigate('/wizard');
          } catch (navError) {
            Logger('Failed to navigate to wizard: ' + navError);
            window.location.href = '/wizard';
          }
        }
      }
    };

    checkWizardState();
    const interval = setInterval(checkWizardState, 1500);
    return () => clearInterval(interval);
  }, [state.ui.initialLoading, isWizard, navigate, state.api]);

  const verifyCompletionStep = async () => {
    if (state.ui.activeStep === 2 && !state.ui.initialLoading) {
      try {
        const wizardState = await GetUserInfoStatus();
        if (wizardState.missingNameEmail) {
          return 0;
        } else if (wizardState.rpcUnavailable) {
          return 1;
        }
        return 2;
      } catch (error) {
        Logger('Failed to verify completion step: ' + error);
        return state.ui.activeStep;
      }
    }
    return state.ui.activeStep;
  };

  return {
    verifyCompletionStep,
  };
};
