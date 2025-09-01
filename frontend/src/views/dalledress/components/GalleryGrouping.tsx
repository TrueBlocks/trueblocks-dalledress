import { Box, SimpleGrid, Title } from '@mantine/core';
import { dalle } from '@models';
import { Log } from '@utils';

import { getItemKey } from '../store';
import { DalleDressCard } from './DalleDressCard';

export interface GalleryGroupingProps {
  items: dalle.DalleDress[];
  series: string;
  columns: number;
  selected?: string | null;
  onItemClick?: (item: dalle.DalleDress) => void;
  onItemDoubleClick?: (item: dalle.DalleDress) => void;
}

export const GalleryGrouping = ({
  series,
  items,
  columns,
  onItemClick,
  onItemDoubleClick,
  selected,
}: GalleryGroupingProps) => {
  Log(
    'GalleryGrouping:selected=' +
      String(selected) +
      ' items=' +
      items.map((i) => getItemKey(i)).join(','),
  );
  return (
    <Box mb="lg">
      <Title order={5} mb={6} style={{ fontFamily: 'monospace' }}>
        {series || 'unknown'} ({items.length})
      </Title>
      <SimpleGrid cols={columns} spacing={6} verticalSpacing={6}>
        {items.map((it) => {
          const itemKey = getItemKey(it);
          return (
            <DalleDressCard
              key={itemKey}
              item={it}
              onClick={onItemClick}
              onDoubleClick={onItemDoubleClick}
              selected={itemKey === selected}
            />
          );
        })}
      </SimpleGrid>
    </Box>
  );
};
