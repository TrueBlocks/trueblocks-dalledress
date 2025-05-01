import { SetInitialized } from '@app';
import { About, Data, Home, Names, Settings } from '@views';
import { Wizard } from '@wizards';
import {
  FaCog,
  FaDatabase,
  FaHome,
  FaInfoCircle,
  FaUser,
} from 'react-icons/fa';

export interface MenuItem {
  icon: React.ComponentType<{ size: number; style?: React.CSSProperties }>;
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
    icon: FaHome,
    label: 'Home',
    path: '/',
    position: 'top',
    component: Home,
    hotkey: 'mod+1',
    altHotkey: 'alt+1',
    type: 'navigation',
  },
  {
    icon: FaInfoCircle,
    label: 'About',
    path: '/about',
    position: 'top',
    component: About,
    hotkey: 'mod+2',
    altHotkey: 'alt+2',
    type: 'navigation',
  },
  {
    icon: FaDatabase,
    label: 'Data',
    path: '/data',
    position: 'top',
    component: Data,
    hotkey: 'mod+3',
    altHotkey: 'alt+3',
    type: 'navigation',
  },
  {
    icon: FaUser,
    label: 'Names',
    path: '/names',
    position: 'top',
    component: Names,
    hotkey: 'mod+4',
    altHotkey: 'alt+4',
    type: 'navigation',
  },
  {
    icon: FaCog,
    label: 'Settings',
    path: '/settings',
    position: 'bottom',
    component: Settings,
    hotkey: 'mod+5',
    altHotkey: 'alt+5',
    type: 'navigation',
  },
  {
    icon: FaCog,
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
