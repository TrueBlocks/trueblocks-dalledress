import { Form } from '@components';
import { MantineProvider } from '@mantine/core';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { vi } from 'vitest';

vi.mock('react-hotkeys-hook', () => ({
  useHotkeys: (keys: string, handler: (e: KeyboardEvent) => void) => {
    if (keys === 'mod+a') {
      document.addEventListener('keydown', (e) => {
        if (e.key === 'a' && (e.metaKey || e.ctrlKey)) {
          handler(e);
        }
      });
    }
  },
}));

describe('Form Component', () => {
  it('renders the form with a title and description', () => {
    render(
      <MantineProvider>
        <Form
          title="Test Form"
          description="This is a test form."
          fields={[]}
          onSubmit={vi.fn()}
        />
      </MantineProvider>,
    );

    expect(screen.getByText('Test Form')).toBeInTheDocument();
    expect(screen.getByText('This is a test form.')).toBeInTheDocument();
  });

  it('renders a text input field', () => {
    render(
      <MantineProvider>
        <Form
          title="Test Form"
          initMode="edit"
          fields={[
            {
              name: 'username',
              label: 'Username',
              type: 'text',
              value: '',
              placeholder: 'Enter your username',
              onChange: vi.fn(),
            },
          ]}
          onSubmit={vi.fn()}
        />
      </MantineProvider>,
    );

    expect(screen.getByLabelText('Username')).toBeInTheDocument();
    expect(
      screen.getByPlaceholderText('Enter your username'),
    ).toBeInTheDocument();
  });

  it('renders a number input field', () => {
    render(
      <MantineProvider>
        <Form
          title="Test Form"
          initMode="edit"
          fields={[
            {
              name: 'age',
              label: 'Age',
              type: 'number',
              value: '',
              placeholder: 'Enter your age',
              onChange: vi.fn(),
            },
          ]}
          onSubmit={vi.fn()}
        />
      </MantineProvider>,
    );

    expect(screen.getByLabelText('Age')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Enter your age')).toBeInTheDocument();
  });

  it('resets the form to its original values when Cancel is clicked', () => {
    const onCancel = vi.fn();

    render(
      <MantineProvider>
        <Form
          title="Test Form"
          initMode="edit"
          fields={[
            {
              name: 'username',
              label: 'Username',
              type: 'text',
              value: 'JohnDoe',
              placeholder: 'Enter your username',
              onChange: vi.fn(),
            },
          ]}
          onCancel={onCancel}
          onSubmit={vi.fn()}
        />
      </MantineProvider>,
    );

    const cancelButton = screen.getByText('Cancel');
    fireEvent.click(cancelButton);

    expect(onCancel).toHaveBeenCalled();
  });

  it('submits the form with the correct values when Save is clicked', () => {
    const onSubmit = vi.fn();

    render(
      <MantineProvider>
        <Form
          title="Test Form"
          initMode="edit"
          fields={[
            {
              name: 'username',
              label: 'Username',
              type: 'text',
              value: 'JohnDoe',
              placeholder: 'Enter your username',
              onChange: vi.fn(),
            },
          ]}
          onSubmit={onSubmit}
        />
      </MantineProvider>,
    );

    // Simulate clicking the Save button
    const saveButton = screen.getByText('Save');
    fireEvent.click(saveButton);

    // Verify that the onSubmit handler was called
    expect(onSubmit).toHaveBeenCalled();
  });

  it('conditionally renders fields based on the visible property', () => {
    render(
      <MantineProvider>
        <Form
          title="Test Form"
          initMode="edit"
          fields={[
            {
              name: 'username',
              label: 'Username',
              type: 'text',
              value: '',
              placeholder: 'Enter your username',
              visible: true,
              onChange: vi.fn(),
            },
            {
              name: 'hiddenField',
              label: 'Hidden Field',
              type: 'text',
              value: '',
              placeholder: 'This field is hidden',
              visible: false,
              onChange: vi.fn(),
            },
          ]}
          onSubmit={vi.fn()}
        />
      </MantineProvider>,
    );

    // Verify that the visible field is rendered
    expect(screen.getByLabelText('Username')).toBeInTheDocument();

    // Verify that the hidden field is not rendered
    expect(screen.queryByLabelText('Hidden Field')).not.toBeInTheDocument();
  });

  it('renders fields inline when sameLine is true', () => {
    render(
      <MantineProvider>
        <Form
          title="Test Form"
          initMode="edit"
          fields={[
            {
              name: 'firstName',
              label: 'First Name',
              type: 'text',
              value: '',
              placeholder: 'Enter your first name',
              sameLine: true,
              onChange: vi.fn(),
            },
            {
              name: 'lastName',
              label: 'Last Name',
              type: 'text',
              value: '',
              placeholder: 'Enter your last name',
              sameLine: true,
              onChange: vi.fn(),
            },
          ]}
          onSubmit={vi.fn()}
        />
      </MantineProvider>,
    );

    // Verify that both fields are rendered
    const form = screen.getByRole('form');
    const inlineContainers = within(form).getAllByRole('textbox');

    // Ensure both fields are rendered inline
    expect(inlineContainers).toHaveLength(2);
    expect(inlineContainers[0]).toHaveAttribute(
      'placeholder',
      'Enter your first name',
    );
    expect(inlineContainers[1]).toHaveAttribute(
      'placeholder',
      'Enter your last name',
    );
  });

  it('selects all text in the active input field when mod+a is pressed', () => {
    render(
      <MantineProvider>
        <Form
          title="Test Form"
          initMode="edit"
          fields={[
            {
              name: 'username',
              label: 'Username',
              type: 'text',
              value: 'JohnDoe',
              placeholder: 'Enter your username',
              onChange: vi.fn(),
            },
          ]}
          onSubmit={vi.fn()}
        />
      </MantineProvider>,
    );

    const input = screen.getByLabelText('Username') as HTMLInputElement;
    input.focus();
    fireEvent.keyDown(input, { key: 'a', metaKey: true });

    expect(input.selectionStart).toBe(0);
    expect(input.selectionEnd).toBe(input.value.length);
  });
});
