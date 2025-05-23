import { ChangeEvent, forwardRef } from 'react';

import { FormField } from '@components';
import { Fieldset, Stack, Text, TextInput } from '@mantine/core';

export interface FieldRendererProps {
  field: FormField<Record<string, unknown>>;
  mode?: 'display' | 'edit';
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void;
  loading?: boolean;
  keyProp?: string | number;
  autoFocus?: boolean;
}

export const FieldRenderer = forwardRef<HTMLInputElement, FieldRendererProps>(
  ({ field, mode, onChange, onBlur, loading, keyProp, autoFocus }, ref) => {
    if (field.fields && field.fields.length > 0) {
      return (
        <Fieldset key={keyProp}>
          {field.label && <legend>{field.label}</legend>}
          <Stack>
            {field.fields.map((nestedField, nestedIndex) => (
              <FieldRenderer
                key={nestedField.name || nestedIndex}
                field={nestedField}
                mode={mode}
                onChange={onChange}
                onBlur={onBlur}
                loading={loading}
              />
            ))}
          </Stack>
        </Fieldset>
      );
    }

    if (field.customRender) {
      return <div key={keyProp}>{field.customRender}</div>;
    }

    if (mode === 'display') {
      return (
        <div key={keyProp}>
          <Text size="sm" fw={500}>
            {field.label}: {field.value || 'N/A'}
          </Text>
        </div>
      );
    }

    return (
      <div key={keyProp}>
        <TextInput
          ref={ref}
          label={field.label}
          placeholder={field.placeholder}
          withAsterisk={field.required}
          value={field.value as string}
          onChange={(e) => {
            if (!field.readOnly) {
              field.onChange?.(e);
            }
            if (onChange) {
              onChange(e);
            }
          }}
          onBlur={(e) => {
            field.onBlur?.(e);
            if (onBlur) {
              onBlur(e);
            }
          }}
          error={
            (!loading && field.error) ||
            (field.required && !field.value && `${field.label} is required`)
          }
          styles={{
            input: {
              ...(field.error
                ? {
                    borderColor: '#fa5252',
                    backgroundColor: 'rgba(250, 82, 82, 0.1)',
                  }
                : {}),
              ...(field.readOnly
                ? {
                    color: 'var(--mantine-color-text)', // Use Mantine's text color variable for theme adaptability
                    opacity: 0.6, // Slightly reduce opacity to differentiate but keep readable
                  }
                : {}),
            },
            error: {
              fontWeight: 500,
            },
          }}
          rightSection={field.rightSection}
          name={field.name}
          readOnly={field.readOnly}
          disabled={field.readOnly}
          tabIndex={field.readOnly ? -1 : 0}
          autoFocus={autoFocus}
        />
        {field.hint && (
          <Text size="sm" c="dimmed">
            {field.hint}
          </Text>
        )}
      </div>
    );
  },
);

FieldRenderer.displayName = 'FieldRenderer';
