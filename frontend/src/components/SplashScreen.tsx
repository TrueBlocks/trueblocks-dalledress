import { Center, Loader, Stack, Text } from '@mantine/core';

interface SplashScreenProps {
  message?: string;
  showSpinner?: boolean;
}

export const SplashScreen = ({
  message = 'Loading TrueBlocks Dalledress...',
  showSpinner = true,
}: SplashScreenProps) => {
  return (
    <Center
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        width: '100vw',
        height: '100vh',
        backgroundColor: 'var(--mantine-color-body)',
        zIndex: 9999,
      }}
    >
      <Stack align="center" gap="md">
        {showSpinner && <Loader size="lg" />}
        <Text size="lg" c="dimmed">
          {message}
        </Text>
      </Stack>
    </Center>
  );
};
