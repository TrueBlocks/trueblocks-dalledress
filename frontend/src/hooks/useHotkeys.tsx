import { Logger, SetHelpCollapsed, SetMenuCollapsed } from '@app';
import { useAppContext } from '@contexts';
import { msgs } from '@models';
import { emitEvent, registerHotkeys } from '@utils';
import { MenuItems } from 'src/Menu';
import { useLocation } from 'wouter';

interface BaseHotkey {
  type: 'navigation' | 'dev' | 'toggle';
  hotkey: string;
  label: string;
}

interface NavigationHotkey extends BaseHotkey {
  type: 'navigation';
  path: string;
}

interface DevHotkey extends BaseHotkey {
  type: 'dev';
  path: string;
  action?: () => Promise<void>;
}

interface ToggleHotkey extends BaseHotkey {
  type: 'toggle';
  action: () => void;
}

type Hotkey = NavigationHotkey | DevHotkey | ToggleHotkey;

export const useAppHotkeys = (): void => {
  const { currentLocation } = useAppContext();
  const { menuCollapsed, setMenuCollapsed } = useAppContext();
  const { helpCollapsed, setHelpCollapsed } = useAppContext();
  const [, navigate] = useLocation();

  const handleHotkey = async (
    hkType: Hotkey,
    e: KeyboardEvent,
  ): Promise<void> => {
    e.preventDefault();
    e.stopPropagation();

    try {
      switch (hkType.type) {
        case 'navigation':
          if (currentLocation === hkType.path) {
            emitEvent(msgs.EventType.TAB_CYCLE, {
              route: hkType.path,
              key: hkType.hotkey,
            });
          } else {
            navigate(hkType.path);
          }
          break;

        case 'dev':
          if (!import.meta.env.DEV) return;
          if (hkType.action) await hkType.action();
          navigate(hkType.path);
          break;

        case 'toggle':
          hkType.action();
          break;
      }
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : String(error);
      Logger(errorMessage);

      if (
        (hkType.type === 'navigation' || hkType.type === 'dev') &&
        hkType.path
      ) {
        window.location.href = hkType.path;
      }
    }
  };

  const toggleHotkeys = [
    {
      key: 'mod+h',
      handler: (e: KeyboardEvent) =>
        handleHotkey(
          {
            type: 'toggle',
            hotkey: 'mod+h',
            label: 'Toggle help panel',
            action: () => {
              const next = !helpCollapsed;
              setHelpCollapsed(next);
              SetHelpCollapsed(next);
            },
          },
          e,
        ),
      options: { enableOnFormTags: true },
    },
    {
      key: 'mod+m',
      handler: (e: KeyboardEvent) =>
        handleHotkey(
          {
            type: 'toggle',
            hotkey: 'mod+m',
            label: 'Toggle menu panel',
            action: () => {
              const next = !menuCollapsed;
              setMenuCollapsed(next);
              SetMenuCollapsed(next);
            },
          },
          e,
        ),
      options: { enableOnFormTags: true },
    },
  ];

  const menuItemHotkeys = MenuItems.flatMap((item) => {
    const hotkeyConfigs = [];

    if (item.hotkey) {
      hotkeyConfigs.push({
        key: item.hotkey,
        handler: (e: KeyboardEvent) => {
          let hotkeyObj: Hotkey;
          const hotkey = item.hotkey || '';

          switch (item.type) {
            case 'dev':
              hotkeyObj = {
                type: 'dev',
                hotkey,
                path: item.path,
                label: `Navigate to ${item.label}`,
                action: item.action as () => Promise<void>,
              };
              break;
            case 'toggle':
              hotkeyObj = {
                type: 'toggle',
                hotkey,
                label: `Toggle ${item.label}`,
                action: item.action as () => void,
              };
              break;
            case 'navigation':
            default:
              hotkeyObj = {
                type: 'navigation',
                hotkey,
                path: item.path,
                label: `Navigate to ${item.label}`,
              };
              break;
          }
          handleHotkey(hotkeyObj, e);
        },
      });
    }

    if (item.altHotkey) {
      hotkeyConfigs.push({
        key: item.altHotkey,
        handler: (e: KeyboardEvent) => {
          let hotkeyObj: Hotkey;
          const hotkey = item.altHotkey || '';

          switch (item.type) {
            case 'dev':
              hotkeyObj = {
                type: 'dev',
                hotkey,
                path: item.path,
                label: `Navigate to ${item.label} (reverse)`,
                action: item.action as () => Promise<void>,
              };
              break;
            case 'toggle':
              hotkeyObj = {
                type: 'toggle',
                hotkey,
                label: `Toggle ${item.label} (reverse)`,
                action: item.action as () => void,
              };
              break;
            case 'navigation':
            default:
              hotkeyObj = {
                type: 'navigation',
                hotkey,
                path: item.path,
                label: `Navigate to ${item.label} (reverse)`,
              };
              break;
          }
          handleHotkey(hotkeyObj, e);
        },
      });
    }

    return hotkeyConfigs;
  });

  const hotkeysConfig = [...menuItemHotkeys, ...toggleHotkeys];
  registerHotkeys(hotkeysConfig);
};
