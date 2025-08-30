import { Box, SimpleGrid, Title } from '@mantine/core';
import { dalle } from '@models';

import { DalleDressCard } from './DalleDressCard';

export interface SeriesGalleryProps {
  series: string;
  items: dalle.DalleDress[];
  columns: number;
  onItemClick?: (item: dalle.DalleDress) => void;
  onItemDoubleClick?: (item: dalle.DalleDress) => void;
  selectedRelPath?: string | null;
}

export const SeriesGallery = ({
  series,
  items,
  columns,
  onItemClick,
  onItemDoubleClick,
  selectedRelPath,
}: SeriesGalleryProps) => (
  <Box mb="lg">
    <Title order={5} mb={6} style={{ fontFamily: 'monospace' }}>
      {series || 'unknown'} ({items.length})
    </Title>
    <SimpleGrid cols={columns} spacing={6} verticalSpacing={6}>
      {items.map((it) => (
        <DalleDressCard
          key={it.annotatedPath}
          item={it}
          onClick={onItemClick}
          onDoubleClick={onItemDoubleClick}
          selected={it.annotatedPath === selectedRelPath}
        />
      ))}
    </SimpleGrid>
  </Box>
);
