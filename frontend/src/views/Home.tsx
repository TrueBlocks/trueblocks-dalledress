import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { ProjectsList } from '@components';
import {
  Button,
  Card,
  Container,
  Grid,
  Image,
  Paper,
  Space,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { msgs } from '@models';

import { GetImageURL } from '../../wailsjs/go/app/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export function Home() {
  const [_, setAppName] = useState('Your App');
  const [sampleImageUrl, setSampleImageUrl] = useState<string>('');
  const [refreshKey, setRefreshKey] = useState(0);

  useEffect(() => {
    GetAppId().then((id) => {
      setAppName(id.appName);
    });

    const cacheBuster = new Date().getTime();
    GetImageURL(`samples/sample1.png?v=${cacheBuster}`)
      .then((url) => {
        setSampleImageUrl(url || '');
      })
      .catch((err) => {
        setSampleImageUrl('');
      });

    const off = EventsOn(msgs.EventType.IMAGES_CHANGED, () => {
      setRefreshKey((prev) => prev + 1);
    });

    return () => {
      off();
    };
  }, [refreshKey]);

  return (
    <Container size="lg" py="xl">
      <Space h="md" />
      <Text>
        Generate stunning images from Ethereum addresses, block hashes, or
        transaction hashes.
      </Text>
      <Space h="xl" />

      <Paper p="md" withBorder>
        <Stack>
          <Title order={3}>Sample Image</Title>
          {sampleImageUrl ? (
            <>
              <Image
                src={sampleImageUrl}
                alt="Sample Image"
                maw={300}
                key={refreshKey}
                radius="md"
                fit="contain"
                fallbackSrc="https://placehold.co/300x200?text=No+Image"
                loading="lazy"
                style={{ objectFit: 'contain' }}
              />
              <Text size="sm" color="dimmed">
                Sample image from file server
              </Text>
            </>
          ) : (
            <Text>Loading sample image...</Text>
          )}
          <Text size="sm">
            Images are served from the internal file server and can be found in
            your documents folder.
          </Text>
        </Stack>
      </Paper>

      <Grid>
        <Grid.Col span={6}>
          <Card shadow="sm" p="lg" radius="md" withBorder>
            <Title order={3} mb="md">
              Your Projects
            </Title>
            <ProjectsList />
          </Card>
        </Grid.Col>

        <Grid.Col span={6}>
          <Card shadow="sm" p="lg" radius="md" withBorder>
            <Title order={3} mb="md">
              Transactions
            </Title>
            <Text mb="md">
              Search and view detailed transaction information.
            </Text>
            <Button fullWidth>View Transactions</Button>
          </Card>
        </Grid.Col>
      </Grid>
    </Container>
  );
}
