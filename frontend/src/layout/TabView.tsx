import { Fragment, useCallback, useEffect, useState } from 'react';

import { TabDivider } from '@components';
import { useActiveProject, useEvent } from '@hooks';
import { Tabs } from '@mantine/core';
import { msgs, types } from '@models';

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
  const { lastFacetMap, setLastFacet } = useActiveProject();

  const getInitialTab = (): string => {
    const savedTab = lastFacetMap[route];
    if (savedTab) {
      return savedTab;
    }
    return tabs[0]?.value || '';
  };

  const [activeTab, setActiveTab] = useState<string>(getInitialTab());

  useEffect(() => {
    const savedTab = lastFacetMap[route];
    if (savedTab && savedTab !== activeTab) {
      setActiveTab(savedTab);
    }
  }, [lastFacetMap, route, activeTab]);

  useEffect(() => {
    if (activeTab && !lastFacetMap[route]) {
      setLastFacet(route, activeTab as types.DataFacet);
      if (onTabChange) {
        onTabChange(activeTab);
      }
    }
  }, [activeTab, lastFacetMap, route, setLastFacet, onTabChange]);

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
      if (event?.route === route) {
        const newTab = event.key.startsWith('alt+') ? prevTab() : nextTab();
        setActiveTab(newTab);
        setLastFacet(route, newTab as types.DataFacet);
        if (onTabChange) {
          onTabChange(newTab);
        }
      }
    },
  );

  const handleTabChange = (newTab: string | null) => {
    if (newTab === null) return;
    setActiveTab(newTab);
    setLastFacet(route, newTab as types.DataFacet);
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
