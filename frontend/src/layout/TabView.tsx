import { Fragment, useCallback, useEffect, useState } from 'react';

import { TabDivider } from '@components';
import { useActiveProject, useEvent } from '@hooks';
import { Tabs } from '@mantine/core';
import { msgs, types } from '@models';
import { Log } from '@utils';

import './TabView.css';

interface Tab {
  label: string;
  value: string;
  content: React.ReactNode;
  dividerBefore?: boolean;
}

interface TabViewProps {
  tabs: Tab[];
  route: string;
  onTabChange?: (newTab: string) => void;
}

export const TabView = ({ tabs, route, onTabChange }: TabViewProps) => {
  Log(`[TabView] Component rendered for route: ${route}`);

  const { lastFacetMap, setLastFacet, loading } = useActiveProject();

  const [activeTab, setActiveTab] = useState<string>('');

  useEffect(() => {
    Log(
      `[TabView] useEffect: route=${route}, loading=${loading}, activeTab="${activeTab}"`,
    );

    if (loading) {
      Log('[TabView] Still loading, skipping tab setup');
      return;
    }

    const vR = route.replace(/^\/+/, '');
    const savedTab = lastFacetMap[vR];
    const isValidSavedTab =
      savedTab &&
      String(savedTab) !== 'undefined' &&
      tabs.some((tab) => tab.value === savedTab);
    const targetTab = isValidSavedTab ? savedTab : tabs[0]?.value || '';

    Log(
      `[TabView] savedTab="${savedTab}", isValidSavedTab=${isValidSavedTab}, targetTab="${targetTab}"`,
    );

    if (targetTab && targetTab !== activeTab) {
      Log(`[TabView] Setting activeTab to "${targetTab}"`);
      setActiveTab(targetTab);
    } else {
      Log('[TabView] No tab change needed');
    }
  }, [lastFacetMap, route, tabs, activeTab, loading]);

  const nextTab = useCallback((): string => {
    const currentIndex = tabs.findIndex((tab) => tab.value === activeTab);
    const nextIndex = (currentIndex + 1) % tabs.length;
    return tabs[nextIndex]?.value || activeTab;
  }, [tabs, activeTab]);

  const prevTab = useCallback((): string => {
    const currentIndex = tabs.findIndex((tab) => tab.value === activeTab);
    const prevIndex = (currentIndex - 1 + tabs.length) % tabs.length;
    return tabs[prevIndex]?.value || activeTab;
  }, [tabs, activeTab]);

  useEvent<{ route: string; key: string }>(
    msgs.EventType.TAB_CYCLE,
    (_message: string, event?: { route: string; key: string }) => {
      if (!event) return;
      const vER = event.route?.replace(/^\/+/, '') || '';
      const vR = route.replace(/^\/+/, '');
      if (vER === vR) {
        const newTab = event.key.startsWith('alt+') ? prevTab() : nextTab();
        Log(`[TabView] Tab cycle event: "${newTab}" (route=${route})`);
        setActiveTab(newTab);
        setLastFacet(vR, newTab as types.DataFacet);
        if (onTabChange) {
          onTabChange(newTab);
        }
      }
    },
  );

  const handleTabChange = (newTab: string | null) => {
    if (newTab === null) return;
    Log(`[TabView] handleTabChange: "${newTab}" (route=${route})`);
    setActiveTab(newTab);
    const vR = route.replace(/^\/+/, '');
    setLastFacet(vR, newTab as types.DataFacet);
    if (onTabChange) {
      onTabChange(newTab);
    }
  };

  return (
    <div className="tab-view-container">
      <Tabs
        value={activeTab}
        onChange={(value) => {
          handleTabChange(value);
        }}
      >
        <Tabs.List>
          {tabs.map((tab, index) => (
            <Fragment key={`tab-${index}`}>
              {tab.dividerBefore && <TabDivider key={`divider-${index}`} />}
              <Tabs.Tab key={tab.value} value={tab.value}>
                {tab.label}
              </Tabs.Tab>
            </Fragment>
          ))}
        </Tabs.List>

        {tabs.map((tab) => (
          <Tabs.Panel key={tab.value} value={tab.value}>
            {tab.content}
          </Tabs.Panel>
        ))}
      </Tabs>
    </div>
  );
};
