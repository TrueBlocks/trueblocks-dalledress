import { ChangeEvent, ReactNode, useCallback, useMemo } from 'react';

import { TextInput } from '@mantine/core';

export interface FormField {
  name?: string;
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
  fields?: FormField[];
  isButtonGroup?: boolean;
  buttonAlignment?: 'left' | 'center' | 'right';
  customRender?: ReactNode;
  readOnly?: boolean;
  disabled?: boolean;
  sameLine?: boolean;
  flex?: number;
}

export const usePreprocessedFields = (
  fields: FormField[],
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void,
  formData: Record<string, unknown> = {},
): FormField[] => {
  const combineOnOneLine = useCallback(
    (fields: (FormField & { flex?: number })[]): FormField => ({
      customRender: (
        <div style={{ display: 'flex', gap: '1rem' }}>
          {fields.map((field) => (
            <TextInput
              key={field.name}
              label={field.label}
              placeholder={field.placeholder}
              withAsterisk={field.required}
              value={field.value as string}
              onChange={field.readOnly ? undefined : field.onChange || onChange}
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
    [onChange],
  );

  const preprocessFields = useCallback(
    (fieldsToProcess: FormField[]): FormField[] => {
      const acceptedFields: FormField[] = [];
      let currentGroup: FormField[] = [];

      fieldsToProcess.forEach((field, index) => {
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

export const createFieldRenderer = (
  field: FormField,
  index: number,
  renderComponent: (
    field: FormField,
    keyProp: string | number,
    autoFocus?: boolean,
  ) => ReactNode,
) => {
  return renderComponent(field, field.name || index, index === 0);
};
