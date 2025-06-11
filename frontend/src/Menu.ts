import { SetInitialized } from '@app';
import { About, DalleDress, History, Home, Khedra, Settings } from '@views';
import {
  // ABIS_ROUTE
  Abis,
  // ABIS_ROUTE
  // MONITORS_ROUTE
  Monitors,
  // MONITORS_ROUTE
  // NAMES_ROUTE
  Names,
  // NAMES_ROUTE
} from '@views';
import { Wizard } from '@wizards';

export interface MenuItem {
  label: string;
  path: string;
  position: 'top' | 'bottom' | 'hidden';
  component?: React.ComponentType;
  hotkey?: string;
  altHotkey?: string;
  type?: 'navigation' | 'dev' | 'toggle';
  action?: () => void | Promise<void>;
}

export const MenuItems: MenuItem[] = [
  {
    label: 'Home',
    path: '/',
    position: 'top',
    component: Home,
    hotkey: 'mod+1',
    altHotkey: 'alt+1',
    type: 'navigation',
  },
  {
    label: 'About',
    path: '/about',
    position: 'top',
    component: About,
    hotkey: 'mod+2',
    altHotkey: 'alt+2',
    type: 'navigation',
  },
  {
    label: 'History',
    path: '/history/:address',
    position: 'top',
    component: History,
    hotkey: 'mod+3',
    altHotkey: 'alt+3',
    type: 'navigation',
  },
  {
    label: 'Khedra',
    path: '/khedra',
    position: 'top',
    component: Khedra,
    hotkey: 'mod+4',
    altHotkey: 'alt+4',
    type: 'navigation',
  },
  {
    label: 'DalleDress',
    path: '/dalledress',
    position: 'top',
    component: DalleDress,
    hotkey: 'mod+5',
    altHotkey: 'alt+5',
    type: 'navigation',
  },
  // ABIS_ROUTE
  {
    label: 'ABIs',
    path: '/abis',
    position: 'top',
    component: Abis,
    hotkey: 'mod+6',
    altHotkey: 'alt+6',
    type: 'navigation',
  },
  // ABIS_ROUTE
  // MONITORS_ROUTE
  {
    label: 'Monitors',
    path: '/monitors',
    position: 'top',
    component: Monitors,
    hotkey: 'mod+7',
    altHotkey: 'alt+7',
    type: 'navigation',
  },
  // MONITORS_ROUTE
  // NAMES_ROUTE
  {
    label: 'Names',
    path: '/names',
    position: 'top',
    component: Names,
    hotkey: 'mod+8',
    altHotkey: 'alt+8',
    type: 'navigation',
  },
  // NAMES_ROUTE
  {
    label: 'Settings',
    path: '/settings',
    position: 'bottom',
    component: Settings,
    hotkey: 'mod+9',
    altHotkey: 'alt+9',
    type: 'navigation',
  },
  {
    path: '/wizard',
    label: 'Wizard',
    position: 'hidden',
    component: Wizard,
    hotkey: 'mod+shift+w',
    type: 'dev',
    action: async () => {
      await SetInitialized(false);
    },
  },
];
