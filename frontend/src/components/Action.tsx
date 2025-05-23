import { useIcons } from '@hooks';
import { ActionIcon, ActionIconProps } from '@mantine/core';

type IconName = keyof ReturnType<typeof useIcons>;

interface ActionProps extends Omit<ActionIconProps, 'children' | 'onClick'> {
  icon: IconName;
  iconOff?: IconName;
  isOn?: boolean;
  onClick: () => void;
  disabled?: boolean;
  title?: string;
}

export const Action = ({
  icon,
  iconOff,
  isOn = true,
  onClick,
  disabled = false,
  title,
  ...mantineProps
}: ActionProps) => {
  const icons = useIcons();

  const currentIcon = iconOff && !isOn ? iconOff : icon;
  const IconComponent = icons[currentIcon];

  const handleClick = () => {
    if (!disabled) {
      onClick();
    }
  };

  return (
    <ActionIcon
      onClick={handleClick}
      disabled={disabled}
      title={title}
      {...mantineProps}
    >
      <IconComponent />
    </ActionIcon>
  );
};
