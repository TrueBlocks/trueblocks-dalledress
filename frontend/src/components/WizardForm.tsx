import { ChangeEvent, FormEvent, useEffect, useState } from 'react';

import { FieldRenderer, FormField, usePreprocessedFields } from '@components';
import { useFormHotkeys } from '@hooks';
import { Button, Group, Stack, Text, Title } from '@mantine/core';

export interface WizardFormProps {
  title?: string;
  description?: string;
  fields: FormField[];
  onSubmit: (e: FormEvent) => void;
  onBack?: () => void;
  onCancel?: () => void;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  submitText?: string;
  submitButtonRef?: React.RefObject<HTMLButtonElement | null>;
}

export const WizardForm = ({
  title,
  description,
  fields,
  onSubmit,
  onBack,
  onCancel,
  onChange,
  submitText = 'Next',
  submitButtonRef,
}: WizardFormProps) => {
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(false);
  }, []);

  useFormHotkeys({
    onCancel,
    submitButtonRef,
  });

  const processedFields = usePreprocessedFields(fields, onChange);

  const renderField = (field: FormField, index: number) => (
    <FieldRenderer
      key={field.name || index}
      field={field}
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
      <form role="form" onSubmit={onSubmit}>
        <Stack>
          {processedFields.map((field, index) => renderField(field, index))}
          <Group justify="flex-end" mt="md">
            {onBack && (
              <Button tabIndex={0} variant="outline" onClick={onBack}>
                Back
              </Button>
            )}
            <Button
              type="submit"
              tabIndex={0}
              ref={submitButtonRef as React.RefObject<HTMLButtonElement>}
            >
              {submitText}
            </Button>
          </Group>
        </Stack>
      </form>
    </Stack>
  );
};
