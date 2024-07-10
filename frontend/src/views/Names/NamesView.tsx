import React, { useState, useEffect } from "react";
import { GetNames, MaxNames, GetFilelist } from "@gocode/app/App";
import { useHotkeys } from "react-hotkeys-hook";
import classes from "@/App.module.css";
import { View, ViewStatus } from "@/components/view";
import { Stack, Text, Title } from "@mantine/core";

function NamesView() {
  const [names, setName] = useState<string[]>();
  const [curName, setCurName] = useState<number>(0);
  const [maxNames, setMaxNames] = useState<number>(0);
  const [fileList, setFileList] = useState<string[]>();

  useHotkeys("left", (event) => {
    event.preventDefault();
    setCurName(curName - 1 < 0 ? 0 : curName - 1);
  });
  useHotkeys("up", (event) => {
    event.preventDefault();
    setCurName(curName - 20 < 0 ? 0 : curName - 20);
  });
  useHotkeys("right", (event) => {
    event.preventDefault();
    setCurName(curName + 1 > maxNames ? maxNames : curName + 1);
  });
  useHotkeys("down", (event) => {
    event.preventDefault();
    setCurName(curName + 20 > maxNames ? maxNames - 20 : curName + 20);
  });
  useHotkeys("home", (event) => {
    event.preventDefault();
    setCurName(0);
  });
  useHotkeys("end", (event) => {
    event.preventDefault();
    setCurName(maxNames - 20);
  });

  useEffect(() => {
    console.log("useEffect", curName);
    GetNames(curName, 20).then((names: string[]) => setName(names));
    MaxNames().then((maxNames: number) => setMaxNames(maxNames));
    GetFilelist().then((fileList: string[]) => setFileList(fileList));
  }, [curName]);

  return (
    <View>
      <Title order={3}>Home View Title</Title>
      <Stack className={classes.mainContent}>
        {fileList?.map((item, index) => <div key={index}>{item}</div>)}
        <pre>Number of records: {maxNames}</pre>
        <pre>{JSON.stringify(names, null, 4)}</pre>
        {/* <Text>Home View Content</Text> */}
      </Stack>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}

export default NamesView;
