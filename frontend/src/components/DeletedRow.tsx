// DeletedIndicator.tsx - A component to display visual indicators for deleted rows
import { FC } from 'react';

import { Logger } from '@app';
import { ModifyIcons } from '@components';
import { types } from '@models';

interface DeletedRowProps {
  row: types.Name;
  onAction: (type: string, row: types.Name, tabName?: string) => void;
  tabName?: string;
}

export const DeletedRow: FC<DeletedRowProps> = ({ row, onAction, tabName }) => {
  const address = String(row.address || '');
  const isDeleted = row.deleted === true;

  // Debug the deletion status
  Logger(
    `Rendering row ${address}: deleted=${row.deleted}, isDeleted=${isDeleted}`,
  );

  return (
    <div style={{ display: 'flex', alignItems: 'center' }}>
      {isDeleted && (
        <span
          style={{
            backgroundColor: '#f44336',
            color: 'white',
            padding: '2px 6px',
            borderRadius: '4px',
            fontSize: '10px',
            marginRight: '8px',
          }}
        >
          DELETED
        </span>
      )}
      <ModifyIcons
        key={`${address}-${isDeleted ? 'deleted' : 'active'}-${Date.now()}`}
        row={row}
        onAction={onAction}
        tabName={tabName}
      />
    </div>
  );
};
