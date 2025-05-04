import { ChangeEvent, forwardRef } from 'react';

import { FormField } from '@components';
import { Fieldset, Stack, Text, TextInput } from '@mantine/core';

export interface FieldRendererProps {
  field: FormField;
  mode?: 'display' | 'edit';
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  loading?: boolean;
  keyProp?: string | number;
  autoFocus?: boolean;
}

export const FieldRenderer = forwardRef<HTMLInputElement, FieldRendererProps>(
  ({ field, mode, onChange, loading, keyProp, autoFocus }, ref) => {
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
          onChange={field.readOnly ? undefined : field.onChange || onChange}
          error={
            (!loading && field.error) ||
            (field.required && !field.value && `${field.label} is required`)
          }
          rightSection={field.rightSection}
          onBlur={field.onBlur}
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
