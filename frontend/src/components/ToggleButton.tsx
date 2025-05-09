import { ActionIcon } from '@mantine/core';
import { FaAngleDoubleLeft, FaAngleDoubleRight } from 'react-icons/fa';

export const ToggleButton = ({
  collapsed,
  onToggle,
  direction,
}: {
  collapsed: boolean;
  onToggle: () => void;
  direction: 'left' | 'right';
}) => {
  const icon =
    direction === 'left' ? (
      collapsed ? (
        <FaAngleDoubleRight />
      ) : (
        <FaAngleDoubleLeft />
      )
    ) : collapsed ? (
      <FaAngleDoubleLeft />
    ) : (
      <FaAngleDoubleRight />
    );

  const alignment =
    direction === 'left'
      ? collapsed
        ? 'flex-start'
        : 'flex-end'
      : collapsed
        ? 'flex-end'
        : 'flex-start';

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
        }}
      >
        {icon}
      </ActionIcon>
    </div>
  );
};
