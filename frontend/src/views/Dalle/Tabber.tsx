import React from 'react';
import { Tabs, Text } from '@mantine/core';
import classes from '../View.module.css';

interface TabItem {
  value: string;
  label: string;
  content: string;
}

interface TabViewProps {
  items: TabItem[];
  activeTab: string;
  onTabChange: (value: string | null) => void;
}

const Tabber: React.FC<TabViewProps> = ({ items, activeTab, onTabChange }) => {
  return (
    <Tabs value={activeTab} onChange={onTabChange}>
      <Tabs.List>
        {items.map((item) => (
          <Tabs.Tab key={item.value} value={item.value}>
            {item.label}
          </Tabs.Tab>
        ))}
      </Tabs.List>

      {items.map((item) => (
        <Tabs.Panel key={item.value} value={item.value} className={classes.tabPanel}>
          <Text mt="md">
            <pre>{item.content}</pre>
          </Text>
        </Tabs.Panel>
      ))}
    </Tabs>
  );
};

export default Tabber;
