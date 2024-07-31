import React from "react";

// Find: NewViews
import { IconHome, IconSpider, IconList, IconArticle, IconTag, IconDeviceDesktop, IconServer, IconSettings } from "@tabler/icons-react";
import { HomeView, DalleView, SeriesView, HistoryView, NamesView, MonitorsView, ServersView, SettingsView } from "@views";

// Note:
//  Change with care. The order of the items in this list matters (the last one is the default).
//  The order field is used to sort the menu items.
export const routeItems = [
  {
    order: 2,
    route: "/dalle",
    label: "Dalle",
    icon: <IconSpider />,
    component: DalleView,
  },
  {
    order: 3,
    route: "/series",
    label: "Series",
    icon: <IconList />,
    component: SeriesView,
  },
  {
    order: 10,
    route: "/history/:address",
    label: "History",
    icon: <IconArticle />,
    component: HistoryView,
  },
  {
    order: 20,
    route: "/monitors",
    label: "Monitors",
    icon: <IconDeviceDesktop />,
    component: MonitorsView,
  },
  {
    order: 30,
    route: "/names",
    label: "Names",
    icon: <IconTag />,
    component: NamesView,
  },
  {
    order: 40,
    route: "/servers",
    label: "Servers",
    icon: <IconServer />,
    component: ServersView,
  },
  {
    order: 50,
    route: "/settings",
    label: "Settings",
    icon: <IconSettings />,
    component: SettingsView,
  },
  {
    order: 0,
    route: "/",
    label: "Home",
    icon: <IconHome />,
    component: HomeView,
  },
];
