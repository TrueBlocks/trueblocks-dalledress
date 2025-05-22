import { ReactElement, ReactNode } from 'react';

import { MantineProvider } from '@mantine/core';
import {
  RenderOptions,
  fireEvent,
  render as rtlRender,
  screen,
} from '@testing-library/react';
import { vi } from 'vitest';

type AppBridgeFunctions = {
  GetAppId: ReturnType<typeof vi.fn>;
  GetWizardReturn: ReturnType<typeof vi.fn>;
  GetUserPreferences: ReturnType<typeof vi.fn>;
  SetUserPreferences: ReturnType<typeof vi.fn>;
  IsReady: ReturnType<typeof vi.fn>;
  GetAppPreferences: ReturnType<typeof vi.fn>;
  SetLastView: ReturnType<typeof vi.fn>;
};

type RuntimeBridgeFunctions = {
  EventsEmit: ReturnType<typeof vi.fn>;
};

const createInitialMockAppBridge = (): AppBridgeFunctions => ({
  GetAppId: vi.fn().mockResolvedValue('mockAppId'),
  GetWizardReturn: vi.fn().mockResolvedValue('/mock-wizard-return'),
  GetUserPreferences: vi
    .fn()
    .mockResolvedValue({ name: 'Mock User', email: 'mock@example.com' }),
  SetUserPreferences: vi.fn().mockResolvedValue({}),
  IsReady: vi.fn().mockResolvedValue(true),
  GetAppPreferences: vi.fn().mockResolvedValue({
    lastView: '/mock-default-view',
    menuCollapsed: false,
  }),
  SetLastView: vi.fn().mockResolvedValue(undefined),
});

const createInitialMockRuntimeBridge = (): RuntimeBridgeFunctions => ({
  EventsEmit: vi.fn(),
});

export let mockAppBridge = createInitialMockAppBridge();
export let mockRuntimeBridge = createInitialMockRuntimeBridge();

vi.mock('../../wailsjs/go/app/App', () => mockAppBridge);
vi.mock('../../wailsjs/runtime/runtime', () => mockRuntimeBridge);

export const setupWailsMocks = (
  appOverrides?: Partial<AppBridgeFunctions>,
  runtimeOverrides?: Partial<RuntimeBridgeFunctions>,
) => {
  if (appOverrides) {
    for (const key in appOverrides) {
      const typedKey = key as keyof AppBridgeFunctions;
      if (
        Object.prototype.hasOwnProperty.call(mockAppBridge, typedKey) &&
        Object.prototype.hasOwnProperty.call(appOverrides, typedKey)
      ) {
        const overrideFn = appOverrides[typedKey];
        if (
          typeof mockAppBridge[typedKey]?.mockImplementation === 'function' &&
          typeof overrideFn === 'function'
        ) {
          mockAppBridge[typedKey].mockImplementation(overrideFn as any);
        } else {
          (mockAppBridge as any)[typedKey] = overrideFn;
        }
      } else {
        console.warn(
          `[setupWailsMocks] Attempted to override App.${typedKey}, but it's not defined in the central mockAppBridge or not provided in overrides.`,
        );
      }
    }
  }
  if (runtimeOverrides) {
    for (const key in runtimeOverrides) {
      const typedKey = key as keyof RuntimeBridgeFunctions;
      if (
        Object.prototype.hasOwnProperty.call(mockRuntimeBridge, typedKey) &&
        Object.prototype.hasOwnProperty.call(runtimeOverrides, typedKey)
      ) {
        const overrideFn = runtimeOverrides[typedKey];
        if (
          typeof mockRuntimeBridge[typedKey]?.mockImplementation ===
            'function' &&
          typeof overrideFn === 'function'
        ) {
          mockRuntimeBridge[typedKey].mockImplementation(overrideFn as any);
        } else {
          (mockRuntimeBridge as any)[typedKey] = overrideFn;
        }
      } else {
        console.warn(
          `[setupWailsMocks] Attempted to override Runtime.${typedKey}, but it's not defined in the central mockRuntimeBridge or not provided in overrides.`,
        );
      }
    }
  }
};

export function AllTheProviders({
  children,
}: {
  children: ReactNode;
}): ReactElement {
  return <MantineProvider>{children}</MantineProvider>;
}

function customRender(
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) {
  return rtlRender(ui, { wrapper: AllTheProviders, ...options });
}

export * from '@testing-library/react';
export { customRender as render, screen, fireEvent };

const registeredHotkeys = new Map<string, (e: KeyboardEvent) => void>();

const initialMockUseHotkeysImplementation = (
  key: string | string[],
  callback: (e: KeyboardEvent) => void,
) => {
  const keys = Array.isArray(key) ? key : [key];
  keys.forEach((k) => registeredHotkeys.set(k, callback));
};

export let mockUseHotkeys = vi.fn(initialMockUseHotkeysImplementation);

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

const createInitialViewContextDefaultValue = () => ({
  currentView: 'mockView',
  setCurrentView: vi.fn(),
  viewPagination: {},
  getPagination: vi.fn(() => ({ currentPage: 0, pageSize: 10, totalItems: 0 })),
  updatePagination: vi.fn(),
});
export let mockViewContextDefaultValue = createInitialViewContextDefaultValue();

export function setupContextMocks({
  customViewContext,
}: {
  customViewContext?: Partial<
    ReturnType<typeof createInitialViewContextDefaultValue>
  >;
} = {}) {
  const newDefaults = createInitialViewContextDefaultValue();
  if (customViewContext) {
    mockViewContextDefaultValue = {
      ...newDefaults,
      ...customViewContext,
      setCurrentView:
        customViewContext.setCurrentView || newDefaults.setCurrentView,
      getPagination:
        customViewContext.getPagination || newDefaults.getPagination,
      updatePagination:
        customViewContext.updatePagination || newDefaults.updatePagination,
    };
  } else {
    mockViewContextDefaultValue = newDefaults;
  }

  vi.mock('@contexts', async (importOriginal) => {
    const original = (await importOriginal()) as any;
    return {
      ...original,
      useViewContext: () => mockViewContextDefaultValue,
    };
  });
}

const createInitialTableContextDefaultValue = () => ({
  focusState: 'table' as 'table' | 'controls',
  selectedRowIndex: -1,
  setSelectedRowIndex: vi.fn(),
  focusTable: vi.fn(),
  focusControls: vi.fn(),
  tableRef: { current: null as HTMLTableElement | null },
});
export let mockTableContextDefaultValue =
  createInitialTableContextDefaultValue();

export function setupComponentHookMocks({
  customTableContext,
}: {
  customTableContext?: Partial<
    ReturnType<typeof createInitialTableContextDefaultValue>
  >;
} = {}) {
  const newDefaults = createInitialTableContextDefaultValue();
  if (customTableContext) {
    mockTableContextDefaultValue = {
      ...newDefaults,
      ...customTableContext,
      setSelectedRowIndex:
        customTableContext.setSelectedRowIndex ||
        newDefaults.setSelectedRowIndex,
      focusTable: customTableContext.focusTable || newDefaults.focusTable,
      focusControls:
        customTableContext.focusControls || newDefaults.focusControls,
    };
  } else {
    mockTableContextDefaultValue = newDefaults;
  }

  vi.mock('@components', async (importOriginal) => {
    const original = (await importOriginal()) as any;
    return {
      ...original,
      useTableContext: vi.fn(() => mockTableContextDefaultValue),
      useFormHotkeys: vi.fn(),
    };
  });
}

export function resetAllCentralMocks() {
  // Reset Wails Bridge mocks by re-initializing them to their default state
  mockAppBridge = createInitialMockAppBridge();
  mockRuntimeBridge = createInitialMockRuntimeBridge();

  // Reset Hotkeys mock
  mockUseHotkeys.mockReset(); // Clears history, calls, and implementation
  mockUseHotkeys.mockImplementation(initialMockUseHotkeysImplementation); // Restore initial impl.
  clearRegisteredHotkeys();

  // Reset Context mock values to their initial state
  mockViewContextDefaultValue = createInitialViewContextDefaultValue();
  mockTableContextDefaultValue = createInitialTableContextDefaultValue();
}
