import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { ActionIcon, Flex } from '@mantine/core';
import { BrowserOpenURL } from '@runtime';
import { FaEnvelope, FaGithub, FaGlobe, FaTwitter } from 'react-icons/fa';

export const Socials = () => {
  const [appName, setAppName] = useState<string>('');
  const [twitter, setTwitter] = useState<string>('');
  const [github, setGithub] = useState<string>('');

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
        <FaGlobe />
      </ActionIcon>
      <ActionIcon
        variant="subtle"
        size="sm"
        onClick={() => handleClick(`${github}`)}
      >
        <FaGithub />
      </ActionIcon>
      <ActionIcon
        variant="subtle"
        size="sm"
        onClick={() => handleClick(`https://x.com/${twitter}`)}
      >
        <FaTwitter />
      </ActionIcon>
      <ActionIcon variant="subtle" size="sm" onClick={handleEmailClick}>
        <FaEnvelope />
      </ActionIcon>
    </Flex>
  );
};
