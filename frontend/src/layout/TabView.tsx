import { useCallback, useState } from 'react';

import { useAppContext } from '@contexts';
import { useEvent } from '@hooks';
import { Tabs } from '@mantine/core';
import { msgs } from '@models';

interface Tab {
  label: string;
  content: React.ReactNode;
}

interface TabViewProps {
  tabs: Tab[];
  route: string;
}

export const TabView = ({ tabs, route }: TabViewProps) => {
  const { lastTab, setLastTab } = useAppContext();
  const defaultTab = lastTab[route] || tabs[0]?.label || '';
  const [activeTab, setActiveTab] = useState<string>(defaultTab);

  const nextTab = useCallback((): string => {
    const currentIndex = tabs.findIndex((tab) => tab.label === activeTab);
    const nextIndex = (currentIndex + 1) % tabs.length;
    return tabs[nextIndex]?.label || activeTab;
  }, [tabs, activeTab]);

  const prevTab = useCallback((): string => {
    const currentIndex = tabs.findIndex((tab) => tab.label === activeTab);
    const prevIndex = (currentIndex - 1 + tabs.length) % tabs.length;
    return tabs[prevIndex]?.label || activeTab;
  }, [tabs, activeTab]);

  useEvent<{ route: string; key: string }>(
    msgs.EventType.TAB_CYCLE,
    (event) => {
      if (event.route === route) {
        const newTab = event.key.startsWith('alt+') ? prevTab() : nextTab();
        setActiveTab(newTab);
        setLastTab(route, newTab); // Synchronize with AppContext
      }
    },
  );

  return (
    <div>
      <Tabs
        value={activeTab}
        onChange={(value) => {
          setActiveTab(value || tabs[0]?.label || '');
          setLastTab(route, value || tabs[0]?.label || '');
        }}
      >
        <Tabs.List>
          {tabs.map((tab) => (
            <Tabs.Tab key={tab.label} value={tab.label}>
              {tab.label}
            </Tabs.Tab>
          ))}
        </Tabs.List>

        {tabs.map((tab) => (
          <Tabs.Panel key={tab.label} value={tab.label}>
            {tab.content}
          </Tabs.Panel>
        ))}
      </Tabs>
    </div>
  );
};
