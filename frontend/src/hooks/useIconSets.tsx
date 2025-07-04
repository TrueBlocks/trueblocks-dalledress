import { FC, createElement, useMemo } from 'react';

import { IconType } from 'react-icons';

import * as Icons from './Icons';

export type IconSet = {
  Monitors: IconType;
  Names: IconType;
  Chunks: IconType;
  Exports: IconType;
  Abis: IconType;
  Status: IconType;

  Home: IconType;
  Khedra: IconType;
  DalleDress: IconType;
  Settings: IconType;
  Wizard: IconType;

  Switch: IconType;
  File: IconType;

  Twitter: IconType;
  Github: IconType;
  Website: IconType;
  Email: IconType;

  Add: IconType;
  Edit: IconType;
  Delete: IconType;
  Undelete: IconType;
  Remove: IconType;
  Autoname: IconType;

  ChevronLeft: IconType;
  ChevronRight: IconType;
  ChevronUp: IconType;
  ChevronDown: IconType;

  Light: IconType;
  Dark: IconType;
};

const faIcons: IconSet = {
  // Collections
  Monitors: Icons.FaMonitors,
  Names: Icons.FaNames,
  Chunks: Icons.FaChunks,
  Exports: Icons.FaExports,
  Abis: Icons.FaAbis,
  Status: Icons.FaStatus,

  // App navigation
  Home: Icons.FaHome,
  Khedra: Icons.FaIndustry,
  DalleDress: Icons.FaPalette,
  Settings: Icons.FaCog,
  Wizard: Icons.FaHatWizard,

  // General
  Switch: Icons.FaRandom,
  File: Icons.FaFile,

  // Social
  Twitter: Icons.FaTwitter,
  Github: Icons.FaGithub,
  Website: Icons.FaGlobe,
  Email: Icons.FaEnvelope,

  // Actions
  Add: Icons.FaPlus,
  Edit: Icons.FaEdit,
  Delete: Icons.FaTimes,
  Undelete: Icons.FaUndo,
  Remove: Icons.FaEraser,
  Autoname: Icons.FaMagic,

  // Navigation
  ChevronLeft: Icons.FaAngleDoubleLeft,
  ChevronRight: Icons.FaAngleDoubleRight,
  ChevronUp: Icons.FaAngleDoubleUp,
  ChevronDown: Icons.FaAngleDoubleDown,

  // Theme
  Light: Icons.FaSun,
  Dark: Icons.FaMoon,
};

const biIcons: IconSet = {
  // Collections
  Monitors: Icons.BiMonitors,
  Names: Icons.BiNames,
  Chunks: Icons.BiChunks,
  Exports: Icons.BiExports,
  Abis: Icons.BiAbis,
  Status: Icons.BiStatus,

  // App navigation
  Home: Icons.BiHome,
  Khedra: Icons.BiBuildings,
  DalleDress: Icons.BiPalette,
  Settings: Icons.BiCog,
  Wizard: Icons.BiCog,

  // General
  Switch: Icons.BiTransfer,
  File: Icons.BiFile,

  // Social
  Twitter: Icons.BiLogoTwitter,
  Github: Icons.BiLogoGithub,
  Website: Icons.BiGlobe,
  Email: Icons.BiEnvelope,

  // Actions
  Add: Icons.BiPlus,
  Edit: Icons.BiPencil,
  Delete: Icons.BiX,
  Undelete: Icons.BiUndo,
  Remove: Icons.BiTrash,
  Autoname: Icons.BiBot,

  // Navigation
  ChevronLeft: Icons.BiChevronsLeft,
  ChevronRight: Icons.BiChevronsRight,
  ChevronUp: Icons.BiChevronsUp,
  ChevronDown: Icons.BiChevronsDown,

  // Theme
  Light: Icons.BiSun,
  Dark: Icons.BiMoon,
};

const iconSets: Record<string, IconSet> = {
  fa: faIcons,
  bi: biIcons,
};

export const getIconSet = (name: string): IconSet => {
  const iconSet = iconSets[name];
  if (!iconSet) {
    throw new Error(`Icon set "${name}" not found`);
  }
  return iconSet;
};

const DEFAULT_SIZE = 16;
const DEFAULT_ICON_SET = 'fa';

// Define types for icon components
type IconProps = React.SVGAttributes<SVGElement> & {
  style?: React.CSSProperties;
  size?: number;
};

export const useIconSets = (
  iconSetName = DEFAULT_ICON_SET,
  size = DEFAULT_SIZE,
) => {
  const iconSet = getIconSet(iconSetName);

  return useMemo(() => {
    const Monitors: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Monitors, { size, ...props });
    const Names: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Names, { size, ...props });
    const Chunks: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Chunks, { size, ...props });
    const Exports: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Exports, { size, ...props });
    const Abis: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Abis, { size, ...props });
    const Status: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Status, { size, ...props });

    const Home: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Home, { size, ...props });
    const Khedra: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Khedra, { size, ...props });
    const DalleDress: FC<IconProps> = (props = {}) =>
      createElement(iconSet.DalleDress, { size, ...props });
    const Settings: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Settings, { size, ...props });
    const Wizard: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Wizard, { size, ...props });

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

    return {
      Monitors,
      Names,
      Chunks,
      Exports,
      Abis,
      Status,

      Home,
      Khedra,
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
      Autoname,

      ChevronLeft,
      ChevronRight,
      ChevronUp,
      ChevronDown,
      Light,
      Dark,
    };
  }, [iconSet, size]);
};
