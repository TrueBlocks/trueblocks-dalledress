import React from "react";
import { Route, Switch } from "wouter";
import classes from "./Routes.module.css";
import DalleView from "@/views/Dalle/DalleView";
import HomeView from "@/views/Home/HomeView";
import NamesView from "@/views/Names/NamesView";
import SettingsView from "@/views/Settings/SettingsView";

function Routes() {
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
        <Route>
          <HomeView />
        </Route>
      </Switch>
    </div>
  );
}

export default Routes;
