import '@testing-library/jest-dom';
import { beforeEach, vi } from 'vitest';

beforeEach(() => {
  vi.clearAllMocks();
});

vi.mock('@utils', async (importOriginal) => {
  try {
    const original = await importOriginal();
    return {
      ...(original as object),
      Log: vi.fn(),
      // Wizard-specific utilities
      checkAndNavigateToWizard: () => Promise.resolve(null),
      useEmitters: () => ({ emitStatus: vi.fn(), emitError: vi.fn() }),
    };
  } catch {
    return {
      Log: vi.fn(),
      checkAndNavigateToWizard: () => Promise.resolve(null),
      useEmitters: () => ({ emitStatus: vi.fn(), emitError: vi.fn() }),
    };
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
