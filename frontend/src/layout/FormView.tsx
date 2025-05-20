import { Logger } from '@app';
import { Form, FormField, FormProps } from '@components';
import { Container } from '@mantine/core';

export interface FormViewProps<T> {
  formFields: FormField<T>[];
  title?: string;
  description?: string;
  onSubmit: FormProps<T>['onSubmit'];
  onChange?: FormProps<T>['onChange'];
  onCancel?: FormProps<T>['onCancel'];
}

export const FormView = <T extends Record<string, unknown>>({
  formFields,
  title,
  description,
  onSubmit,
  onChange,
  onCancel,
}: FormViewProps<T>) => {
  Logger('DEBUGGING1 FormView: ' + JSON.stringify(formFields, null, 2));
  return (
    <Container
      size="md"
      mt="xl"
      style={{
        backgroundColor: '#1a1a1a',
        padding: '1rem',
        borderRadius: '8px',
      }}
    >
      <Form<T>
        title={title}
        description={description}
        fields={formFields}
        onSubmit={onSubmit}
        onChange={onChange}
        onCancel={onCancel}
        submitText="Save"
      />
    </Container>
  );
};
