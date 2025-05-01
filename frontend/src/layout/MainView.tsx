import React, { useEffect, useRef } from 'react';

import { useAppContext } from '@contexts';
import { Breadcrumb } from '@layout';
import { AppShell } from '@mantine/core';
import { MenuItem, MenuItems } from 'src/Menu';
import { Route } from 'wouter';

import { StatusBar } from './StatusBar';

function isComponentMenuItem(
  item: MenuItem,
): item is MenuItem & { component: React.ComponentType } {
  return !!item.component;
}

export const MainView = () => {
  const { currentLocation } = useAppContext();
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollTop = 0;
    }
  }, [currentLocation]);

  return (
    <AppShell.Main
      style={{
        width: '100%',
        height: '100vh',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
      }}
    >
      <div>
        <Breadcrumb />
      </div>

      <div
        ref={scrollContainerRef}
        style={{
          width: '100%',
          flex: 1,
          overflowY: 'auto',
          overflowX: 'hidden',
          position: 'relative',
          backgroundColor: 'black',
        }}
      >
        {MenuItems.filter(isComponentMenuItem).map((item) => (
          <Route key={item.path} path={item.path}>
            {React.createElement(item.component)}
          </Route>
        ))}
      </div>

      <div
        style={{
          width: '100%',
          backgroundColor: 'white',
        }}
      >
        <StatusBar />
      </div>
    </AppShell.Main>
  );
};
