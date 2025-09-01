import React from 'react';

import { describe, expect, it, vi } from 'vitest';

import { useAppHotkeys } from '../hooks/useAppHotkeys';
import { render, setupFocusedHookMocks, triggerHotkey } from './mocks';

describe('chromeCollapsed hotkey behavior only', () => {
  it('Hotkeys mod+m and mod+h inert while minimal mode active', () => {
    const setMenuCollapsed = vi.fn();
    const setHelpCollapsed = vi.fn();
    setupFocusedHookMocks({
      customPreferences: {
        chromeCollapsed: true,
        setMenuCollapsed,
        setHelpCollapsed,
      },
    });
    const HotkeyHost = () => {
      useAppHotkeys();
      return <div data-testid="hotkey-host" />;
    };
    render(<HotkeyHost />);
    triggerHotkey('mod+m');
    triggerHotkey('mod+h');
    expect(setMenuCollapsed).not.toHaveBeenCalled();
    expect(setHelpCollapsed).not.toHaveBeenCalled();
  });
});
