import { Box, SimpleGrid, Title } from '@mantine/core';
import { dalledress } from '@models';

import { DalleDressCard } from './DalleDressCard';

export interface SeriesGalleryProps {
  series: string;
  items: dalledress.GalleryItem[];
  columns: number;
  onItemClick?: (item: dalledress.GalleryItem) => void;
  onItemDoubleClick?: (item: dalledress.GalleryItem) => void;
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
          key={it.relPath}
          item={it}
          onClick={onItemClick}
          onDoubleClick={onItemDoubleClick}
          selected={it.relPath === selectedRelPath}
        />
      ))}
    </SimpleGrid>
  </Box>
);
