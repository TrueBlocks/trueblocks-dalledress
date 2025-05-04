import { fireEvent } from '@testing-library/react';
import { renderHook } from '@testing-library/react';
import { registerHotkeys } from '@utils';
import { vi } from 'vitest';

const mockUseHotkeys = vi.fn();

vi.mock('react-hotkeys-hook', () => ({
  useHotkeys: (
    key: string,
    callback: (event: KeyboardEvent) => void,
    options?: object,
  ) => {
    mockUseHotkeys(key, callback, options);

    if (typeof key === 'string') {
      document.addEventListener('keydown', (e) => {
        // Mock the behavior of useHotkeys
        if (
          (key === 'mod+a' && e.key === 'a' && (e.metaKey || e.ctrlKey)) ||
          (key === 'enter' && e.key === 'Enter') ||
          (key === 'mod+1' && e.key === '1' && (e.metaKey || e.ctrlKey)) ||
          (key === 'alt+1' && e.key === '1' && e.altKey)
        ) {
          callback(e);
        }
      });
    }
  },
}));

describe('registerHotkeys utility', () => {
  beforeEach(() => {
    mockUseHotkeys.mockClear();
    vi.clearAllMocks();
  });

  test('registers a single hotkey correctly', () => {
    const mockHandler = vi.fn();

    renderHook(() => {
      registerHotkeys([
        {
          key: 'mod+a',
          handler: mockHandler,
        },
      ]);
    });

    expect(mockUseHotkeys).toHaveBeenCalledTimes(1);
    expect(mockUseHotkeys).toHaveBeenCalledWith(
      'mod+a',
      mockHandler,
      undefined,
    );

    fireEvent.keyDown(document, { key: 'a', metaKey: true });
    expect(mockHandler).toHaveBeenCalled();
  });

  test('registers multiple hotkeys correctly', () => {
    const mockHandler1 = vi.fn();
    const mockHandler2 = vi.fn();

    renderHook(() => {
      registerHotkeys([
        {
          key: 'mod+1',
          handler: mockHandler1,
        },
        {
          key: 'enter',
          handler: mockHandler2,
        },
      ]);
    });

    expect(mockUseHotkeys).toHaveBeenCalledTimes(2);
    expect(mockUseHotkeys).toHaveBeenCalledWith(
      'mod+1',
      mockHandler1,
      undefined,
    );
    expect(mockUseHotkeys).toHaveBeenCalledWith(
      'enter',
      mockHandler2,
      undefined,
    );

    fireEvent.keyDown(document, { key: '1', metaKey: true });
    expect(mockHandler1).toHaveBeenCalled();

    fireEvent.keyDown(document, { key: 'Enter' });
    expect(mockHandler2).toHaveBeenCalled();
  });

  test('passes options correctly to useHotkeys', () => {
    const mockHandler = vi.fn();
    const options = { enableOnFormTags: true };

    renderHook(() => {
      registerHotkeys([
        {
          key: 'mod+a',
          handler: mockHandler,
          options,
        },
      ]);
    });

    expect(mockUseHotkeys).toHaveBeenCalledWith('mod+a', mockHandler, options);
  });

  test('fires correct handler for each key', () => {
    const mockHandler1 = vi.fn();
    const mockHandler2 = vi.fn();

    renderHook(() => {
      registerHotkeys([
        {
          key: 'mod+1',
          handler: mockHandler1,
        },
        {
          key: 'alt+1',
          handler: mockHandler2,
        },
      ]);
    });

    fireEvent.keyDown(document, { key: '1', metaKey: true });
    expect(mockHandler1).toHaveBeenCalled();
    expect(mockHandler2).not.toHaveBeenCalled();

    mockHandler1.mockClear();
    mockHandler2.mockClear();

    fireEvent.keyDown(document, { key: '1', altKey: true });
    expect(mockHandler1).not.toHaveBeenCalled();
    expect(mockHandler2).toHaveBeenCalled();
  });
});
