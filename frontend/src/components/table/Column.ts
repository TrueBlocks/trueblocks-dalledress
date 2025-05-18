import { ChangeEvent, ReactNode } from 'react';

export interface Column<T extends Record<string, unknown>> {
  name?: string;
  key: string;
  header: string;
  value?: string | number | boolean;
  label?: string;
  placeholder?: string;
  required?: boolean;
  error?: string;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  onBlur?: () => void;
  rightSection?: ReactNode;
  hint?: string;
  visible?: boolean | ((formData: Record<string, unknown>) => boolean);
  objType?: string;
  type?:
    | 'text'
    | 'number'
    | 'password'
    | 'checkbox'
    | 'radio'
    | 'button'
    | 'textarea'
    | 'select';
  fields?: Column<T>[];
  isButtonGroup?: boolean;
  buttonAlignment?: 'left' | 'center' | 'right';
  customRender?: ReactNode;
  readOnly?: boolean;
  disabled?: boolean;
  sameLine?: boolean;
  flex?: number;
  editable?: boolean;
  width?: string | number;
  className?: string;
  sortable?: boolean;
  accessor?: (row: T) => ReactNode;
  render?: (row: T, rowIndex: number) => ReactNode;
}
