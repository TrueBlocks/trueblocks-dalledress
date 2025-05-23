// A type called IconSet that has Home, Settings, and Khedra keys, all strings
import { IconType } from 'react-icons';
// Import BoxIcons (bi)
import {
  BiBuildings,
  BiChevronsDown,
  BiChevronsLeft,
  BiChevronsRight,
  BiChevronsUp,
  BiCog,
  BiData,
  BiEnvelope,
  BiFile,
  BiGlobe,
  BiHome,
  BiInfoCircle,
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
  FaDatabase,
  FaEdit,
  FaEnvelope,
  FaEraser,
  FaFile,
  FaGithub,
  FaGlobe,
  FaHatWizard,
  FaHome,
  FaIndustry,
  FaMoon,
  FaPalette,
  FaPlus,
  FaQuestion,
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
  About: IconType;
  History: IconType;
  Khedra: IconType;
  Names: IconType;
  DalleDress: IconType;
  Settings: IconType;
  Wizard: IconType;

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
  About: FaQuestion,
  History: FaDatabase,
  Khedra: FaIndustry,
  Names: FaUser,
  DalleDress: FaPalette,
  Settings: FaCog,
  Wizard: FaHatWizard,
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
  About: BiInfoCircle,
  History: BiData,
  Khedra: BiBuildings,
  Names: BiUser,
  DalleDress: BiPalette,
  Settings: BiCog,
  Wizard: BiCog,
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
