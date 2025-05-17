import { ChangeEvent, ReactNode, useCallback, useMemo } from 'react';

import { TextInput } from '@mantine/core';

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
  onBlur?: () => void;
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
  width?: string | number;
  className?: string;
  sortable?: boolean;
  accessor?: (row: T) => ReactNode;
  render?: (row: T, rowIndex: number) => ReactNode;
}

export const usePreprocessedFields = <T,>(
  fields: FormField<T>[],
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void,
  formData: T = {} as T,
): FormField<T>[] => {
  const combineOnOneLine = useCallback(
    (fields: (FormField<T> & { flex?: number })[]): FormField<T> => ({
      customRender: (
        <div style={{ display: 'flex', gap: '1rem' }}>
          {fields.map((field) => (
            <TextInput
              key={field.name}
              label={field.label}
              placeholder={field.placeholder}
              withAsterisk={field.required}
              onChange={(e) => {
                if (!field.readOnly) {
                  field.onChange?.(e);
                }
              }}
              error={field.error}
              rightSection={field.rightSection}
              onBlur={field.onBlur}
              name={field.name}
              readOnly={field.readOnly}
              style={{
                flex: field.flex || 1,
              }}
            />
          ))}
        </div>
      ),
    }),
    [],
  );

  const preprocessFields = useCallback(
    (fieldsToProcess: FormField<T>[]): FormField<T>[] => {
      const acceptedFields: FormField<T>[] = [];
      let currentGroup: FormField<T>[] = [];

      fieldsToProcess.forEach((field, index) => {
        if (typeof field.editable === 'undefined') {
          field.editable = true;
        }
        if (field.editable === false) {
          return;
        }

        const isVisible =
          typeof field.visible === 'function'
            ? field.visible(formData)
            : field.visible !== false;

        if (!isVisible) {
          return;
        }

        if (field.fields && field.fields.length > 0) {
          field.fields = preprocessFields(field.fields);
          acceptedFields.push(field);
          return;
        }

        if (currentGroup.length === 0 || field.sameLine) {
          currentGroup.push(field);
        }

        const nextField = fieldsToProcess[index + 1];
        if (!nextField || !nextField.sameLine) {
          if (currentGroup.length > 1) {
            acceptedFields.push(combineOnOneLine(currentGroup));
          } else if (currentGroup[0]) {
            acceptedFields.push(currentGroup[0]);
          }
          currentGroup = [];
        }
      });

      return acceptedFields;
    },
    [combineOnOneLine, formData],
  );

  return useMemo(() => preprocessFields(fields), [preprocessFields, fields]);
};
