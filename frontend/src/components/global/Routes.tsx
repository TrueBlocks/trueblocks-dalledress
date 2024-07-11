import React from "react";
import { Route, Switch } from "wouter";
import classes from "@/App.module.css";
import HomeView from "@/views/Home/HomeView";
import SettingsView from "@/views/Settings/SettingsView";

function Routes() {
  return (
    <div className={classes.mainContent}>
      <Switch>
        <Route path="/settings">
          <SettingsView />
        </Route>
        <Route>
          <HomeView />
        </Route>
      </Switch>
    </div>
  );
}

export default Routes;
