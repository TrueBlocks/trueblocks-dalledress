export const getBarWidth = (collapsed: boolean, factor: number) =>
  collapsed ? 50 : 150 * factor;
