import { useIcons } from '@hooks';
import { ActionIcon } from '@mantine/core';

export const ToggleButton = ({
  collapsed,
  onToggle,
  direction,
}: {
  collapsed: boolean;
  onToggle: () => void;
  direction: 'left' | 'right';
}) => {
  const { ChevronLeft, ChevronRight } = useIcons();
  const icon =
    direction === 'left' ? (
      collapsed ? (
        <ChevronRight />
      ) : (
        <ChevronLeft />
      )
    ) : collapsed ? (
      <ChevronLeft />
    ) : (
      <ChevronRight />
    );

  const alignment =
    direction === 'left'
      ? collapsed
        ? 'flex-end'
        : 'flex-start'
      : collapsed
        ? 'flex-start'
        : 'flex-end';

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: alignment,
      }}
    >
      <ActionIcon
        onClick={onToggle}
        variant="subtle"
        size="sm"
        radius="md"
        style={{
          fontWeight: 'normal',
          paddingRight: direction === 'right' ? '0.5rem' : '',
          paddingLeft: direction !== 'right' ? '0.5rem' : '',
        }}
      >
        {icon}
      </ActionIcon>
    </div>
  );
};
