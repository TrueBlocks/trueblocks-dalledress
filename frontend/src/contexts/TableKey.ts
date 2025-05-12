export interface TableKey {
  viewName: string;
  tabName: string;
  tableId?: number;
}

export const tableKeyToString = (key: TableKey): string => {
  return `${key.viewName}/${key.tabName}/${key.tableId ?? 0}`;
};
