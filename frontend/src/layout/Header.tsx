import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { WalletConnectButton } from '@components';
import { useActiveProject } from '@hooks';
import { useIconSets } from '@hooks';
import {
  ActionIcon,
  AppShell,
  Group,
  Text,
  useMantineColorScheme,
} from '@mantine/core';

export const Header = () => {
  const [appName, setAppName] = useState('AppName');
  const { colorScheme, setColorScheme } = useMantineColorScheme();
  const { toggleDarkMode, isDarkMode } = useActiveProject();
  const { Light, Dark } = useIconSets();

  useEffect(() => {
    GetAppId().then((id) => {
      setAppName(id.appName);
    });
  }, []);

  useEffect(() => {
    setColorScheme(isDarkMode ? 'dark' : 'light');
  }, [isDarkMode, setColorScheme]);

  const handleToggleTheme = async () => {
    await toggleDarkMode();
  };

  return (
    <AppShell.Header>
      <Group justify="space-between" p="md" h="100%">
        <Text size="xl" fw={700}>
          {appName}
        </Text>
        <Group justify="flex-end" align="center" gap="xs">
          <ActionIcon
            variant="default"
            color={colorScheme === 'dark' ? 'yellow' : 'blue'}
            onClick={handleToggleTheme}
            title="Toggle color scheme"
          >
            {colorScheme === 'dark' ? <Light size={18} /> : <Dark size={18} />}
          </ActionIcon>
          <WalletConnectButton />
        </Group>
      </Group>
    </AppShell.Header>
  );
};
