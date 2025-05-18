import { useEffect, useState } from 'react';

import { Column, Form, FormField } from '@components';

import './EditableRowForm.css';

interface EditableRowFormProps<T extends Record<string, unknown>> {
  row: T;
  columns: Column<T>[];
  onClose: () => void;
  onSubmit?: (updatedData: Partial<T>) => void;
}

export const EditableRowForm = <T extends Record<string, unknown>>({
  row,
  columns,
  onClose,
  onSubmit,
}: EditableRowFormProps<T>) => {
  const [formData, setFormData] = useState<Partial<T>>(() => {
    const initial: Partial<T> = {};
    columns.forEach((col) => {
      if (col.accessor) {
        try {
          initial[col.key as keyof T] = col.accessor(row) as T[keyof T];
        } catch {
          initial[col.key as keyof T] = row[col.key as keyof T];
        }
      } else {
        initial[col.key as keyof T] = row[col.key as keyof T];
      }
    });
    return initial;
  });

  // Update formData if row changes
  useEffect(() => {
    const initial: Partial<T> = {};
    columns.forEach((col) => {
      if (col.accessor) {
        try {
          initial[col.key as keyof T] = col.accessor(row) as T[keyof T];
        } catch {
          initial[col.key as keyof T] = row[col.key as keyof T];
        }
      } else {
        initial[col.key as keyof T] = row[col.key as keyof T];
      }
    });
    setFormData(initial);
  }, [row, columns]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (onSubmit) onSubmit(formData);
    onClose();
  };

  const formFields: FormField<T>[] = columns
    .filter((col) => col.editable !== false)
    .map((col) => ({
      name: col.key,
      label: col.header,
      value: String(formData[col.key as keyof T] ?? ''),
      type: 'text',
      readOnly: col.readOnly,
    }));

  return (
    <div className="editable-row-form">
      <Form<T>
        fields={formFields}
        initMode="edit"
        onSubmit={handleSubmit}
        onCancel={onClose}
        onChange={handleChange}
        submitText="Save"
        compact={true}
      />
    </div>
  );
};
