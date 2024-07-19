import React, { useState, useEffect } from "react";
import { GetHistory, GetHistoryCnt } from "@gocode/app/App";
import classes from "@/App.module.css";
import { View, ViewStatus } from "@/components/view";
import { Group, Stack, Title, Text } from "@mantine/core";
import { types } from "@gocode/models";
import { useHotkeys } from "react-hotkeys-hook";

export function HistoryView() {
  const [address, setAddress] = useState<string>("0x054993ab0f2b1acc0fdc65405ee203b4271bebe6");
  const [items, setItems] = useState<types.Transaction[]>([]);
  const [nItems, setNItems] = useState<number>(0);
  const [curItem, setCurItem] = useState<number>(0);
  const perPage = 20;

  useHotkeys("left", (event) => {
    setCurItem((cur) => Math.max(cur - 1, 0));
    event.preventDefault();
  });
  useHotkeys("up", (event) => {
    setCurItem((cur) => Math.max(cur - perPage, 0));
    event.preventDefault();
  });
  useHotkeys("home", (event) => {
    setCurItem(0);
    event.preventDefault();
  });

  useHotkeys("right", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem((cur) => Math.min(max, cur + 1));
    event.preventDefault();
  });
  useHotkeys("down", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem((cur) => Math.min(max, cur + perPage));
    event.preventDefault();
  });
  useHotkeys("end", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem(max);
    event.preventDefault();
  });

  useEffect(() => {
    const fetchHistory = async () => {
      const fetchedItems = await GetHistory(address, curItem, perPage);
      setItems(fetchedItems);
      const cnt = await GetHistoryCnt(address);
      setNItems(cnt);
    };
    fetchHistory();
  }, [address, curItem, perPage]);

  useEffect(() => {
    setCurItem(0);
  }, [address]);

  return (
    <View>
      <Title order={3}>History Title</Title>
      <Stack className={classes.mainContent}>
        <Text>
          {address} at record {curItem} of {nItems} records
        </Text>
        <Group>
          {items.map((item, idx) => (
            <Text key={idx}>
              {item.blockNumber} {item.from} {item.to} {item.value}
            </Text>
          ))}
        </Group>
      </Stack>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}
