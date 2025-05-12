export interface TableKey {
  viewName: string;
  tabName: string;
  tableId?: string;
}

export const tableKeyToString = (key: TableKey): string => {
  return `${key.viewName}/${key.tabName}/${key.tableId ?? ''}`;
};
