import '@testing-library/jest-dom';
import { vi } from 'vitest';

// Mock @utils globally
vi.mock('@utils', async (importOriginal) => {
  try {
    const original = await importOriginal();
    return {
      ...(original as object),
      Log: vi.fn(),
    };
  } catch {
    return { Log: vi.fn() };
  }
});

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {}, // Deprecated
    removeListener: () => {}, // Deprecated
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  }),
});
