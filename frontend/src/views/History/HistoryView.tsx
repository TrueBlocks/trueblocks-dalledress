import React, { useState, useEffect } from "react";
import { GetHistory, GetHistoryCnt } from "@gocode/app/App";
import classes from "@/App.module.css";
import { View, ViewStatus } from "@/components/view";
import { Stack, Title } from "@mantine/core";
import { types } from "@gocode/models";

export function HistoryView() {
  const [address, setAddress] = useState<string>("");
  const [cnt, setCnt] = useState<number>(0);
  const [txs, setTxs] = useState<types.Transaction[]>([]);

  useEffect(() => {
    setAddress("0xf503017d7baf7fbc0fff7492b751025c6a78179b");
    GetHistoryCnt(address).then((cnt: number) => {
      setCnt(cnt);
    });
    GetHistory(address, 0, 20).then((txs: types.Transaction[]) => {
      setTxs(txs);
    });
  }, []);

  return (
    <View>
      <Title order={3}>History Title</Title>
      <Stack className={classes.mainContent}>
        <div>{address}</div>
        <div>{cnt}</div>
        <div>{JSON.stringify(txs)}</div>
      </Stack>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}
