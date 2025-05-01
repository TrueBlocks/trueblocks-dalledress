import { useEffect } from 'react';

import { MantineProvider } from '@mantine/core';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';

import { Wizard } from '../Wizard';

// Minimal mock for external APIs only
vi.mock('@app', () => ({
  GetAppId: () => Promise.resolve({ appName: 'Test App' }),
  GetWizardReturn: () => Promise.resolve('/'),
  GetUserPreferences: () => Promise.resolve({ name: '', email: '' }),
  SetUserPreferences: () => Promise.resolve({}),
}));

vi.mock('@runtime', () => ({
  EventsEmit: vi.fn(),
}));

const mockNavigate = vi.fn();
const mockCompleteWizard = vi.fn().mockResolvedValue(undefined);

// Simple mock for AppContext
vi.mock('@contexts', () => ({
  useAppContext: () => ({
    navigate: mockNavigate,
    isWizard: true,
  }),
}));

// Minimal mock for wizardUtils
vi.mock('@utils', () => ({
  checkAndNavigateToWizard: () => Promise.resolve(null),
}));

// We need to make the step navigation actually work in our test
let activeStep = 0;
const mockUpdateUI = vi.fn().mockImplementation((data) => {
  if (data && typeof data.activeStep === 'number') {
    activeStep = data.activeStep;
  }
});

// Mock for hooks with better state handling - include ALL required exports
vi.mock('../hooks', () => ({
  useWizardState: () => ({
    state: {
      data: {
        name: '',
        email: '',
        isFirstTimeSetup: true,
      },
      ui: { activeStep },
      validation: {},
    },
    updateData: vi.fn(),
    updateValidation: vi.fn(),
    updateUI: mockUpdateUI,
    submitUserInfo: vi.fn().mockImplementation(() => {
      // This simulates moving to step 1 when user info is submitted
      mockUpdateUI({ activeStep: 1 });
      return Promise.resolve();
    }),
    submitChainInfo: vi.fn().mockImplementation(() => {
      // This simulates moving to step 2 when RPC is submitted
      mockUpdateUI({ activeStep: 2 });
      return Promise.resolve();
    }),
    completeWizard: mockCompleteWizard,
  }),
  useWizardNavigation: () => ({}),
  useWizardValidation: () => ({
    validateName: vi.fn().mockReturnValue(true),
    validateEmail: vi.fn().mockReturnValue(true),
    validateRpc: vi.fn().mockReturnValue(true),
    validateUserInfo: vi.fn().mockReturnValue(true),
  }),
}));

// Mock for useFormHotkeys to handle ESC key
vi.mock('@hooks', () => ({
  useFormHotkeys: ({ onCancel }: { onCancel?: () => void }) => {
    useEffect(() => {
      const handleKeyDown = (event: KeyboardEvent) => {
        if (event.key === 'Escape') {
          onCancel?.();
        }
      };

      document.addEventListener('keydown', handleKeyDown);
      return () => {
        document.removeEventListener('keydown', handleKeyDown);
      };
    }, [onCancel]);
  },
}));

describe('Wizard', () => {
  beforeEach(() => {
    // Reset active step before each test
    activeStep = 0;
    vi.clearAllMocks();
  });

  test('renders first step with user info form', async () => {
    render(
      <MantineProvider>
        <Wizard />
      </MantineProvider>,
    );

    // Wait for initial step to be visible and inputs to appear
    await waitFor(() => {
      expect(
        screen.getByText('User Information', {
          selector: 'span.mantine-Stepper-stepLabel',
        }),
      ).toBeInTheDocument();
    });

    expect(screen.getByPlaceholderText('Enter your name')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Enter your email')).toBeInTheDocument();
  });

  test('can navigate to next step when fields are filled', async () => {
    render(
      <MantineProvider>
        <Wizard />
      </MantineProvider>,
    );

    // Fill in required fields using more robust selectors
    const nameInput = screen.getByPlaceholderText('Enter your name');
    const emailInput = screen.getByPlaceholderText('Enter your email');

    fireEvent.change(nameInput, { target: { value: 'Test User' } });
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });

    // Click Next button using getByRole instead of direct node access
    const nextButton = screen.getByRole('button', { name: /Next/i });
    fireEvent.click(nextButton);

    // Verify navigation to Chain Configuration step
    await waitFor(() => {
      expect(
        screen.getByText('Chain Configuration', {
          selector: 'span.mantine-Stepper-stepLabel',
        }),
      ).toBeInTheDocument();
    });
  });

  // test('can navigate to next step and then back', async () => {
  //   render(
  //     <MantineProvider>
  //       <Wizard />
  //     </MantineProvider>,
  //   );

  //   // Fill in required fields
  //   const nameInput = screen.getByPlaceholderText('Enter your name');
  //   const emailInput = screen.getByPlaceholderText('Enter your email');

  //   fireEvent.change(nameInput, { target: { value: 'Test User' } });
  //   fireEvent.change(emailInput, { target: { value: 'test@example.com' } });

  //   // Click Next button
  //   const nextButton = screen.getByRole('button', { name: 'Next' });
  //   fireEvent.click(nextButton);

  //   // Wait for active step to update
  //   await waitFor(() => {
  //     expect(activeStep).toBe(1);
  //   });

  //   // Look for Chain Configuration in the Stepper label
  //   await waitFor(() => {
  //     const stepLabel = screen.getByText('Chain Configuration', {
  //       selector: 'span.mantine-Stepper-stepLabel',
  //     });
  //     expect(stepLabel).toBeInTheDocument();
  //   });

  //   // Check that Back button now exists and click it
  //   const backButton = screen.getByRole('button', { name: 'Back' });
  //   expect(backButton).toBeInTheDocument();

  //   fireEvent.click(backButton);

  //   // Verify we navigate back to User Information step
  //   await waitFor(() => {
  //     expect(activeStep).toBe(0);
  //   });

  //   await waitFor(() => {
  //     expect(
  //       screen.getByText('User Information', {
  //         selector: 'span.mantine-Stepper-stepLabel',
  //       }),
  //     ).toBeInTheDocument();
  //   });
  // });

  // test.skip('escape key closes the wizard', async () => {
  //   render(
  //     <MantineProvider>
  //       <Wizard />
  //     </MantineProvider>,
  //   );

  //   // First, ensure the component is fully mounted
  //   await waitFor(() => {
  //     expect(
  //       screen.getByText('User Information', {
  //         selector: 'span.mantine-Stepper-stepLabel',
  //       }),
  //     ).toBeInTheDocument();
  //   });

  //   // Now trigger the ESC key on the document - this should call our mocked event listener
  //   fireEvent.keyDown(document, { key: 'Escape', code: 'Escape' });

  //   // Verify navigation was called
  //   await waitFor(() => {
  //     expect(mockNavigate).toHaveBeenCalled();
  //   });
  // });
});
