import { useEffect, useState } from 'react';
import { Button, Checkbox, Group, Paper, Stack, Tabs, Text, Textarea, Title } from '@mantine/core';
import { usePersistedTab } from '@trueblocks/ui';
import { DASHBOARD_PREFS, booleanPref, serializeBooleanPref } from '../dallePrefs';
import {
  GetPref,
  GetRuntimeInfo,
  GetTab,
  SetPref,
  SetTab,
  ValidateDalle,
} from '../../wailsjs/go/app/App';
import { app } from '../../wailsjs/go/models';

const TABS = ['general', 'defaults'];

function messageFromError(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

export function Settings() {
  const [status, setStatus] = useState('');
  const [runtime, setRuntime] = useState<app.RuntimeInfo | null>(null);
  const [defaultInput, setDefaultInput] = useState('Person Tour Coordinates');
  const [enhance, setEnhance] = useState(false);
  const [generateImage, setGenerateImage] = useState(false);
  const [annotate, setAnnotate] = useState(false);
  const [loaded, setLoaded] = useState(false);
  const { activeTab, setActiveTab } = usePersistedTab({
    key: 'settings',
    defaultTab: 'general',
    loadTab: GetTab,
    saveTab: SetTab,
    tabs: TABS,
    cycleViewId: 'settings',
  });

  useEffect(() => {
    Promise.all([
      GetRuntimeInfo(),
      GetPref(DASHBOARD_PREFS.input),
      GetPref(DASHBOARD_PREFS.enhance),
      GetPref(DASHBOARD_PREFS.generateImage),
      GetPref(DASHBOARD_PREFS.annotate),
    ])
      .then(([info, savedInput, savedEnhance, savedGenerateImage, savedAnnotate]) => {
        setRuntime(info);
        setDefaultInput(savedInput || 'Person Tour Coordinates');
        setEnhance(booleanPref(savedEnhance));
        setGenerateImage(booleanPref(savedGenerateImage));
        setAnnotate(booleanPref(savedAnnotate));
        setLoaded(true);
      })
      .catch((error: unknown) => setStatus(messageFromError(error)));
  }, []);

  useEffect(() => {
    if (!loaded) return;
    SetPref(DASHBOARD_PREFS.input, defaultInput);
  }, [defaultInput, loaded]);

  useEffect(() => {
    if (!loaded) return;
    SetPref(DASHBOARD_PREFS.enhance, serializeBooleanPref(enhance));
  }, [enhance, loaded]);

  useEffect(() => {
    if (!loaded) return;
    SetPref(DASHBOARD_PREFS.generateImage, serializeBooleanPref(generateImage));
  }, [generateImage, loaded]);

  useEffect(() => {
    if (!loaded) return;
    SetPref(DASHBOARD_PREFS.annotate, serializeBooleanPref(annotate));
  }, [annotate, loaded]);

  const validate = () => {
    ValidateDalle()
      .then(() => setStatus('Dalle data is valid.'))
      .catch((error: unknown) => setStatus(messageFromError(error)));
  };

  return (
    <Stack>
      <Title order={2}>Settings</Title>
      <Tabs value={activeTab} onChange={(value) => value && setActiveTab(value)}>
        <Tabs.List mb="md">
          <Tabs.Tab value="general">General</Tabs.Tab>
          <Tabs.Tab value="defaults">Generation Defaults</Tabs.Tab>
        </Tabs.List>
        <Tabs.Panel value="general">
          <Stack gap="md">
            <Group>
              <Button onClick={validate}>Validate Dalle Data</Button>
              {status && <Text>{status}</Text>}
            </Group>
            <Paper withBorder p="md">
              <Stack gap="xs">
                <Text fw={700}>Runtime</Text>
                <Text size="sm">Data directory: {runtime?.dataDir || 'unknown'}</Text>
                <Text size="sm">Database version: {runtime?.databaseVersion || 'unknown'}</Text>
                <Text size="xs" c="dimmed">
                  Archive: {runtime?.archiveHash || 'unknown'}
                </Text>
              </Stack>
            </Paper>
          </Stack>
        </Tabs.Panel>
        <Tabs.Panel value="defaults">
          <Paper withBorder p="md">
            <Stack gap="md">
              <Textarea
                label="Default seed text"
                value={defaultInput}
                minRows={3}
                onChange={(event) => setDefaultInput(event.currentTarget.value)}
              />
              <Group>
                <Checkbox
                  label="Enhance prompt"
                  checked={enhance}
                  onChange={(event) => setEnhance(event.currentTarget.checked)}
                />
                <Checkbox
                  label="Generate image"
                  checked={generateImage}
                  onChange={(event) => setGenerateImage(event.currentTarget.checked)}
                />
                <Checkbox
                  label="Annotate"
                  checked={annotate}
                  disabled={!generateImage}
                  onChange={(event) => setAnnotate(event.currentTarget.checked)}
                />
              </Group>
            </Stack>
          </Paper>
        </Tabs.Panel>
      </Tabs>
    </Stack>
  );
}
