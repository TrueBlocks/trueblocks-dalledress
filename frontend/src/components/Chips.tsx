import { Badge, Group, Tooltip } from '@mantine/core';
import { types } from '@models';

// Define a ChipItem type for clarity
export interface ChipItem {
  id: string; // Unique identifier
  label: string; // Text to display
  color?: string; // Badge color (defaults to 'blue' if not provided)
  tooltip?: string; // Optional explanation tooltip
}

interface ChipsProps {
  items: ChipItem[];
}

export const Chips: React.FC<ChipsProps> = ({ items }) => {
  if (items.length === 0) {
    return (
      <Badge size="xs" color="gray" variant="outline" opacity={0.6}>
        None
      </Badge>
    );
  }

  return (
    <Group gap={4} justify="flex-start" align="center">
      {' '}
      {/* Replaced SimpleGrid with Group, using gap prop */}
      {items.map((chip) => (
        <Tooltip
          key={chip.id}
          label={chip.tooltip || ''}
          disabled={!chip.tooltip}
        >
          <Badge
            size="xs"
            color={chip.color || 'blue'}
            variant="light"
            styles={(theme) => ({
              root: {
                transition: 'transform 0.2s ease, box-shadow 0.2s ease',
                '&:hover': {
                  transform: 'translateY(-1px)',
                  boxShadow: theme.shadows.xs,
                  cursor: chip.tooltip ? 'help' : 'default',
                },
              },
            })}
          >
            {chip.label}
          </Badge>
        </Tooltip>
      ))}
    </Group>
  );
};

// Helper function to map a Name object to ChipItems
export function mapNameToChips(name: types.Name): ChipItem[] {
  const chipItems: ChipItem[] = [];

  if (name.isContract) {
    chipItems.push({
      id: 'contract',
      label: 'CNTRCT', // Shortened label to 6 chars
      color: 'blue',
      tooltip: 'This address is a smart contract',
    });
  }

  if (name.isErc20) {
    chipItems.push({
      id: 'erc20',
      label: 'ERC20', // Already <= 6 chars
      color: 'green',
      tooltip: 'This contract implements the ERC-20 token standard',
    });
  }

  if (name.isErc721) {
    chipItems.push({
      id: 'erc721',
      label: 'NFT', // Already <= 6 chars
      color: 'grape',
      tooltip: 'This contract implements the ERC-721 NFT standard',
    });
  }

  if (name.isPrefund) {
    chipItems.push({
      id: 'prefund',
      label: 'PRFUND', // Label is 6 chars
      color: 'orange',
      tooltip: 'This address was allocated ETH in the genesis block',
    });
  }

  if (name.isCustom) {
    chipItems.push({
      id: 'custom',
      label: 'CUSTOM', // Shortened label to 6 chars
      color: 'pink',
      tooltip: 'This address was edited by the user',
    });
  }

  // Add chips for tags, avoiding duplicates with boolean flags
  if (name.tags && typeof name.tags === 'string' && name.tags.trim() !== '') {
    const tagsArray = name.tags
      .split(',')
      .map((tag) => tag.trim())
      .filter((tag) => tag !== '');

    tagsArray.forEach((tag: string) => {
      const lowerTag = tag.toLowerCase();
      // Check if the tag corresponds to an existing boolean chip
      if (name.isErc20 && lowerTag.includes('erc20')) return;
      if (
        name.isErc721 &&
        (lowerTag.includes('erc721') || lowerTag.includes('nft'))
      )
        return;
      if (name.isPrefund && lowerTag.includes('prefund')) return;
      if (name.isContract && lowerTag.includes('contract')) return;
      // Consider adding deduplication for 'custom' if a 'custom' tag might appear
      // if (name.isCustom && lowerTag.includes('custom')) return;

      chipItems.push({
        id: `tag-${tag}`,
        label: tag.substring(0, 6), // Shortened label to 6 chars
        color: 'gray',
        tooltip: `Tag: ${tag}`,
      });
    });
  }

  return chipItems;
}
