import React from "react";
import { Tabs, Text } from "@mantine/core";
import classes from "../View.module.css";

interface TabItem {
  value: string;
  label: string;
  content: string;
}

interface TabberItem {
  item: TabItem;
  pre: boolean;
}
interface TabViewProps {
  items: TabberItem[];
  activeTab: string;
  onTabChange: (value: string | null) => void;
}

const Tabber: React.FC<TabViewProps> = ({ items, activeTab, onTabChange }) => {
  return (
    <Tabs value={activeTab} onChange={onTabChange}>
      <Tabs.List>
        {items.map((tItem) => {
          var item = tItem.item;
          return (
            <Tabs.Tab key={item.value} value={item.value}>
              {item.label}
            </Tabs.Tab>
          );
        })}
      </Tabs.List>

      {items.map((tItem) => {
        var item = tItem.item;
        return (
          <Tabs.Panel key={item.value} value={item.value} className={classes.tabPanel}>
            {tItem.pre ? (
              <Text mt="sm">
                <small>
                  <pre> {item.content}</pre>
                </small>
              </Text>
            ) : (
              <Text mt="sm">{item.content}</Text>
            )}
          </Tabs.Panel>
        );
      })}
    </Tabs>
  );
};

export default Tabber;
