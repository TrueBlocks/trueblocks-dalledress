import { Box, Card, Image, Stack, Text } from '@mantine/core';
import { dalledress } from '@models';
import { Log, getDisplayAddress } from '@utils';

export interface DalleDressCardProps {
  item: dalledress.GalleryItem;
  onClick?: (item: dalledress.GalleryItem) => void;
  selected?: boolean;
}

export const DalleDressCard = ({
  item,
  onClick,
  selected,
}: DalleDressCardProps) => {
  const handleClick = () => {
    Log('gallery:click:' + item.relPath);
    onClick?.(item);
  };
  return (
    <Card
      p={4}
      radius="sm"
      withBorder
      style={{
        cursor: 'pointer',
        borderColor: selected ? 'var(--mantine-color-blue-5)' : undefined,
        boxShadow: selected
          ? '0 0 0 2px var(--mantine-color-blue-6), 0 0 6px 2px rgba(51,154,240,0.55)'
          : undefined,
        background: selected
          ? 'linear-gradient(135deg, rgba(51,154,240,0.18), rgba(51,154,240,0.05))'
          : undefined,
        transition: 'box-shadow 120ms, border-color 120ms, background 160ms',
      }}
      onClick={handleClick}
      data-relpath={item.relPath}
    >
      <Stack gap={4} align="stretch">
        <Box
          style={{
            position: 'relative',
            width: '100%',
            aspectRatio: '1 / 1',
            background: selected ? 'rgba(51,154,240,0.12)' : undefined,
            borderRadius: 4,
            overflow: 'hidden',
          }}
        >
          <Image
            src={item.url}
            alt={item.fileName}
            radius="xs"
            fit="cover"
            loading="lazy"
            style={{ width: '100%', height: '100%' }}
          />
        </Box>
        <Text size="xs" fw={500} truncate>
          {item.index >= 0 ? `#${item.index} ` : ''}
          {getDisplayAddress(item.address || '')}
        </Text>
      </Stack>
    </Card>
  );
};
