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

// Mock useIcons hook
vi.mock('@hooks', async (importOriginal) => {
  try {
    const original = await importOriginal();
    return {
      ...(original as object),
      useIcons: () => ({
        Home: (props: any) => (
          <div data-testid="home-icon" {...props}>
            Home Icon
          </div>
        ),
        Khedra: (props: any) => (
          <div data-testid="khedra-icon" {...props}>
            Khedra Icon
          </div>
        ),
        // ADD_ROUTE
        Abis: (props: any) => (
          <div data-testid="abis-icon" {...props}>
            Abis Icon
          </div>
        ),
        Monitors: (props: any) => (
          <div data-testid="monitors-icon" {...props}>
            Monitors Icon
          </div>
        ),
        Chunks: (props: any) => (
          <div data-testid="chunks-icon" {...props}>
            Chunks Icon
          </div>
        ),
        Export: (props: any) => (
          <div data-testid="export-icon" {...props}>
            Export Icon
          </div>
        ),
        Names: (props: any) => (
          <div data-testid="names-icon" {...props}>
            Names Icon
          </div>
        ),
        // ADD_ROUTE
        DalleDress: (props: any) => (
          <div data-testid="dalledress-icon" {...props}>
            DalleDress Icon
          </div>
        ),
        Settings: (props: any) => (
          <div data-testid="settings-icon" {...props}>
            Settings Icon
          </div>
        ),
        Wizard: (props: any) => (
          <div data-testid="wizard-icon" {...props}>
            Wizard Icon
          </div>
        ),
        Switch: (props: any) => (
          <div data-testid="switch-icon" {...props}>
            Switch Icon
          </div>
        ),
        File: (props: any) => (
          <div data-testid="file-icon" {...props}>
            File Icon
          </div>
        ),
        Twitter: (props: any) => (
          <div data-testid="twitter-icon" {...props}>
            Twitter Icon
          </div>
        ),
        Github: (props: any) => (
          <div data-testid="github-icon" {...props}>
            Github Icon
          </div>
        ),
        Website: (props: any) => (
          <div data-testid="website-icon" {...props}>
            Website Icon
          </div>
        ),
        Email: (props: any) => (
          <div data-testid="email-icon" {...props}>
            Email Icon
          </div>
        ),
        Add: (props: any) => (
          <div data-testid="add-icon" {...props}>
            Add Icon
          </div>
        ),
        Edit: (props: any) => (
          <div data-testid="edit-icon" {...props}>
            Edit Icon
          </div>
        ),
        Delete: (props: any) => (
          <div data-testid="delete-icon" {...props}>
            Delete Icon
          </div>
        ),
        Undelete: (props: any) => (
          <div data-testid="undelete-icon" {...props}>
            Undelete Icon
          </div>
        ),
        Remove: (props: any) => (
          <div data-testid="remove-icon" {...props}>
            Remove Icon
          </div>
        ),
        ChevronLeft: (props: any) => (
          <div data-testid="chevronleft-icon" {...props}>
            ChevronLeft Icon
          </div>
        ),
        ChevronRight: (props: any) => (
          <div data-testid="chevronright-icon" {...props}>
            ChevronRight Icon
          </div>
        ),
        ChevronUp: (props: any) => (
          <div data-testid="chevronup-icon" {...props}>
            ChevronUp Icon
          </div>
        ),
        ChevronDown: (props: any) => (
          <div data-testid="chevrondown-icon" {...props}>
            ChevronDown Icon
          </div>
        ),
        Light: (props: any) => (
          <div data-testid="light-icon" {...props}>
            Light Icon
          </div>
        ),
        Dark: (props: any) => (
          <div data-testid="dark-icon" {...props}>
            Dark Icon
          </div>
        ),
        // Adding the new icons that will be needed for task #6
        Autoname: (props: any) => (
          <div data-testid="autoname-icon" {...props}>
            Autoname Icon
          </div>
        ),
      }),
    };
  } catch {
    return {
      useIcons: () => ({
        Home: (props: any) => (
          <div data-testid="home-icon" {...props}>
            Home Icon
          </div>
        ),
        Khedra: (props: any) => (
          <div data-testid="khedra-icon" {...props}>
            Khedra Icon
          </div>
        ),
        Names: (props: any) => (
          <div data-testid="names-icon" {...props}>
            Names Icon
          </div>
        ),
        DalleDress: (props: any) => (
          <div data-testid="dalledress-icon" {...props}>
            DalleDress Icon
          </div>
        ),
        Settings: (props: any) => (
          <div data-testid="settings-icon" {...props}>
            Settings Icon
          </div>
        ),
        Wizard: (props: any) => (
          <div data-testid="wizard-icon" {...props}>
            Wizard Icon
          </div>
        ),
        Switch: (props: any) => (
          <div data-testid="switch-icon" {...props}>
            Switch Icon
          </div>
        ),
        File: (props: any) => (
          <div data-testid="file-icon" {...props}>
            File Icon
          </div>
        ),
        Twitter: (props: any) => (
          <div data-testid="twitter-icon" {...props}>
            Twitter Icon
          </div>
        ),
        Github: (props: any) => (
          <div data-testid="github-icon" {...props}>
            Github Icon
          </div>
        ),
        Website: (props: any) => (
          <div data-testid="website-icon" {...props}>
            Website Icon
          </div>
        ),
        Email: (props: any) => (
          <div data-testid="email-icon" {...props}>
            Email Icon
          </div>
        ),
        Add: (props: any) => (
          <div data-testid="add-icon" {...props}>
            Add Icon
          </div>
        ),
        Edit: (props: any) => (
          <div data-testid="edit-icon" {...props}>
            Edit Icon
          </div>
        ),
        Delete: (props: any) => (
          <div data-testid="delete-icon" {...props}>
            Delete Icon
          </div>
        ),
        Undelete: (props: any) => (
          <div data-testid="undelete-icon" {...props}>
            Undelete Icon
          </div>
        ),
        Remove: (props: any) => (
          <div data-testid="remove-icon" {...props}>
            Remove Icon
          </div>
        ),
        ChevronLeft: (props: any) => (
          <div data-testid="chevronleft-icon" {...props}>
            ChevronLeft Icon
          </div>
        ),
        ChevronRight: (props: any) => (
          <div data-testid="chevronright-icon" {...props}>
            ChevronRight Icon
          </div>
        ),
        ChevronUp: (props: any) => (
          <div data-testid="chevronup-icon" {...props}>
            ChevronUp Icon
          </div>
        ),
        ChevronDown: (props: any) => (
          <div data-testid="chevrondown-icon" {...props}>
            ChevronDown Icon
          </div>
        ),
        Light: (props: any) => (
          <div data-testid="light-icon" {...props}>
            Light Icon
          </div>
        ),
        Dark: (props: any) => (
          <div data-testid="dark-icon" {...props}>
            Dark Icon
          </div>
        ),
        // Adding the new icons that will be needed for task #6
        Autoname: (props: any) => (
          <div data-testid="autoname-icon" {...props}>
            Autoname Icon
          </div>
        ),
      }),
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
