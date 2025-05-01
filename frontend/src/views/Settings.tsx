import { SettingsApp, SettingsOrg, SettingsUser } from '@views';

import { TabView } from '../layout/TabView';

export const Settings = () => {
  const tabs = [
    { label: 'Org', content: <SettingsOrg /> },
    { label: 'User', content: <SettingsUser /> },
    { label: 'App', content: <SettingsApp /> },
  ];

  return <TabView tabs={tabs} route="/settings" />;
};
