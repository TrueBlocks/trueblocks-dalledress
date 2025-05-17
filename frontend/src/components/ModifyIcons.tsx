import { FC } from 'react';

import { base, types } from '@models';

import { ModifyIcon } from './ModifyIcon';

export interface ModifyIconsProps {
  row: types.Name;
  onAction: (type: string, row: types.Name, tabName?: string) => void;
  tabName?: string;
}

export const ModifyIcons: FC<ModifyIconsProps> = ({
  row,
  onAction,
  tabName,
}) => {
  const isDeleted = row.deleted === true;

  // Correctly convert base.Address to a hex string
  const typedAddress = row.address as base.Address;
  let addressString = '';
  if (
    typedAddress &&
    typedAddress.address &&
    Array.isArray(typedAddress.address)
  ) {
    addressString =
      '0x' +
      typedAddress.address.map((b) => b.toString(16).padStart(2, '0')).join('');
  }

  const hasValidAddress =
    addressString &&
    addressString.length === 42 &&
    addressString.startsWith('0x');

  return (
    <div
      className="modify-icons-container"
      style={{ display: 'flex', alignItems: 'center' }}
    >
      {/* Edit icon - available for non-deleted entries */}
      {!isDeleted && (
        <ModifyIcon
          type="edit"
          onClick={() => onAction('edit', row, tabName)}
          title="Edit this name entry"
        />
      )}

      {/* Delete icon - available for non-deleted entries */}
      {!isDeleted && (
        <ModifyIcon
          type="delete"
          onClick={() => onAction('delete', row, tabName)}
          title="Mark this name entry as deleted"
        />
      )}

      {/* Undelete icon - available only for deleted entries */}
      {isDeleted && (
        <ModifyIcon
          type="undelete"
          onClick={() => onAction('undelete', row, tabName)}
          title="Restore this deleted name entry"
        />
      )}

      {/* Remove icon - available only for deleted entries */}
      {isDeleted && (
        <ModifyIcon
          type="remove"
          onClick={() => onAction('remove', row, tabName)}
          title="Permanently remove this name entry"
        />
      )}

      {/* Autoname icon - available for valid Ethereum addresses */}
      {hasValidAddress && (
        <ModifyIcon
          type="autoname"
          onClick={() => onAction('autoname', row, tabName)}
          title="Auto-generate name from chain data"
        />
      )}
    </div>
  );
};
