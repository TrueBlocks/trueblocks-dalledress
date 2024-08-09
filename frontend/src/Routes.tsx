import React from "react";

// Find: NewViews
import {
  // Icons
  HomeIcon,
  HistoryIcon,
  MonitorIcon,
  NamesIcon,
  IndexesIcon,
  ManifestIcon,
  AbisIcon,
  StatusIcon,
  DaemonsIcon,
  SettingsIcon,
  SeriesIcon,
  DalleIcon,
  // Views
  HomeView,
  HistoryView,
  NamesView,
  MonitorsView,
  IndexesView,
  ManifestView,
  AbisView,
  DaemonsView,
  StatusView,
  SettingsView,
  SeriesView,
  DalleView,
} from "@views";

// Note:
//  Change with care. The order of the items in this list matters (the last one is the default).
//  The order field is used to sort the menu items.
export const routeItems = [
  {
    order: 2,
    route: "/dalle",
    label: "Dalle",
    icon: DalleIcon,
    component: DalleView,
  },
  {
    order: 3,
    route: "/series",
    label: "Series",
    icon: SeriesIcon,
    component: SeriesView,
  },
  {
    order: 10,
    route: "/history/:address",
    label: "History",
    icon: HistoryIcon,
    component: HistoryView,
  },
  {
    order: 20,
    route: "/monitors",
    label: "Monitors",
    icon: MonitorIcon,
    component: MonitorsView,
  },
  {
    order: 30,
    route: "/names",
    label: "Names",
    icon: NamesIcon,
    component: NamesView,
  },
  {
    order: 40,
    route: "/indexes",
    label: "Indexes",
    icon: IndexesIcon,
    component: IndexesView,
  },
  {
    order: 50,
    route: "/manifest",
    label: "Manifest",
    icon: ManifestIcon,
    component: ManifestView,
  },
  {
    order: 60,
    route: "/abis",
    label: "Abis",
    icon: AbisIcon,
    component: AbisView,
  },
  {
    order: 70,
    route: "/status",
    label: "Status",
    icon: StatusIcon,
    component: StatusView,
  },
  {
    order: 80,
    route: "/daemons",
    label: "Daemons",
    icon: DaemonsIcon,
    component: DaemonsView,
  },
  {
    order: 90,
    route: "/settings",
    label: "Settings",
    icon: SettingsIcon,
    component: SettingsView,
  },
  {
    order: 0,
    route: "/",
    label: "Home",
    icon: HomeIcon,
    component: HomeView,
  },
];
