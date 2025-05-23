import { ChangeEvent, FormEvent, useEffect, useState } from 'react';

import { GetOrgPreferences, Logger, SetOrgPreferences } from '@app';
import { FormField } from '@components';
import { FormView } from '@layout';
import { msgs, preferences } from '@models';
import { emitEvent } from '@utils';

type IndexableOrg = preferences.OrgPreferences & { [key: string]: unknown };

export const SettingsOrg = () => {
  const [formData, setFormData] = useState<IndexableOrg>({});
  const [originalData, setOriginalData] = useState<IndexableOrg>({});

  useEffect(() => {
    GetOrgPreferences().then((data) => {
      setFormData(data as IndexableOrg);
      setOriginalData({ ...data });
    });
  }, []);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    SetOrgPreferences(formData).then(() => {
      setOriginalData({ ...formData });
    });
  };

  const handleCancel = () => {
    setFormData({ ...originalData });
    emitEvent(msgs.EventType.STATUS, `Changes were discarded`);
  };

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: value,
    }));
    Logger(`${name} = ${JSON.stringify(value)}`);
  };

  const formFields: FormField<Record<string, unknown>>[] = [
    {
      name: 'developerName',
      value: formData.developerName || '',
      label: 'Organization Name',
      placeholder: 'Enter your organization name',
      required: true,
    },
    {
      name: 'supportUrl',
      value: formData.supportUrl || '',
      label: 'Support',
      placeholder: 'Enter your organization support url',
      required: true,
    },
    {
      label: 'Options',
      fields: [
        {
          name: 'language',
          value: formData.language || '',
          label: 'Language',
          placeholder: 'Enter your language',
          required: true,
        },
        {
          name: 'theme',
          value: formData.theme || '',
          label: 'Theme',
          placeholder: 'Enter your theme',
          required: true,
          sameLine: true,
        },
        {
          name: 'telemetry',
          value: formData.telemetry ? 'true' : 'false',
          label: 'Telemetry',
          placeholder: 'Enter your telemetry',
          sameLine: true,
        },
        {
          name: 'logLevel',
          value: formData.logLevel || '',
          label: 'Log Level',
          placeholder: 'Enter your log level',
          required: true,
        },
        {
          name: 'experimental',
          value: formData.experimental ? 'true' : 'false',
          label: 'Experimental',
          placeholder: 'Enter your experimental',
          sameLine: true,
        },
        {
          name: 'version',
          value: formData.version || '',
          label: 'Version',
          placeholder: 'Enter your version',
          required: true,
          readOnly: true,
          sameLine: true,
        },
      ],
    },
  ];

  return (
    <FormView<IndexableOrg>
      title="Edit / Manage Your Settings"
      formFields={formFields}
      onSubmit={handleSubmit}
      onChange={handleChange}
      onCancel={handleCancel}
    />
  );
};
