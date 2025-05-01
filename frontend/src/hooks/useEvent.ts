import { useEffect } from 'react';

import { EventsOff, EventsOn } from '@runtime';

export const useEvent = function <T = string>(
  eventType: string,
  callback: (data: T) => void,
) {
  useEffect(() => {
    EventsOn(eventType, callback);
    return () => {
      EventsOff(eventType);
    };
  }, [eventType, callback]);
};
