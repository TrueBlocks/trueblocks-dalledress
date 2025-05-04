import { TabView } from '@layout';
import { InfoView } from '@layout';
import { FormView } from '@layout';

export const Names = () => {
  const tabs = [
    {
      label: 'Custom',
      content: <InfoView title="Custom Info" />,
    },
    {
      label: 'Named',
      content: (
        <FormView
          title="Named Form"
          formFields={[
            {
              name: 'name',
              label: 'Name',
              placeholder: 'Enter a name',
              required: true,
            },
            {
              name: 'description',
              label: 'Description',
              placeholder: 'Enter a description',
            },
          ]}
          onSubmit={(data) => {
            console.log('Form submitted:', data);
          }}
        />
      ),
    },
    {
      label: 'Prefund',
      content: <div>This is the Prefund tab.</div>,
    },
  ];

  return <TabView tabs={tabs} route="/names" />;
};
