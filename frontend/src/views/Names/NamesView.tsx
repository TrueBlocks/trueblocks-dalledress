import React, { useState, useEffect } from "react";
import { GetNames } from "@gocode/app/App";
import { types } from "@gocode/models";
import { useHotkeys } from "react-hotkeys-hook";
import classes from "@/App.module.css";
import { View, ViewStatus } from "@/components/view";
import { Stack, Title } from "@mantine/core";

export function NamesView() {
  const [names, setName] = useState<types.Name[]>();
  const [nNames, setNamesCount] = useState<number>(0);
  const [curName, setCurName] = useState<number>(-1);

  useHotkeys("left", (event) => {
    setCurName(curName - 1 < 0 ? 0 : curName - 1);
    event.preventDefault();
  });
  useHotkeys("up", (event) => {
    setCurName(curName - 20 < 0 ? 0 : curName - 20);
    event.preventDefault();
  });
  useHotkeys("right", (event) => {
    setCurName(curName + 1 > nNames ? nNames : curName + 1);
    event.preventDefault();
  });
  useHotkeys("down", (event) => {
    setCurName(curName + 20 > nNames ? nNames - 20 : curName + 20);
    event.preventDefault();
  });
  useHotkeys("home", (event) => {
    setCurName(0);
    event.preventDefault();
  });
  useHotkeys("end", (event) => {
    setCurName(nNames - 20);
    event.preventDefault();
  });

  useEffect(() => {
    GetNames(curName, 20).then((names) => {
      setName(names);
      setNamesCount(names?.length);
    });
  }, [curName]);

  useEffect(() => {
    setCurName(0);
  }, []);

  return (
    <View>
      <Title order={3}>Names Title</Title>
      <Stack className={classes.mainContent}>
        <pre>Number of records: {nNames}</pre>
        <pre>{JSON.stringify(names, null, 4)}</pre>
      </Stack>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}
