import { ReactElement } from 'react';

import { MantineProvider } from '@mantine/core';
import { RenderOptions, render } from '@testing-library/react';
import { vi } from 'vitest';

vi.mock('@app', () => ({
  Logger: vi.fn(),
  SetHelpCollapsed: vi.fn(),
  SetInitialized: vi.fn(),
  SetMenuCollapsed: vi.fn(),
}));

vi.mock('@runtime', () => ({
  EventsEmit: vi.fn(),
  EventsOn: vi.fn(),
  EventsOff: vi.fn(),
}));

const AllTheProviders = ({ children }: { children: React.ReactNode }) => {
  return <MantineProvider>{children}</MantineProvider>;
};

const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) => render(ui, { wrapper: AllTheProviders, ...options });

// Re-export everything
export * from '@testing-library/react';
export { customRender as render };
