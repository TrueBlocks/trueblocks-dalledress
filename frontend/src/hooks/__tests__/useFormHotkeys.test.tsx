import { fireEvent, renderHook } from '@testing-library/react';
import { vi } from 'vitest';

import { useFormHotkeys } from '../useFormHotkeys';

vi.mock('react-hotkeys-hook', () => ({
  useHotkeys: (key: string, handler: (e: KeyboardEvent) => void) => {
    document.addEventListener('keydown', (e) => {
      if (
        (key === 'mod+a' && (e.metaKey || e.ctrlKey) && e.key === 'a') ||
        (key === 'enter' && e.key === 'Enter') ||
        (key === 'esc' && e.key === 'Escape')
      ) {
        handler(e as KeyboardEvent);
      }
    });
  },
}));

describe('useFormHotkeys', () => {
  test('ESC key calls onCancel in edit mode', () => {
    const setMode = vi.fn();
    const onCancel = vi.fn();

    renderHook(() =>
      useFormHotkeys({
        mode: 'edit',
        setMode,
        onCancel,
      }),
    );

    // Simulate ESC key press
    fireEvent.keyDown(document, { key: 'Escape', code: 'Escape' });

    // Verify that it changes mode and calls onCancel
    expect(setMode).toHaveBeenCalledWith('display');
    expect(onCancel).toHaveBeenCalled();
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

    // Simulate Enter key press
    fireEvent.keyDown(document, { key: 'Enter', code: 'Enter' });

    // Verify mode change
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

    fireEvent.keyDown(input, { key: 'a', metaKey: true });

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

    fireEvent.keyDown(document, { key: 'Enter' });

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

    fireEvent.keyDown(document, { key: 'Escape' });

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

    fireEvent.keyDown(input, { key: 'Enter' });

    expect(submitButton.click).toHaveBeenCalled();

    document.body.removeChild(form);
  });
});
