// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
import { FC, createElement, useMemo } from 'react';

import { IconType } from 'react-icons';

import * as Icons from './Icons';

export type IconSet = {
  Exports: IconType;
  Monitors: IconType;
  Abis: IconType;
  Names: IconType;
  Chunks: IconType;
  Contracts: IconType;
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
  Publish: IconType;
  Pin: IconType;

  ChevronLeft: IconType;
  ChevronRight: IconType;
  ChevronUp: IconType;
  ChevronDown: IconType;

  Light: IconType;
  Dark: IconType;

  DebugOn: IconType;
  DebugOff: IconType;

  Missing: IconType;
};

const faIcons: IconSet = {
  // Collections
  Exports: Icons.FaExports,
  Monitors: Icons.FaMonitors,
  Abis: Icons.FaAbis,
  Names: Icons.FaNames,
  Chunks: Icons.FaChunks,
  Contracts: Icons.FaContracts,
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
  Publish: Icons.FaGlobe,
  Pin: Icons.FaListAlt,

  // Navigation
  ChevronLeft: Icons.FaAngleDoubleLeft,
  ChevronRight: Icons.FaAngleDoubleRight,
  ChevronUp: Icons.FaAngleDoubleUp,
  ChevronDown: Icons.FaAngleDoubleDown,

  // Theme
  Light: Icons.FaSun,
  Dark: Icons.FaMoon,

  // Debug
  DebugOn: Icons.FaEye,
  DebugOff: Icons.FaEyeSlash,

  // Fallback
  Missing: Icons.FaMissing,
};

const biIcons: IconSet = {
  // Collections
  Exports: Icons.BiExports,
  Monitors: Icons.BiMonitors,
  Abis: Icons.BiAbis,
  Names: Icons.BiNames,
  Chunks: Icons.BiChunks,
  Contracts: Icons.BiContracts,
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
  Publish: Icons.BiGlobe,
  Pin: Icons.BiListUl,

  // Navigation
  ChevronLeft: Icons.BiChevronsLeft,
  ChevronRight: Icons.BiChevronsRight,
  ChevronUp: Icons.BiChevronsUp,
  ChevronDown: Icons.BiChevronsDown,

  // Theme
  Light: Icons.BiSun,
  Dark: Icons.BiMoon,

  // Debug
  DebugOn: Icons.BiShow,
  DebugOff: Icons.BiHide,

  // Fallback
  Missing: Icons.BiMissing,
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
    const Exports: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Exports, { size, ...props });
    const Monitors: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Monitors, { size, ...props });
    const Abis: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Abis, { size, ...props });
    const Names: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Names, { size, ...props });
    const Chunks: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Chunks, { size, ...props });
    const Contracts: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Contracts, { size, ...props });
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
    const Publish: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Publish, { size, ...props });
    const Pin: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Pin, { size, ...props });

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
    const DebugOn: FC<IconProps> = (props = {}) =>
      createElement(iconSet.DebugOn, { size, ...props });
    const DebugOff: FC<IconProps> = (props = {}) =>
      createElement(iconSet.DebugOff, { size, ...props });
    const Missing: FC<IconProps> = (props = {}) =>
      createElement(iconSet.Missing, { size, ...props });

    return {
      Exports,
      Monitors,
      Abis,
      Names,
      Chunks,
      Contracts,
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
      Publish,
      Pin,

      ChevronLeft,
      ChevronRight,
      ChevronUp,
      ChevronDown,
      Light,
      Dark,
      DebugOn,
      DebugOff,
      Missing,
    };
  }, [iconSet, size]);
};
