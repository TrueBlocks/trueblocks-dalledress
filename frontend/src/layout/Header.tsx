import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { AppShell, Group, Text } from '@mantine/core';

export const Header = () => {
  const [appName, setAppName] = useState('AppName');

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
        <Text>Header Content</Text>
      </Group>
    </AppShell.Header>
  );
};
