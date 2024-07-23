import React, { useEffect } from "react";
import { Route, Switch, useLocation } from "wouter";
import classes from "@/App.module.css";
import { HomeView, NamesView, HistoryView, DalleView, SeriesView, SettingsView } from "@views";
import { GetLast } from "@gocode/app/App";

export const Routes = () => {
  const [, setLocation] = useLocation();

  useEffect(() => {
    const lastRoute = (GetLast("route") || "/").then((route) => {
      setLocation(route);
    });
  }, [setLocation]);

  var menuItems = [
    { route: "/dalle", component: DalleView },
    { route: "/series", component: SeriesView },
    { route: "/history", component: HistoryView },
    { route: "/names", component: NamesView },
    { route: "/settings", component: SettingsView },
    { route: "/", component: HomeView },
  ];

  return (
    <div className={classes.mainContent}>
      <Switch>
        {menuItems.map((item) => (
          <Route key={item.route} path={item.route}>
            <item.component />
          </Route>
        ))}
      </Switch>
    </div>
  );
};
