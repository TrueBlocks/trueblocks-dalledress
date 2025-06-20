/*
List of other icons sets
Ai: Ant Design Icons
Bs: Bootstrap Icons
Cg: css.gg
Ci: CoreUI Icons
Di: Devicon
Fi: Feather Icons
Fc: Flat Color Icons
Gi: Game Icons
Go: GitHub Octicons
Gr: Grommet Icons
Hi, Hi2: Heroicons
Im: IcoMoon Free
Io, Io5: Ionicons (4 and 5)
Lu: Lucide Icons
Md: Material Design Icons
Pi: Phosphor Icons
Ri: Remix Icons
Rx: Radix Icons
Si: Simple Icons
Sl: Simple Line Icons
Tb: Tabler Icons
Tfi: Themify Icons
Ti: Typicons
Vsc: VS Code Icons
Wi: Weather Icons
Fa6: Font Awesome 6
*/
// A type called IconSet that has Home, Settings, and Khedra keys, all strings
import { IconType } from 'react-icons';
// Import BoxIcons (bi)
import {
  BiBot,
  BiBuildings,
  BiChevronsDown,
  BiChevronsLeft,
  BiChevronsRight,
  BiChevronsUp,
  BiCog,
  BiDesktop,
  BiEnvelope,
  BiFile,
  BiGlobe,
  BiHistory,
  BiHome,
  BiListUl,
  BiLogoGithub,
  BiLogoTwitter,
  BiMoon,
  BiPalette,
  BiPencil,
  BiPlus,
  BiSun,
  BiTransfer,
  BiTrash,
  BiUndo,
  BiUser,
  BiX,
} from 'react-icons/bi';
// A function that returns an IconSet using icons from react-icons/fa
import {
  FaAngleDoubleDown,
  FaAngleDoubleLeft,
  FaAngleDoubleRight,
  FaAngleDoubleUp,
  FaCog,
  FaDesktop,
  FaEdit,
  FaEnvelope,
  FaEraser,
  FaFile,
  FaGithub,
  FaGlobe,
  FaHatWizard,
  FaHistory,
  FaHome,
  FaIndustry,
  FaListAlt,
  FaMagic,
  FaMoon,
  FaPalette,
  FaPlus,
  FaRandom,
  FaSun,
  FaTimes,
  FaTwitter,
  FaUndo,
  FaUser,
} from 'react-icons/fa';

// An IconSet using icons from react-icons/fa to represent Home, Settings, and Khedra
export type IconSet = {
  // Menu options
  Home: IconType;
  Khedra: IconType;
  DalleDress: IconType;
  Settings: IconType;
  Wizard: IconType;

  // ADD_ROUTE
  Names: IconType;
  ABIs: IconType;
  Monitors: IconType;
  Chunks: IconType;
  Exports: IconType;
  // ADD_ROUTE

  // File operations
  Switch: IconType;
  File: IconType;

  // Social media icons
  Twitter: IconType;
  Github: IconType;
  Website: IconType;
  Email: IconType;

  // Editing operations
  Add: IconType;
  Edit: IconType;
  Delete: IconType;
  Undelete: IconType;
  Remove: IconType;
  Autoname: IconType;

  // Chevrons
  ChevronLeft: IconType;
  ChevronRight: IconType;
  ChevronUp: IconType;
  ChevronDown: IconType;

  // Color modes
  Light: IconType;
  Dark: IconType;
};

// An IconSet using icons from fa react to represent Home, Settings, and Khedra
const faIcons: IconSet = {
  Home: FaHome,
  Khedra: FaIndustry,
  DalleDress: FaPalette,
  Settings: FaCog,
  Wizard: FaHatWizard,

  // ADD_ROUTE
  Names: FaUser,
  ABIs: FaListAlt,
  Monitors: FaDesktop,
  Chunks: FaIndustry,
  Exports: FaHistory,
  // ADD_ROUTE

  Switch: FaRandom,
  File: FaFile,

  Twitter: FaTwitter,
  Github: FaGithub,
  Website: FaGlobe,
  Email: FaEnvelope,

  Add: FaPlus,
  Edit: FaEdit,
  Delete: FaTimes,
  Undelete: FaUndo,
  Remove: FaEraser,
  Autoname: FaMagic,

  ChevronLeft: FaAngleDoubleLeft,
  ChevronRight: FaAngleDoubleRight,
  ChevronUp: FaAngleDoubleUp,
  ChevronDown: FaAngleDoubleDown,
  Light: FaSun,
  Dark: FaMoon,
};

// An IconSet using icons from bi react to represent Home, Settings, and Khedra
const biIcons: IconSet = {
  Home: BiHome,
  Khedra: BiBuildings,
  DalleDress: BiPalette,
  Settings: BiCog,
  Wizard: BiCog,

  // ADD_ROUTE
  Names: BiUser,
  ABIs: BiListUl,
  Monitors: BiDesktop,
  Chunks: BiBuildings,
  Exports: BiHistory,
  // ADD_ROUTE

  Switch: BiTransfer,
  File: BiFile,

  Twitter: BiLogoTwitter,
  Github: BiLogoGithub,
  Website: BiGlobe,
  Email: BiEnvelope,

  Add: BiPlus,
  Edit: BiPencil,
  Delete: BiX,
  Undelete: BiUndo,
  Remove: BiTrash,
  Autoname: BiBot,

  ChevronLeft: BiChevronsLeft,
  ChevronRight: BiChevronsRight,
  ChevronUp: BiChevronsUp,
  ChevronDown: BiChevronsDown,
  Light: BiSun,
  Dark: BiMoon,
};

// A map keyed by the name of the icon set, with values being the icon set itself
const iconSets: Record<string, IconSet> = {
  fa: faIcons,
  bi: biIcons,
};

// An exported routine that, given the name of an icon set, returns the icon set from the map
export const getIconSet = (name: string): IconSet => {
  const iconSet = iconSets[name];
  if (!iconSet) {
    throw new Error(`Icon set "${name}" not found`);
  }
  return iconSet;
};
