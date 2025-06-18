import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';

import { GetAppPreferences, IsReady } from '@app';
import { useLocation } from 'wouter';

interface AppContextType {
  currentLocation: string;
  navigate: (to: string) => void;
  isWizard: boolean;
  ready: boolean;
}

export const AppContext = createContext<AppContextType | undefined>(undefined);

export const AppContextProvider = ({ children }: { children: ReactNode }) => {
  const [location, navigate] = useLocation();
  const [lastView, setLastView] = useState('/');
  const [ready, setReady] = useState(false);
  const hasRedirected = useRef(false);

  const isWizard = location.startsWith('/wizard');

  useEffect(() => {
    const initializeApp = async () => {
      let attempts = 0;
      const maxAttempts = 200;
      while (attempts < maxAttempts) {
        const isReady = await IsReady();
        if (isReady) {
          const prefs = await GetAppPreferences();
          setLastView(prefs.lastView || '/');
          setReady(true);
          return;
        }

        await new Promise((resolve) => setTimeout(resolve, 50));
        attempts++;
      }
    };
    initializeApp();
  }, []);

  useEffect(() => {
    if (ready && location === '/' && !hasRedirected.current) {
      hasRedirected.current = true;
      navigate(lastView);
    }
  }, [ready, location, lastView, navigate]);

  const contextValue = useMemo(
    () => ({
      currentLocation: location,
      navigate,
      isWizard,
      ready,
    }),
    [location, navigate, isWizard, ready],
  );

  return (
    <AppContext.Provider value={contextValue}>{children}</AppContext.Provider>
  );
};

export const useAppContext = () => {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error('useAppContext must be used within an AppContextProvider');
  }
  return context;
};
