import { msgs } from '@models';
import { EventsEmit } from '@runtime';

export const emitEvent = (eventType: string, data?: unknown) => {
  EventsEmit(eventType, data);
};

export const createEventEmitters = () => {
  return {
    emitStatus: (message: string) => emitEvent(msgs.EventType.STATUS, message),
    emitError: (message: string) => emitEvent(msgs.EventType.ERROR, message),
    emitManager: (reason: string) => emitEvent(msgs.EventType.MANAGER, reason),
    emitAppInit: () => emitEvent(msgs.EventType.APP_INIT),
    emitAppReady: () => emitEvent(msgs.EventType.APP_READY),
    emitViewChange: (view: string) =>
      emitEvent(msgs.EventType.VIEW_CHANGE, view),
  };
};

export type EventEmitters = ReturnType<typeof createEventEmitters>;

export const useEventEmitters = (): EventEmitters => {
  return createEventEmitters();
};
