import { SetInitialized } from '@app';
import {
  Abis,
  About,
  DalleDress,
  History,
  Home,
  Khedra,
  Names,
  Settings,
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
    label: 'Names',
    path: '/names',
    position: 'top',
    component: Names,
    hotkey: 'mod+5',
    altHotkey: 'alt+5',
    type: 'navigation',
  },
  {
    label: 'DalleDress',
    path: '/dalledress',
    position: 'top',
    component: DalleDress,
    hotkey: 'mod+6',
    altHotkey: 'alt+6',
    type: 'navigation',
  },
  {
    label: 'ABIs',
    path: '/abis',
    position: 'top',
    component: Abis,
    hotkey: 'mod+7',
    altHotkey: 'alt+7',
    type: 'navigation',
  },
  {
    label: 'Settings',
    path: '/settings',
    position: 'bottom',
    component: Settings,
    hotkey: 'mod+8',
    altHotkey: 'alt+8',
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
