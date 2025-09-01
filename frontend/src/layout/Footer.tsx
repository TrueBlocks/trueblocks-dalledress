import { useEffect, useState } from 'react';

import { GetFilename, GetOrgPreferences } from '@app';
import { ChevronButton, Socials, getBarSize } from '@components';
import { useEvent, usePreferences } from '@hooks';
import { AppShell, Flex, Group, Text } from '@mantine/core';
import { msgs, preferences, project } from '@models';

export const Footer = () => {
  var [org, setOrg] = useState<preferences.OrgPreferences>({});
  const { menuCollapsed, chromeCollapsed, setChromeCollapsed } =
    usePreferences();

  useEffect(() => {
    const fetchOrgName = async () => {
      GetOrgPreferences().then((data) => {
        setOrg(data);
      });
    };
    fetchOrgName();
  }, []);

  return (
    <AppShell.Footer ml={getBarSize('menu', menuCollapsed) - 1} p={0}>
      <Flex
        h="100%"
        px={chromeCollapsed ? 4 : 'md'}
        align="center"
        justify="space-between"
        w="100%"
      >
        <Group gap={4} align="center">
          <ChevronButton
            collapsed={chromeCollapsed}
            onToggle={() => setChromeCollapsed(!chromeCollapsed)}
            direction={chromeCollapsed ? 'right' : 'left'}
          />
        </Group>
        {!chromeCollapsed && (
          <>
            <FilePanel />
            <Text size="sm">{org.developerName} © 2025</Text>
            <Socials />
          </>
        )}
      </Flex>
    </AppShell.Footer>
  );
};

export const FilePanel = () => {
  const [status, setStatus] = useState<project.Project | null>(null);

  useEffect(() => {
    const fetchFilename = async () => {
      setStatus(await GetFilename());
    };
    fetchFilename();
  }, []);

  useEvent(msgs.EventType.MANAGER, async (_message?: string) => {
    setStatus(await GetFilename());
  });

  return (
    <>{status ? <Text>{status.name}</Text> : <Text>No Open Project</Text>}</>
  );
};
