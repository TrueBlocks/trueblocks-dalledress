import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { useIcons } from '@hooks';
import {
  ActionIcon,
  AppShell,
  Group,
  Text,
  useMantineColorScheme,
} from '@mantine/core';

export const Header = () => {
  const [appName, setAppName] = useState('AppName');
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();
  const { Light, Dark } = useIcons();

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
            {colorScheme === 'dark' ? <Light size={18} /> : <Dark size={18} />}
          </ActionIcon>
          <Text>Header Content</Text>
        </Group>
      </Group>
    </AppShell.Header>
  );
};
