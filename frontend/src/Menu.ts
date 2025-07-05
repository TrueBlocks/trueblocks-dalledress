import { SetInitialized } from '@app';
import { DalleDress, Home, Khedra, Settings } from '@views';
import { Abis, Chunks, Exports, Monitors, Names, Status } from '@views';
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
    label: 'Exports',
    path: '/exports',
    position: 'top',
    component: Exports,
    hotkey: 'mod+2',
    altHotkey: 'alt+2',
    type: 'navigation',
  },
  {
    label: 'Monitors',
    path: '/monitors',
    position: 'top',
    component: Monitors,
    hotkey: 'mod+3',
    altHotkey: 'alt+3',
    type: 'navigation',
  },
  {
    label: 'Abis',
    path: '/abis',
    position: 'top',
    component: Abis,
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
    label: 'Chunks',
    path: '/chunks',
    position: 'top',
    component: Chunks,
    hotkey: 'mod+6',
    altHotkey: 'alt+6',
    type: 'navigation',
  },
  {
    label: 'Status',
    path: '/status',
    position: 'top',
    component: Status,
    hotkey: 'mod+7',
    altHotkey: 'alt+7',
    type: 'navigation',
  },
  {
    label: 'DalleDress',
    path: '/dalledress',
    position: 'top',
    component: DalleDress,
    hotkey: 'mod+8',
    altHotkey: 'alt+8',
    type: 'navigation',
  },
  {
    label: 'Khedra',
    path: '/khedra',
    position: 'bottom',
    component: Khedra,
    hotkey: 'mod+9',
    altHotkey: 'alt+9',
    type: 'navigation',
  },
  {
    label: 'Settings',
    path: '/settings',
    position: 'bottom',
    component: Settings,
    hotkey: 'mod+0',
    altHotkey: 'alt+0',
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
