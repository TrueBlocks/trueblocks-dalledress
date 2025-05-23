import { useFormHotkeys } from '@components';
import { renderHook } from '@testing-library/react';
import { vi } from 'vitest';

import { resetAllCentralMocks, triggerHotkey } from '../../../__tests__/mocks';

vi.mock('react-hotkeys-hook', async () => {
  const mocks = await import('../../../__tests__/mocks');
  return { useHotkeys: mocks.mockUseHotkeys };
});

describe('useFormHotkeys', () => {
  beforeEach(() => {
    resetAllCentralMocks();
  });

  test('Enter key changes to edit mode when in display mode', () => {
    const setMode = vi.fn();

    renderHook(() =>
      useFormHotkeys({
        mode: 'display',
        setMode,
        onCancel: vi.fn(),
      }),
    );

    triggerHotkey('enter');

    expect(setMode).toHaveBeenCalledWith('edit');
  });

  test('should select all text in input when mod+a is pressed', () => {
    const mockSetSelectionRange = vi.fn();

    const input = document.createElement('input');
    input.value = 'test value';
    input.setSelectionRange = mockSetSelectionRange;
    document.body.appendChild(input);
    input.focus();

    renderHook(() =>
      useFormHotkeys({
        mode: 'edit',
        setMode: vi.fn(),
        onCancel: vi.fn(),
      }),
    );

    triggerHotkey('mod+a');

    expect(mockSetSelectionRange).toHaveBeenCalledWith(0, input.value.length);

    document.body.removeChild(input);
  });

  test('should switch from display to edit mode when enter is pressed', () => {
    const setMode = vi.fn();

    renderHook(() =>
      useFormHotkeys({
        mode: 'display',
        setMode,
        onCancel: vi.fn(),
      }),
    );

    triggerHotkey('enter');

    expect(setMode).toHaveBeenCalledWith('edit');
  });

  test('should call onCancel and switch to display mode when escape is pressed in edit mode', () => {
    const setMode = vi.fn();
    const onCancel = vi.fn();

    renderHook(() =>
      useFormHotkeys({
        mode: 'edit',
        setMode,
        onCancel,
      }),
    );

    triggerHotkey('esc');

    expect(setMode).toHaveBeenCalledWith('display');
    expect(onCancel).toHaveBeenCalled();
  });

  test('should find and click submit button in the form when no submitButtonRef is provided', () => {
    const form = document.createElement('form');
    const submitButton = document.createElement('button');
    submitButton.type = 'submit';
    submitButton.click = vi.fn();

    form.appendChild(submitButton);
    document.body.appendChild(form);

    const input = document.createElement('input');
    form.appendChild(input);
    input.focus();

    renderHook(() =>
      useFormHotkeys({
        mode: 'edit',
        setMode: vi.fn(),
        onCancel: vi.fn(),
      }),
    );

    triggerHotkey('enter');

    expect(submitButton.click).toHaveBeenCalled();

    document.body.removeChild(form);
  });
});
