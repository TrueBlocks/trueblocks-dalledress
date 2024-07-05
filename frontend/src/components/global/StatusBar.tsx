import React from 'react';

function StatusBar() {
  return (
    <div
      style={{
        width: '100%',
        position: 'fixed',
        bottom: 0,
        color: '#fff',
        borderTop: '1px solid #fff',
        paddingLeft: '-10px',
      }}
    >
      Status bar
    </div>
  );
}

export default StatusBar;
