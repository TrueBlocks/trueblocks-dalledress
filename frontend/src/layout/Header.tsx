import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import {
  ActionIcon,
  AppShell,
  Group,
  Text,
  useMantineColorScheme,
} from '@mantine/core';
import { FaMoon, FaSun } from 'react-icons/fa';

export const Header = () => {
  const [appName, setAppName] = useState('AppName');
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();

  useEffect(() => {
    GetAppId().then((id) => {
      setAppName(id.appName);
    });
  }, []);

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
            onClick={() => toggleColorScheme()}
            title="Toggle color scheme"
          >
            {colorScheme === 'dark' ? (
              <FaSun size={18} />
            ) : (
              <FaMoon size={18} />
            )}
          </ActionIcon>
          <Text>Header Content</Text>
        </Group>
      </Group>
    </AppShell.Header>
  );
};
