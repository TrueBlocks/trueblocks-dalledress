import { useAppContext } from '@contexts';
import { Box, Stack, Text } from '@mantine/core';

export const Data = () => {
  const { currentLocation } = useAppContext();
  const color = 'purple';
  const bgColor = 'white';

  return (
    <Box style={{ backgroundColor: bgColor, minHeight: '100%' }}>
      <Stack style={{ color: color }}>
        <Text px="xl">{`THIS IS THE ${currentLocation.toUpperCase()} SCREEN`}</Text>
      </Stack>
    </Box>
  );
};
