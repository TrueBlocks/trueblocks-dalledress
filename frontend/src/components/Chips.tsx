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
    <Group gap="xs">
      {items.map((chip) => (
        <Tooltip
          key={chip.id}
          label={chip.tooltip || ''}
          disabled={!chip.tooltip}
        >
          <Badge
            size="xs"
            color={chip.color || 'blue'}
            variant="filled"
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
      label: 'Contract',
      color: 'blue',
      tooltip: 'This address is a smart contract',
    });
  }

  if (name.isErc20) {
    chipItems.push({
      id: 'erc20',
      label: 'ERC20',
      color: 'green',
      tooltip: 'This contract implements the ERC-20 token standard',
    });
  }

  if (name.isErc721) {
    chipItems.push({
      id: 'erc721',
      label: 'NFT',
      color: 'purple',
      tooltip: 'This contract implements the ERC-721 NFT standard',
    });
  }

  if (name.isPrefund) {
    chipItems.push({
      id: 'prefund',
      label: 'Prefund',
      color: 'orange',
      tooltip: 'This address was allocated ETH in the genesis block',
    });
  }

  return chipItems;
}
