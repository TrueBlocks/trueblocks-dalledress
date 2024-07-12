import React, { useEffect } from "react";
import { Route, Switch, useLocation } from "wouter";
import classes from "@/App.module.css";
import { DalleView, NamesView, HomeView, SeriesView, SettingsView } from "@/views";
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
    { route: "/names", component: NamesView },
    { route: "/series", component: SeriesView },
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

export default Routes;
