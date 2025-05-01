import { useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { ProjectsList } from '@components';
import { Button, Card, Container, Grid, Text, Title } from '@mantine/core';

export function Home() {
  const [_, setAppName] = useState('Your App');
  useEffect(() => {
    GetAppId().then((id) => {
      setAppName(id.appName);
    });
  }, []);

  return (
    <Container size="lg" py="xl">
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
