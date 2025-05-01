import { ChangeEvent, useEffect } from 'react';

import { GetUserPreferences } from '@app';
import { FormField, WizardForm } from '@components';
import { msgs } from '@models';
import { emitEvent } from '@utils';

import { WizardStepProps } from '.';

export const StepUserInfo = ({
  state,
  onSubmit,
  updateData,
  validateName,
  validateEmail,
  onCancel,
}: WizardStepProps) => {
  const { name, email } = state.data;
  const { nameError, emailError } = state.validation;

  useEffect(() => {
    const loadUserData = async () => {
      try {
        const userPrefs = await GetUserPreferences();
        if (updateData) {
          updateData({
            name: userPrefs.name || '',
            email: userPrefs.email || '',
          });
        }
      } catch (error) {
        emitEvent(
          msgs.EventType.STATUS,
          `Error trying to load user info: ${error}`,
        );
      }
    };

    if (!name && !email) {
      loadUserData();
    }
  }, [name, email, updateData]);

  const handleNameChange = (e: ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    if (updateData) {
      updateData({ name: newValue });
    }
  };

  const handleEmailChange = (e: ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    if (updateData) {
      updateData({ email: newValue });
    }
  };

  const formFields: FormField[] = [
    {
      name: 'name',
      value: name,
      label: 'Name',
      placeholder: 'Enter your name',
      required: true,
      error: nameError,
      onChange: handleNameChange,
      onBlur: validateName,
    },
    {
      name: 'email',
      value: email,
      label: 'Email',
      placeholder: 'Enter your email',
      required: true,
      error: emailError,
      onChange: handleEmailChange,
      onBlur: validateEmail,
    },
  ];

  return (
    <WizardForm
      title="User Information"
      description="Please provide your name and email address."
      fields={formFields}
      onSubmit={onSubmit}
      onCancel={onCancel}
    />
  );
};
