import { SetMenuCollapsed } from '@app';
import { ToggleButton, getBarWidth } from '@components';
import { useAppContext } from '@contexts';
import { AppShell, Button, Stack } from '@mantine/core';
import { MenuItem, MenuItems } from 'src/Menu';

interface MenuBarProps {
  disabled?: boolean;
}

export const MenuBar = ({ disabled = false }: MenuBarProps) => {
  const { currentLocation, navigate } = useAppContext();
  const { menuCollapsed, setMenuCollapsed } = useAppContext();

  const toggleMenu = (open: boolean) => {
    setMenuCollapsed(open);
    SetMenuCollapsed(open).then(() => {});
  };

  const topMenuItems = MenuItems.filter((item) => item.position === 'top');
  const botMenuItems = MenuItems.filter((item) => item.position === 'bottom');

  const renderMenuItem = ({ icon: Icon, label, path }: MenuItem) => (
    <Button
      key={path}
      variant={currentLocation === path ? 'filled' : 'subtle'}
      fullWidth
      h={36}
      w={menuCollapsed ? 36 : '100%'}
      leftSection={
        <Icon size={16} style={{ marginLeft: menuCollapsed ? 9 : 0 }} />
      }
      justify={menuCollapsed ? 'center' : 'flex-start'}
      px={menuCollapsed ? 0 : 'md'}
      style={{
        marginLeft: menuCollapsed ? -9 : 0,
      }}
      disabled={disabled}
      onClick={() => {
        if (!disabled) navigate(path);
      }}
    >
      {!menuCollapsed && label}
    </Button>
  );

  return (
    <AppShell.Navbar
      p="md"
      style={{
        paddingTop: 0,
        paddingBottom: 0,
        height: 'calc(100vh - 30px)',
        width: getBarWidth(menuCollapsed, 1),
        transition: 'width 0.2s ease',
      }}
    >
      <Stack h="100%" justify="space-between" gap="sm">
        <Stack gap="sm">
          <ToggleButton
            collapsed={menuCollapsed}
            onToggle={() => toggleMenu(!menuCollapsed)}
            direction="left"
          />
          {topMenuItems.map(renderMenuItem)}
        </Stack>
        <Stack gap="sm" pb="md">
          {botMenuItems.map(renderMenuItem)}
        </Stack>
      </Stack>
    </AppShell.Navbar>
  );
};
