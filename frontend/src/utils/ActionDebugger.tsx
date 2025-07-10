import React from 'react';

import { ActionDefinition } from '@hooks';
import { isDebugMode } from '@utils';

interface ActionDebuggerProps {
  rowActions: ActionDefinition[];
  headerActions: ActionDefinition[];
}

export const ActionDebugger: React.FC<ActionDebuggerProps> = ({
  rowActions,
  headerActions,
}) => {
  if (!isDebugMode()) {
    return <></>;
  }

  return (
    <div
      style={{
        backgroundColor: '#2c3e50',
        color: 'white',
        padding: '10px',
        marginBottom: '10px',
        fontSize: '13px',
        fontFamily: 'monospace',
        border: '1px solid #34495e',
        borderRadius: '4px',
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
      }}
    >
      <div>
        <ActionsList prompt="Row Actions:" actions={rowActions} />
        <ActionsList prompt="Header Actions:" actions={headerActions} />
      </div>
    </div>
  );
};

const ActionsList: React.FC<{
  prompt: string;
  actions: ActionDefinition[];
}> = ({ prompt, actions }) => {
  return (
    <>
      <strong style={{ color: '#ecf0f1' }}>{prompt}:</strong>{' '}
      {actions.length === 0 ? (
        <span style={{ color: '#e74c3c', fontStyle: 'italic' }}>None</span>
      ) : (
        actions.map((action) => {
          const actionTypeStyles: Record<
            string,
            { bgColor: string; textColor: string }
          > = {
            delete: { bgColor: '#e74c3c', textColor: 'white' },
            remove: { bgColor: '#e74c3c', textColor: 'white' },
            add: { bgColor: '#2ecc71', textColor: 'white' },
            update: { bgColor: '#2ecc71', textColor: 'white' },
            autoname: { bgColor: '#f39c12', textColor: 'white' },
            publish: { bgColor: '#9b59b6', textColor: 'white' },
            pin: { bgColor: '#9b59b6', textColor: 'white' },
          };

          const { bgColor, textColor } = actionTypeStyles[action.type] || {
            bgColor: '#3498db',
            textColor: 'white',
          };

          return (
            <span
              key={action.type}
              style={{
                display: 'inline-block',
                backgroundColor: bgColor,
                color: textColor,
                padding: '3px 8px',
                margin: '0 5px',
                borderRadius: '3px',
                fontWeight: 'bold',
                boxShadow: '0 1px 3px rgba(0,0,0,0.2)',
              }}
            >
              {action.type}
            </span>
          );
        })
      )}
    </>
  );
};
