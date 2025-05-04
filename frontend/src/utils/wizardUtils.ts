import { GetUserInfoStatus, IsInitialized, Logger } from '@app';

const disabled = true;

export const checkAndNavigateToWizard = async (
  navigate: ReturnType<typeof import('wouter').useLocation>[1],
  inWizard: boolean,
): Promise<number | void> => {
  try {
    if (disabled) {
      return;
    }

    const initialized = await IsInitialized();
    if (!initialized) {
      if (!inWizard) {
        try {
          navigate('/wizard', { replace: true });
        } catch (navError) {
          Logger('Failed to navigate to wizard: ' + navError);
          window.location.href = '/wizard';
        }
      }
      return;
    }

    const state = await GetUserInfoStatus();
    const needsWizard = state.missingNameEmail || state.rpcUnavailable;
    if (needsWizard && !inWizard) {
      try {
        navigate('/wizard', { replace: true });
      } catch (navError) {
        Logger('Failed to navigate to wizard: ' + navError);
        window.location.href = '/wizard';
      }
    }

    return state.missingNameEmail ? 0 : state.rpcUnavailable ? 1 : 2;
  } catch (err) {
    Logger('Failed to check wizard state: ' + err);
    if (!inWizard) {
      try {
        navigate('/wizard', { replace: true });
      } catch (navError) {
        Logger('Failed to navigate to wizard: ' + navError);
        window.location.href = '/wizard';
      }
    }
  }
};
