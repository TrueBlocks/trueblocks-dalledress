import { useAppContext } from '@contexts';
import { Box, Stack, Text } from '@mantine/core';
import { useRoute } from 'wouter';

export const History = () => {
  const { currentLocation, selectedAddress } = useAppContext();
  const color = 'purple';
  const bgColor = 'white';

  const [, params] = useRoute('/history/:address');
  const addressFromUrl = params?.address;

  const displayAddress = selectedAddress || addressFromUrl;

  return (
    <Box style={{ backgroundColor: bgColor, minHeight: '100%' }}>
      <Stack style={{ color: color }}>
        <Text px="xl">{`THIS IS THE ${currentLocation.toUpperCase()} SCREEN`}</Text>
        {displayAddress ? (
          <Text px="xl">{`Displaying history for address: ${displayAddress}`}</Text>
        ) : (
          <Text px="xl">No address selected or provided in URL.</Text>
        )}
      </Stack>
    </Box>
  );
};
