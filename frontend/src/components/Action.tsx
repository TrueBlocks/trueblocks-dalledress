import { useIcons } from '@hooks';
import { ActionIcon, ActionIconProps } from '@mantine/core';

type IconName = keyof ReturnType<typeof useIcons>;

interface ActionProps extends Omit<ActionIconProps, 'children' | 'onClick'> {
  icon: IconName;
  iconOff?: IconName;
  isOn?: boolean;
  onClick: () => void;
  disabled?: boolean;
  isSubdued?: boolean;
  title?: string;
}

export const Action = ({
  icon,
  iconOff,
  isOn = true,
  onClick,
  disabled = false,
  isSubdued = false,
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

  const style: React.CSSProperties = {};
  if (isSubdued && !disabled) {
    style.opacity = 0.6;
  }

  return (
    <ActionIcon
      onClick={handleClick}
      disabled={disabled}
      title={title}
      style={style}
      {...mantineProps}
    >
      <IconComponent />
    </ActionIcon>
  );
};
