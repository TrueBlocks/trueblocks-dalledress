import { Badge, Group, Tooltip } from '@mantine/core';
import { types } from '@models';

// Define a ChipItem type for clarity
export interface ChipItem {
  id: string; // Unique identifier
  label: string; // Text to display
  color?: string; // Badge color (defaults to 'blue' if not provided)
  tooltip?: string; // Optional explanation tooltip
  clickValue: string; // Value to use for filtering when clicked
}

interface ChipsProps {
  items: ChipItem[];
  onChipClick?: (value: string) => void; // Optional click handler
}

export const Chips: React.FC<ChipsProps> = ({ items, onChipClick }) => {
  if (items.length === 0) {
    return (
      <Badge size="xs" color="gray" variant="outline" opacity={0.6}>
        None
      </Badge>
    );
  }

  return (
    <Group gap={4} justify="flex-start" align="center">
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
                cursor:
                  onChipClick && chip.clickValue
                    ? 'pointer'
                    : chip.tooltip
                      ? 'help'
                      : 'default',
                '&:hover': {
                  transform: 'translateY(-1px)',
                  boxShadow: theme.shadows.xs,
                },
              },
            })}
            onClick={() => {
              if (onChipClick && chip.clickValue) {
                onChipClick(chip.clickValue);
              }
            }}
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
      label: 'CONTRACT',
      color: 'blue',
      tooltip: 'This address is a smart contract',
      clickValue: 'contract',
    });
  }

  if (name.isErc20) {
    chipItems.push({
      id: 'erc20',
      label: 'ERC20',
      color: 'green',
      tooltip: 'This contract implements the ERC-20 token standard',
      clickValue: 'erc20',
    });
  }

  if (name.isErc721) {
    chipItems.push({
      id: 'erc721',
      label: 'NFT',
      color: 'grape',
      tooltip: 'This contract implements the ERC-721 NFT standard',
      clickValue: 'erc721',
    });
  }

  if (name.isPrefund) {
    chipItems.push({
      id: 'prefund',
      label: 'PREFUND',
      color: 'orange',
      tooltip: 'This address was allocated ETH in the genesis block',
      clickValue: 'prefund',
    });
  }

  if (name.isCustom) {
    chipItems.push({
      id: 'custom',
      label: 'CUSTOM',
      color: 'pink',
      tooltip: 'This address was edited by the user',
      clickValue: 'custom',
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
      if (name.isErc20 && lowerTag.includes('erc20')) return;
      if (
        name.isErc721 &&
        (lowerTag.includes('erc721') || lowerTag.includes('nft'))
      )
        return;
      if (name.isPrefund && lowerTag.includes('prefund')) return;
      if (name.isContract && lowerTag.includes('contract')) return;

      chipItems.push({
        id: `tag-${tag}`,
        label: tag.substring(0, 8),
        color: 'gray',
        tooltip: `Tag: ${tag}`,
        clickValue: tag, // Use the full tag for clickValue
      });
    });
  }

  return chipItems;
}
