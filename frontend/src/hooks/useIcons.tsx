import React from 'react';

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
  const Home: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Home, { size, ...props });
  const About: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.About, { size, ...props });
  const History: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.History, { size, ...props });
  const Khedra: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Khedra, { size, ...props });
  const Names: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Names, { size, ...props });
  const DalleDress: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.DalleDress, { size, ...props });
  const Settings: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Settings, { size, ...props });
  const Wizard: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Wizard, { size, ...props });
  const Switch: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Switch, { size, ...props });
  const File: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.File, { size, ...props });
  const Twitter: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Twitter, { size, ...props });
  const Github: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Github, { size, ...props });
  const Website: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Website, { size, ...props });
  const Email: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Email, { size, ...props });
  const Add: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Add, { size, ...props });
  const Edit: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Edit, { size, ...props });
  const Delete: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Delete, { size, ...props });
  const Undelete: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Undelete, { size, ...props });
  const Remove: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Remove, { size, ...props });
  const ChevronLeft: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.ChevronLeft, { size, ...props });
  const ChevronRight: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.ChevronRight, { size, ...props });
  const ChevronUp: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.ChevronUp, { size, ...props });
  const ChevronDown: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.ChevronDown, { size, ...props });
  const Light: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Light, { size, ...props });
  const Dark: React.FC<IconProps> = (props = {}) =>
    React.createElement(iconSet.Dark, { size, ...props });

  // Return all the components
  return {
    Home,
    About,
    History,
    Khedra,
    Names,
    DalleDress,
    Settings,
    Wizard,
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
    ChevronLeft,
    ChevronRight,
    ChevronUp,
    ChevronDown,
    Light,
    Dark,
  };
};
