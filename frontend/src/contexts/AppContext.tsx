import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
} from 'react';

import {
  GetAppPreferences,
  IsReady,
  Logger,
  SetLastTab,
  SetLastView,
} from '@app';
import { useLocation } from 'wouter';

interface AppContextType {
  isDarkMode: boolean;
  toggleDarkMode: () => void;
  currentLocation: string;
  navigate: (to: string) => void;
  isWizard: boolean;
  ready: boolean;
  menuCollapsed: boolean;
  setMenuCollapsed: (collapsed: boolean) => void;
  helpCollapsed: boolean;
  setHelpCollapsed: (collapsed: boolean) => void;
  lastTab: Record<string, string>;
  setLastTab: (route: string, tab: string) => void;
}

export const AppContext = createContext<AppContextType | undefined>(undefined);

export const AppContextProvider = ({ children }: { children: ReactNode }) => {
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [location, navigate] = useLocation();
  const [lastView, setLastView] = useState('/');
  const [ready, setReady] = useState(false);
  const [menuCollapsed, setMenuCollapsed] = useState(true);
  const [helpCollapsed, setHelpCollapsed] = useState(true);
  const [lastTab, setLastTabState] = useState<Record<string, string>>({});
  const hasRedirected = useRef(false);

  const toggleDarkMode = () => {
    setIsDarkMode((prev) => !prev);
  };

  const setLastTab = (route: string, tab: string) => {
    setLastTabState((prev) => {
      const updatedState = {
        ...prev,
        [route]: tab,
      };
      SetLastTab(route, tab).catch((error) => {
        Logger(`Failed to update lastTab in backend: ${error}`);
      });
      return updatedState;
    });
  };

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
          setMenuCollapsed(prefs.menuCollapsed as boolean);
          setHelpCollapsed(prefs.helpCollapsed as boolean);
          setLastTabState(prefs.lastTab || {});
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

  useEffect(() => {
    if (ready) {
      SetLastView(location);
      setLastView(location);
    }
  }, [location, ready]);

  return (
    <AppContext.Provider
      value={{
        isDarkMode,
        toggleDarkMode,
        currentLocation: location,
        navigate,
        isWizard,
        ready,
        menuCollapsed,
        setMenuCollapsed,
        helpCollapsed,
        setHelpCollapsed,
        lastTab,
        setLastTab,
      }}
    >
      {children}
    </AppContext.Provider>
  );
};

export const useAppContext = () => {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error('useAppContext must be used within an AppContextProvider');
  }
  return context;
};
