import React from "react";

// Find: Routes
import { IconHome, IconSpider, IconList, IconArticle, IconTag, IconServer, IconSettings } from "@tabler/icons-react";
import { HomeView, DalleView, SeriesView, HistoryView, NamesView, ServersView, SettingsView } from "@views";

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
    order: 4,
    route: "/history",
    label: "History",
    icon: <IconArticle />,
    component: HistoryView,
  },
  {
    order: 5,
    route: "/names",
    label: "Names",
    icon: <IconTag />,
    component: NamesView,
  },
  {
    order: 6,
    route: "/servers",
    label: "Servers",
    icon: <IconServer />,
    component: ServersView,
  },
  {
    order: 7,
    route: "/settings",
    label: "Settings",
    icon: <IconSettings />,
    component: SettingsView,
  },
  {
    order: 1,
    route: "/",
    label: "Home",
    icon: <IconHome />,
    component: HomeView,
  },
];
