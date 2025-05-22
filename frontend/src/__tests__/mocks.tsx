import { ReactElement, ReactNode } from 'react';

import { MantineProvider } from '@mantine/core';
import {
  RenderOptions,
  fireEvent,
  render as rtlRender,
  screen,
} from '@testing-library/react';
import { vi } from 'vitest';

const mockAppStatic = {
  SetHelpCollapsed: vi.fn(),
  SetInitialized: vi.fn(),
  SetMenuCollapsed: vi.fn(),
};

const mockRuntimeStatic = {
  EventsEmit: vi.fn(),
  EventsOn: vi.fn(),
  EventsOff: vi.fn(),
};

export const mockAppBridge = {
  ...mockAppStatic,
  GetAppId: vi.fn(() => Promise.resolve({ appName: 'Mocked App' })),
  GetWizardReturn: vi.fn(() => Promise.resolve('/mocked-path')),
  GetUserPreferences: vi.fn(() =>
    Promise.resolve({ name: 'Mock User', email: 'mock@example.com' }),
  ),
  SetUserPreferences: vi.fn(() => Promise.resolve({})),
  IsReady: vi.fn(() => Promise.resolve(true)),
  GetAppPreferences: vi.fn(() =>
    Promise.resolve({ lastView: '/', menuCollapsed: false }),
  ),
  SetLastView: vi.fn(() => Promise.resolve()),
};

export const mockRuntimeBridge = {
  ...mockRuntimeStatic,
};

export function setupWailsMocks() {
  vi.mock('@app', () => mockAppBridge);
  vi.mock('@runtime', () => mockRuntimeBridge);
}

const AllTheProviders = ({ children }: { children: ReactNode }) => {
  return <MantineProvider>{children}</MantineProvider>;
};

const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) => rtlRender(ui, { wrapper: AllTheProviders, ...options });

export * from '@testing-library/react';
export { customRender as render, screen, fireEvent };

const registeredHotkeys = new Map<string, (e: KeyboardEvent) => void>();
export const mockUseHotkeys = vi.fn(
  (key: string | string[], callback: (e: KeyboardEvent) => void) => {
    const keys = Array.isArray(key) ? key : [key];
    keys.forEach((k) => registeredHotkeys.set(k, callback));
  },
);

export function setupHotkeysMock() {
  vi.mock('react-hotkeys-hook', () => ({
    useHotkeys: mockUseHotkeys,
  }));
}

export function triggerHotkey(key: string, eventArgs?: Partial<KeyboardEvent>) {
  const handler = registeredHotkeys.get(key);
  if (handler) {
    const mockEvent = {
      key: key.split('+').pop() || key,
      metaKey: key.includes('mod+') || key.includes('meta+'),
      ctrlKey: key.includes('ctrl+'),
      shiftKey: key.includes('shift+'),
      altKey: key.includes('alt+'),
      preventDefault: vi.fn(),
      stopPropagation: vi.fn(),
      ...eventArgs,
    } as KeyboardEvent;
    handler(mockEvent);
    return mockEvent;
  }
  return null;
}

export function clearRegisteredHotkeys() {
  registeredHotkeys.clear();
}

export const mockViewContextDefaultValue = {
  currentView: 'mockView',
  setCurrentView: vi.fn(),
  viewPagination: {},
  getPagination: vi.fn(() => ({ currentPage: 0, pageSize: 10, totalItems: 0 })),
  updatePagination: vi.fn(),
};

export function setupContextMocks({
  customViewContext,
}: { customViewContext?: Partial<typeof mockViewContextDefaultValue> } = {}) {
  vi.mock('@contexts', async (importOriginal) => {
    const original = (await importOriginal()) as any;
    return {
      ...original,
      useViewContext: () => ({
        ...mockViewContextDefaultValue,
        ...customViewContext,
      }),
    };
  });
}

export const mockTableContextDefaultValue = {
  focusState: 'table',
  selectedRowIndex: -1,
  setSelectedRowIndex: vi.fn(),
  focusTable: vi.fn(),
  focusControls: vi.fn(),
  tableRef: { current: null },
};

export function setupComponentHookMocks({
  customTableContext,
}: { customTableContext?: Partial<typeof mockTableContextDefaultValue> } = {}) {
  vi.mock('@components', async (importOriginal) => {
    const original = (await importOriginal()) as any;
    return {
      ...original,
      useTableContext: vi.fn(() => ({
        ...mockTableContextDefaultValue,
        ...customTableContext,
      })),
      useFormHotkeys: vi.fn(),
    };
  });
}

export function resetAllCentralMocks() {
  Object.values(mockAppBridge).forEach(
    (fn) => typeof fn === 'function' && fn.mockClear(),
  );
  Object.values(mockRuntimeBridge).forEach(
    (fn) => typeof fn === 'function' && fn.mockClear(),
  );
  mockUseHotkeys.mockClear();
  clearRegisteredHotkeys();
  mockViewContextDefaultValue.setCurrentView.mockClear();
  mockViewContextDefaultValue.getPagination.mockClear();
  mockViewContextDefaultValue.updatePagination.mockClear();
  mockTableContextDefaultValue.setSelectedRowIndex.mockClear();
  mockTableContextDefaultValue.focusTable.mockClear();
  mockTableContextDefaultValue.focusControls.mockClear();
}
