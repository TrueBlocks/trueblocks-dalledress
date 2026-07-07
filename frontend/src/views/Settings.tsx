import { useState } from 'react';
import { Button, Group, Stack, Tabs, Text, Title } from '@mantine/core';
import { usePersistedTab } from '@trueblocks/ui';
import { GetTab, SetTab, ValidateDalle } from '../../wailsjs/go/app/App';

const TABS = ['general', 'advanced'];

export function Settings() {
  const [status, setStatus] = useState('');
  const { activeTab, setActiveTab } = usePersistedTab({
    key: 'settings',
    defaultTab: 'general',
    loadTab: GetTab,
    saveTab: SetTab,
    tabs: TABS,
    cycleViewId: 'settings',
  });

  const validate = () => {
    ValidateDalle()
      .then(() => setStatus('Dalle data is valid.'))
      .catch((error: unknown) => setStatus(error instanceof Error ? error.message : String(error)));
  };

  return (
    <Stack>
      <Title order={2}>Settings</Title>
      <Tabs value={activeTab} onChange={(value) => value && setActiveTab(value)}>
        <Tabs.List mb="md">
          <Tabs.Tab value="general">General</Tabs.Tab>
          <Tabs.Tab value="advanced">Advanced</Tabs.Tab>
        </Tabs.List>
        <Tabs.Panel value="general">
          <Group>
            <Button onClick={validate}>Validate Dalle Data</Button>
            {status && <Text>{status}</Text>}
          </Group>
        </Tabs.Panel>
        <Tabs.Panel value="advanced">
          <Text c="dimmed">Advanced settings.</Text>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
