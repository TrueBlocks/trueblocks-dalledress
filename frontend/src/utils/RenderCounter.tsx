import React from 'react';

import { isDebugMode } from '@utils';

interface RenderCounterProps {
  count: number;
}

export const RenderCounter: React.FC<RenderCounterProps> = ({ count }) => {
  if (!isDebugMode()) {
    return <></>;
  }

  return (
    <div
      style={{
        backgroundColor: '#8e44ad',
        color: 'white',
        padding: '8px 12px',
        margin: '5px 0',
        fontSize: '12px',
        fontFamily: 'monospace',
        border: '1px solid #9b59b6',
        borderRadius: '4px',
        boxShadow: '0 1px 3px rgba(0,0,0,0.2)',
        display: 'inline-block',
        fontWeight: 'bold',
      }}
    >
      Render Count: {count}
    </div>
  );
};
