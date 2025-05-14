import { ChangeEvent, FormEvent, useEffect, useState } from 'react';

import { FieldRenderer, FormField, usePreprocessedFields } from '@components';
import { useFormHotkeys } from '@hooks';
import { Button, Group, Stack, Text, Title } from '@mantine/core';

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
}

export function Form<T = Record<string, unknown>>({
  title,
  description,
  fields,
  onSubmit,
  onCancel,
  onChange,
  submitButtonRef,
  initMode = 'display',
  compact = false,
}: FormProps<T>) {
  const [loading, setLoading] = useState(true);
  const [mode, setMode] = useState<'display' | 'edit'>(initMode);

  const handleEdit = () => {
    setMode('edit');
  };

  const handleSave = (e: FormEvent) => {
    e?.preventDefault();
    setMode('display');
    onSubmit(e);
  };

  const handleCancel = () => {
    setMode('display');
    if (onCancel) onCancel();
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

  const processedFields = usePreprocessedFields<T>(fields, onChange);

  const renderField = (field: FormField<T>, index: number) => (
    <FieldRenderer
      key={field.name || index}
      field={field as FormField<Record<string, unknown>>}
      mode={mode}
      onChange={onChange}
      loading={loading}
      keyProp={field.name || index}
      autoFocus={index === 0}
    />
  );

  return (
    <Stack gap={compact ? 'xs' : 'md'}>
      {title && <Title order={3}>{title}</Title>}
      {description && <Text>{description}</Text>}
      <form role="form" onSubmit={handleSave}>
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
}
