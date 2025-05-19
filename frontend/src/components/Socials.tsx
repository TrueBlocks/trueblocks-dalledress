import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { useIcons } from '@hooks';
import { ActionIcon, Flex } from '@mantine/core';
import { BrowserOpenURL } from '@runtime';

export const Socials = () => {
  const [appName, setAppName] = useState<string>('');
  const [twitter, setTwitter] = useState<string>('');
  const [github, setGithub] = useState<string>('');
  const { Email, Website, Twitter, Github } = useIcons();

  useEffect(() => {
    GetAppId().then((id) => {
      setAppName(id.appName);
      setTwitter(id.twitter);
      setGithub(id.github);
    });
  }, []);

  const handleClick = (url: string) => {
    BrowserOpenURL(url);
  };

  const handleEmailClick = () => {
    window.location.href = `mailto:info@${appName}`;
  };

  return (
    <Flex gap="sm" align="center">
      <ActionIcon
        variant="subtle"
        size="sm"
        onClick={() => handleClick(`https://${appName}`)}
      >
        <Website />
      </ActionIcon>
      <ActionIcon
        variant="subtle"
        size="sm"
        onClick={() => handleClick(`${github}`)}
      >
        <Github />
      </ActionIcon>
      <ActionIcon
        variant="subtle"
        size="sm"
        onClick={() => handleClick(`https://x.com/${twitter}`)}
      >
        <Twitter />
      </ActionIcon>
      <ActionIcon variant="subtle" size="sm" onClick={handleEmailClick}>
        <Email />
      </ActionIcon>
    </Flex>
  );
};
