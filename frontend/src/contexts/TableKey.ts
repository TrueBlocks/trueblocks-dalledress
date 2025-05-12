export interface TableKey {
  viewName: string;
  tabName: string;
}

export const tableKeyToString = (key: TableKey): string => {
  return `${key.viewName}/${key.tabName}`;
};
