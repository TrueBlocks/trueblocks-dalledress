import { FC } from 'react';

import { IconType } from 'react-icons';
import {
  FaBroom,
  FaCloudUploadAlt,
  FaMagic,
  FaPencilAlt,
  FaPlus,
  FaTimesCircle,
  FaTrash,
  FaTrashRestore,
} from 'react-icons/fa';

export interface ModifyIconProps {
  type:
    | 'add'
    | 'edit'
    | 'delete'
    | 'undelete'
    | 'remove'
    | 'autoname'
    | 'clean'
    | 'publish';
  onClick: () => void;
  disabled?: boolean;
  title?: string;
  size?: number;
}

export const ModifyIcon: FC<ModifyIconProps> = ({
  type,
  onClick,
  disabled = false,
  title,
  size = 16,
}) => {
  // Map of action types to Font Awesome icons
  const iconMap: Record<string, IconType> = {
    add: FaPlus,
    edit: FaPencilAlt,
    delete: FaTrash,
    undelete: FaTrashRestore,
    remove: FaTimesCircle,
    autoname: FaMagic,
    clean: FaBroom,
    publish: FaCloudUploadAlt,
  };

  const Icon = iconMap[type];

  if (!Icon) return null;

  // Default title if not provided
  const tooltipTitle =
    title || `${type.charAt(0).toUpperCase() + type.slice(1)}`;

  return (
    <span
      onClick={disabled ? undefined : onClick}
      title={tooltipTitle}
      style={{
        cursor: disabled ? 'not-allowed' : 'pointer',
        opacity: disabled ? 0.5 : 1,
        margin: '0 4px',
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '4px',
        borderRadius: '4px',
        transition: 'background-color 0.2s ease',
      }}
      className={`modify-icon modify-icon-${type}`}
    >
      <Icon size={size} />
    </span>
  );
};
