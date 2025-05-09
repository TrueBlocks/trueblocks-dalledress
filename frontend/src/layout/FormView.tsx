import { Form, FormField, FormProps } from '@components';
import { Container } from '@mantine/core';

export const FormView = ({
  formFields,
  title,
  description,
  onSubmit,
  onChange,
  onCancel,
}: {
  formFields: FormField[];
  title?: string;
  description?: string;
  onSubmit: FormProps['onSubmit'];
  onChange?: FormProps['onChange'];
  onCancel?: FormProps['onCancel'];
}) => {
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
      <Form
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
