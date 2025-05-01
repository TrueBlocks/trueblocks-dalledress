import { ChangeEvent, FormEvent, useEffect, useState } from 'react';

import { FieldRenderer, FormField, usePreprocessedFields } from '@components';
import { useFormHotkeys } from '@hooks';
import { Button, Group, Stack, Text, Title } from '@mantine/core';

export interface FormProps {
  title?: string;
  description?: string;
  fields: FormField[];
  onSubmit: (e: FormEvent) => void;
  onBack?: () => void;
  onCancel?: () => void;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  submitText?: string;
  submitButtonRef?: React.RefObject<HTMLButtonElement | null>;
  initMode?: 'display' | 'edit';
}

export const Form = ({
  title,
  description,
  fields,
  onSubmit,
  onCancel,
  onChange,
  submitButtonRef,
  initMode = 'display',
}: FormProps) => {
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

  const processedFields = usePreprocessedFields(fields, onChange);

  const renderField = (field: FormField, index: number) => (
    <FieldRenderer
      key={field.name || index}
      field={field}
      mode={mode}
      onChange={onChange}
      loading={loading}
      keyProp={field.name || index}
      autoFocus={index === 0}
    />
  );

  return (
    <Stack>
      {title && <Title order={3}>{title}</Title>}
      {description && <Text>{description}</Text>}
      <form role="form" onSubmit={handleSave}>
        <Stack>
          {processedFields.map((field, index) => renderField(field, index))}
          <Group justify="flex-end" mt="md">
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
