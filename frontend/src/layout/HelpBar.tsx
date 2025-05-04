import { useEffect, useState } from 'react';

import { GetMarkdown, SetHelpCollapsed } from '@app';
import { ToggleButton } from '@components';
import { useAppContext } from '@contexts';
import { AppShell, Stack, Text } from '@mantine/core';
import Markdown from 'markdown-to-jsx';

export const HelpBar = () => {
  const [markdown, setMarkdown] = useState<string>('Loading...');
  const { currentLocation, lastTab } = useAppContext();
  const { helpCollapsed, setHelpCollapsed } = useAppContext();

  const toggleHelp = (open: boolean) => {
    setHelpCollapsed(open);
    SetHelpCollapsed(open);
  };

  useEffect(() => {
    var headerText = currentLocation.startsWith('/')
      ? currentLocation.slice(1)
      : currentLocation;
    headerText = `${headerText.charAt(0).toUpperCase() + headerText.slice(1)} View`;
    const fetchMarkdown = async () => {
      try {
        const content = await GetMarkdown(
          'help',
          currentLocation,
          lastTab[currentLocation] as string,
        );
        setMarkdown(`# ${headerText}\n\n${content}`);
      } catch (rawErr) {
        const errMsg =
          rawErr instanceof Error ? rawErr.message : String(rawErr);
        setMarkdown(`# ${headerText}\n\nError loading help content: ${errMsg}`);
      }
    };

    if (helpCollapsed) {
      setMarkdown(`# ${headerText}\n\nLoading...`);
    } else {
      fetchMarkdown();
    }
  }, [lastTab, currentLocation, helpCollapsed]);

  return (
    <AppShell.Aside p="md" style={{ transition: 'width 0.2s ease' }}>
      <ToggleButton
        collapsed={helpCollapsed}
        onToggle={() => toggleHelp(!helpCollapsed)}
        direction="right"
      />
      {helpCollapsed ? (
        <Text
          style={{
            transform: 'rotate(90deg)',
            whiteSpace: 'nowrap',
            position: 'absolute',
            top: 'calc(20px + 36px)',
            left: '36px',
            transformOrigin: 'left top',
            size: 'xs',
          }}
        >
          {markdown.split('\n')[0]?.replace('# ', '')}
        </Text>
      ) : (
        <Stack gap="sm" style={{ overflowY: 'auto' }}>
          <Markdown>{markdown}</Markdown>
        </Stack>
      )}
    </AppShell.Aside>
  );
};
