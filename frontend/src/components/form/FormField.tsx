import { ChangeEvent, ReactNode } from 'react';

export type ColumnSizeClass =
  // Base semantic sizes (set initial width)
  | 'col-base-address' // 340px for addresses
  | 'col-base-encoding' // 200px for function encodings
  | 'col-base-signature' // 300px for function signatures
  | 'col-base-date' // 140px for dates
  | 'col-base-actions' // 120px for action buttons
  | 'col-base-sm' // 100px small
  | 'col-base-md' // 120px medium
  | 'col-base-lg' // 160px large
  | 'col-base-xl' // 200px extra large
  // Behavior modifiers (override base widths)
  | 'col-min' // min-content (shrink to content)
  | 'col-max' // max-content (expand to content)
  | 'col-fit' // fit-content (fit to content)
  | 'col-expand' // auto (expand to fill space)
  // Backward compatibility (old single-class system)
  | 'col-address' // 340px addresses
  | 'col-encoding' // 200px encodings
  | 'col-signature' // 300px signatures
  | 'col-date' // 140px dates
  | 'col-actions' // 120px actions
  | 'col-content-sm' // 100px small
  | 'col-content-md' // 120px medium
  | 'col-content-lg' // 160px large
  | 'col-content-xl' // 200px extra large
  | 'col-content'; // min-content

export interface FormField<T = Record<string, unknown>> {
  name?: string;
  key?: string;
  header?: string;
  value?: string | number | boolean;
  label?: string;
  placeholder?: string;
  required?: boolean;
  error?: string;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void;
  rightSection?: ReactNode;
  hint?: string;
  visible?: boolean | ((formData: T) => boolean);
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
  fields?: FormField<T>[];
  isButtonGroup?: boolean;
  buttonAlignment?: 'left' | 'center' | 'right';
  customRender?: ReactNode;
  readOnly?: boolean;
  disabled?: boolean;
  sameLine?: boolean;
  flex?: number;
  editable?: boolean;
  width?: string | number | ColumnSizeClass;
  className?: string;
  sortable?: boolean;
  accessor?: (row: T) => ReactNode;
  render?: (row: T, rowIndex: number) => ReactNode;
}
