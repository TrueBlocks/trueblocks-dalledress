import { useEffect, useState } from 'react';

import { useEvent } from '@hooks';
import { msgs } from '@models';

export const StatusBar = () => {
  const [status, setStatus] = useState('');
  const [visible, setVisible] = useState(false);

  useEvent(msgs.EventType.STATUS, (message: string) => {
    setStatus(message);
    setVisible(true);
  });

  useEffect(() => {
    if (!visible) return;
    const timer = setTimeout(() => {
      setVisible(false);
    }, 1500);
    return () => clearTimeout(timer);
  }, [visible, status]);

  if (!visible) return null;

  return (
    <div style={{ backgroundColor: '$cffafe', color: 'black' }}>
      <span>{status}</span>
    </div>
  );
};
