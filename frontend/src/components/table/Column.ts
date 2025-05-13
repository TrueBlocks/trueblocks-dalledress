export interface Column<T> {
  key: string;
  header: string;
  render?: (row: T, rowIndex: number) => React.ReactNode;
  accessor?: (row: T) => React.ReactNode;
  width?: string | number;
  className?: string;
  sortable?: boolean;
  editable?: boolean;
}
