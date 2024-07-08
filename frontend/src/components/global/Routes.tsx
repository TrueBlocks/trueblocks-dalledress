import React, { useEffect } from 'react';
import { Route, Switch, useLocation } from 'wouter';
import classes from './Routes.module.css';
import DalleView from '@/views/Dalle/DalleView';
import NamesView from '@/views/Names/NamesView';
import StatusBar from './StatusBar';
import HomeView from '@/views/Home/HomeView';
import SettingsView from '@/views/Settings/SettingsView';
import { GetLastRoute } from '@gocode/app/App';

export const Routes = () => {
  const [, setLocation] = useLocation();

  useEffect(() => {
    const lastRoute = (GetLastRoute() || '/').then((route) => {
      setLocation(route);
    });
  }, [setLocation]);

  return (
    <div className={classes.container}>
      <Switch>
        <Route path="/dalle">
          <DalleView />
        </Route>
        <Route path="/names">
          <NamesView />
        </Route>
        <Route path="/settings">
          <SettingsView />
        </Route>
        <Route path="/">
          <HomeView />
        </Route>
      </Switch>
      <StatusBar />
    </div>
  );
};
