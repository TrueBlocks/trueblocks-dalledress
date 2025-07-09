import React from 'react';

import { types } from '@models';
import { isDebugMode } from '@utils';

interface ActionDebuggerProps {
  enabledActions: string[];
  setActiveFacet: (facet: types.DataFacet) => void;
}

export const ActionDebugger: React.FC<ActionDebuggerProps> = ({
  enabledActions,
  setActiveFacet,
}) => {
  if (!isDebugMode()) {
    return <></>;
  }

  const handleGoToCustom = () => {
    setActiveFacet(types.DataFacet.CUSTOM);
  };

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
        <strong style={{ color: '#ecf0f1' }}>Enabled Actions:</strong>{' '}
        {enabledActions.length === 0 ? (
          <span style={{ color: '#e74c3c', fontStyle: 'italic' }}>None</span>
        ) : (
          enabledActions.map((action) => {
            // Different colors for different action types
            let bgColor = '#3498db'; // Default blue
            let textColor = 'white';

            if (action === 'delete' || action === 'remove') {
              bgColor = '#e74c3c'; // Red for destructive actions
            } else if (action === 'add' || action === 'update') {
              bgColor = '#2ecc71'; // Green for constructive actions
            } else if (action === 'autoname') {
              bgColor = '#f39c12'; // Orange for automation
            } else if (action === 'publish' || action === 'pin') {
              bgColor = '#9b59b6'; // Purple for publishing actions
            }

            return (
              <span
                key={action}
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
                {action}
              </span>
            );
          })
        )}
      </div>
      <button
        onClick={handleGoToCustom}
        style={{
          backgroundColor: '#16a085',
          color: 'white',
          border: 'none',
          borderRadius: '4px',
          padding: '5px 10px',
          cursor: 'pointer',
          fontWeight: 'bold',
          boxShadow: '0 1px 3px rgba(0,0,0,0.3)',
        }}
      >
        Go to Custom Facet
      </button>
    </div>
  );
};
