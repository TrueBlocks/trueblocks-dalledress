import { useEffect, useRef } from 'react';

import { useLocation } from 'wouter';

import { useActiveProject } from './useActiveProject';
import { useAppReadiness } from './useAppReadiness';

export const useAppNavigation = () => {
  const [location, navigate] = useLocation();
  const ready = useAppReadiness();
  const hasRedirected = useRef(false);
  const { lastView, setLastView } = useActiveProject();

  const isWizard = location.startsWith('/wizard');

  // Handle initial redirect to lastView
  useEffect(() => {
    if (ready && location === '/' && !hasRedirected.current) {
      hasRedirected.current = true;
      navigate(lastView);
    }
  }, [ready, location, lastView, navigate]);

  // Sync location changes to the preferences store
  useEffect(() => {
    if (ready && location !== '/') {
      setLastView(location);
    }
  }, [location, ready, setLastView]);

  return {
    location,
    navigate,
    ready,
    isWizard,
  };
};
