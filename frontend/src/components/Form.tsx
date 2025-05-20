import { ChangeEvent, FormEvent, useEffect, useState } from 'react';

import { Logger } from '@app';
import { FieldRenderer, FormField, usePreprocessedFields } from '@components';
import { useFormHotkeys } from '@hooks';
import { Button, Group, Stack, Text, Title } from '@mantine/core';
import { useForm } from '@mantine/form';

export interface FormProps<T = Record<string, unknown>> {
  title?: string;
  description?: string;
  fields: FormField<T>[];
  onSubmit: (e: FormEvent) => void;
  onBack?: () => void;
  onCancel?: () => void;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  submitText?: string;
  submitButtonRef?: React.RefObject<HTMLButtonElement | null>;
  initMode?: 'display' | 'edit';
  compact?: boolean;
  validate?: Record<
    string,
    (value: unknown, values: Record<string, unknown>) => string | null
  >;
}

export const Form = <T = Record<string, unknown>,>({
  title,
  description,
  fields,
  onSubmit,
  onCancel,
  onChange,
  submitButtonRef,
  initMode = 'display',
  compact = false,
  validate = {},
}: FormProps<T>) => {
  const [loading, setLoading] = useState(true);
  const [mode, setMode] = useState<'display' | 'edit'>(initMode);

  // Initialize Mantine form
  const form = useForm({
    initialValues: {}, // Start with an empty object
    validate,
  });

  // Initialize form values from fields
  useEffect(() => {
    // Build values object from fields
    const values: Record<string, unknown> = {};
    fields.forEach((field) => {
      if (field.name && field.value !== undefined) {
        values[field.name] = field.value;
      }
    });

    try {
      form.setValues(values);
    } catch (e) {
      Logger('Error setting form values:' + e);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Empty dependency array - only run once on mount

  const handleEdit = () => {
    setMode('edit');
  };

  const handleFormSubmit = (e: FormEvent) => {
    Logger('DEBUGGING: onSubmit in Form ' + e);
    e.preventDefault();

    // Validate form before submission
    const { hasErrors } = form.validate();

    // Only change mode to display if validation passes
    if (!hasErrors) {
      setMode('display');

      // Pass the Mantine form values to the onSubmit handler
      // This ensures we're using the values from Mantine's state
      onSubmit(e);
    }
  };

  const handleCancel = () => {
    setMode('display');
    if (onCancel) onCancel();
  };

  // Handle field changes - using Mantine's form state management exclusively
  const handleFieldChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    const fieldValue = type === 'checkbox' ? checked : value;

    // Update only Mantine form values - this is now our single source of truth
    if (name) {
      form.setFieldValue(name, fieldValue);
    }

    // Call parent onChange if provided (for backward compatibility)
    if (onChange) {
      onChange(e);
    }
  };

  useEffect(() => {
    setLoading(false);
  }, []);

  useFormHotkeys({
    mode,
    setMode,
    onCancel,
    submitButtonRef,
  });

  const processedFields = usePreprocessedFields<T>(fields);

  const renderField = (field: FormField<T>, index: number) => {
    // Create a new field with error information from Mantine form
    const fieldWithError = {
      ...field,
      // Use Mantine form values instead of the field's own value
      value:
        field.name && field.name in form.values
          ? (form.values as Record<string, unknown>)[field.name]
          : field.value,
      error: field.name ? form.errors[field.name] : undefined,
    };

    return (
      <FieldRenderer
        key={field.name || index}
        field={fieldWithError as FormField<Record<string, unknown>>}
        mode={mode}
        onChange={handleFieldChange}
        loading={loading}
        keyProp={field.name || index}
        autoFocus={index === 0}
      />
    );
  };

  return (
    <Stack gap={compact ? 'xs' : 'md'}>
      {title && <Title order={3}>{title}</Title>}
      {description && <Text>{description}</Text>}
      <form role="form" onSubmit={handleFormSubmit}>
        <Stack gap={compact ? 'xs' : 'md'}>
          {processedFields.map((field, index) => renderField(field, index))}
          <Group justify="flex-end" mt={compact ? 'xs' : 'md'}>
            {mode === 'display' && (
              <Button tabIndex={0} variant="outline" onClick={handleEdit}>
                Edit
              </Button>
            )}
            {mode === 'edit' && (
              <>
                <Button tabIndex={0} variant="outline" onClick={handleCancel}>
                  Cancel
                </Button>
                <Button
                  type="submit"
                  tabIndex={0}
                  ref={submitButtonRef as React.RefObject<HTMLButtonElement>}
                >
                  Save
                </Button>
              </>
            )}
          </Group>
        </Stack>
      </form>
    </Stack>
  );
};
