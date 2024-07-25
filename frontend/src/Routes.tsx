import React from "react";

// Find: Routes
import { IconHome, IconTag, IconSettings, } from "@tabler/icons-react";
import { HomeView, NamesView, SettingsView } from "@views";

// Note:
//  Change with care. The order of the items in this list matters (the last one is the default).
//  The order field is used to sort the menu items.
export const routeItems = [
  {
    order: 5,
    route: "/names",
    label: "Names",
    icon: <IconTag />,
    component: NamesView,
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
