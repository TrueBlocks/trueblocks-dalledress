import { useEffect } from 'react';

import { NodeStatus, getBarWidth } from '@components';
import { ViewContextProvider, useAppContext } from '@contexts';
import { useAppHealth, useAppHotkeys } from '@hooks';
import { Footer, Header, HelpBar, MainView, MenuBar } from '@layout';
import { AppShell } from '@mantine/core';
import { Router } from 'wouter';

// Add at the top level, outside the component
function globalNavKeySquelcher(e: KeyboardEvent) {
  const navKeys = [
    'ArrowUp',
    'ArrowDown',
    'ArrowLeft',
    'ArrowRight',
    'PageUp',
    'PageDown',
    'Home',
    'End',
    'Enter',
    'Escape',
  ];
  if (navKeys.includes(e.key)) {
    // Only squelch if not handled by a focused control
    e.preventDefault();
  }
}

export const App = () => {
  useEffect(() => {
    window.addEventListener('keydown', globalNavKeySquelcher, {
      capture: true,
    });
    return () => {
      window.removeEventListener('keydown', globalNavKeySquelcher, {
        capture: true,
      });
    };
  }, []);

  const { ready, isWizard, menuCollapsed, helpCollapsed } = useAppContext();

  useAppHotkeys();
  useAppHealth();

  if (!ready) return <div>Not ready</div>;

  const header = { height: 60 };
  const footer = { height: 40 };
  const navbar = {
    width: getBarWidth(menuCollapsed, 1),
    breakpoint: 'sm',
    collapsed: { mobile: !menuCollapsed },
  };
  const aside = {
    width: getBarWidth(helpCollapsed, 2),
    breakpoint: 'sm',
    collapsed: { mobile: !helpCollapsed },
  };

  return (
    <Router>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          height: '100vh',
        }}
      >
        <AppShell
          layout="default"
          header={header}
          footer={footer}
          navbar={navbar}
          aside={aside}
        >
          <Header />
          <MenuBar disabled={isWizard} />
          <ViewContextProvider>
            <MainView />
          </ViewContextProvider>
          <HelpBar />
          <div
            style={{
              position: 'absolute',
              bottom: '40px',
              right: `${getBarWidth(helpCollapsed, 2)}px`,
              zIndex: 1000,
            }}
          >
            <NodeStatus />
          </div>
          <Footer />
        </AppShell>
      </div>
    </Router>
  );
};
