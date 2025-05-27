import { FC, createElement } from 'react';

import { getIconSet } from './Icons';

const DEFAULT_SIZE = 16;
const DEFAULT_ICON_SET = 'fa';

// Define types for icon components - use proper SVG element props
type IconProps = React.SVGAttributes<SVGElement> & {
  style?: React.CSSProperties;
  size?: number;
};

/**
 * Hook that provides access to all icons with consistent size
 */
export const useIcons = (
  iconSetName = DEFAULT_ICON_SET,
  size = DEFAULT_SIZE,
) => {
  const iconSet = getIconSet(iconSetName);

  // Create explicit named components for each icon to preserve component names in React DevTools
  // This properly types the components as React components that can be used in JSX
  const Home: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Home, { size, ...props });
  const About: FC<IconProps> = (props = {}) =>
    createElement(iconSet.About, { size, ...props });
  const History: FC<IconProps> = (props = {}) =>
    createElement(iconSet.History, { size, ...props });
  const Khedra: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Khedra, { size, ...props });
  const DalleDress: FC<IconProps> = (props = {}) =>
    createElement(iconSet.DalleDress, { size, ...props });
  const Settings: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Settings, { size, ...props });
  const Wizard: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Wizard, { size, ...props });

  // ABIS_CODE
  const Names: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Names, { size, ...props });
  const ABIs: FC<IconProps> = (props = {}) =>
    createElement(iconSet.ABIs, { size, ...props });

  const Switch: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Switch, { size, ...props });
  const File: FC<IconProps> = (props = {}) =>
    createElement(iconSet.File, { size, ...props });

  const Twitter: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Twitter, { size, ...props });
  const Github: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Github, { size, ...props });
  const Website: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Website, { size, ...props });
  const Email: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Email, { size, ...props });

  const Add: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Add, { size, ...props });
  const Edit: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Edit, { size, ...props });
  const Delete: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Delete, { size, ...props });
  const Undelete: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Undelete, { size, ...props });
  const Remove: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Remove, { size, ...props });
  const Autoname: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Autoname, { size, ...props });

  const ChevronLeft: FC<IconProps> = (props = {}) =>
    createElement(iconSet.ChevronLeft, { size, ...props });
  const ChevronRight: FC<IconProps> = (props = {}) =>
    createElement(iconSet.ChevronRight, { size, ...props });
  const ChevronUp: FC<IconProps> = (props = {}) =>
    createElement(iconSet.ChevronUp, { size, ...props });
  const ChevronDown: FC<IconProps> = (props = {}) =>
    createElement(iconSet.ChevronDown, { size, ...props });
  const Light: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Light, { size, ...props });
  const Dark: FC<IconProps> = (props = {}) =>
    createElement(iconSet.Dark, { size, ...props });

  // Return all the components
  return {
    Home,
    About,
    History,
    Khedra,
    DalleDress,
    Settings,
    Wizard,

    // ABIS_CODE
    Names,
    ABIs,

    Switch,
    File,

    Twitter,
    Github,
    Website,
    Email,

    Add,
    Edit,
    Delete,
    Undelete,
    Remove,
    Autoname,

    ChevronLeft,
    ChevronRight,
    ChevronUp,
    ChevronDown,
    Light,
    Dark,
  };
};
