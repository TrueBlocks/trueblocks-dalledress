import React, { useState, useEffect } from "react";
import { GetFilelist } from "@gocode/app/App";
import classes from "@/App.module.css";
import { View, ViewStatus, StringTable, DataItem } from "@components";
import { Stack, Text, Title } from "@mantine/core";

export function SeriesView() {
  const [fileList, setFileList] = useState<string[]>([]);

  useEffect(() => {
    GetFilelist("./output/series").then((fileList: string[]) => {
      fileList = fileList.map((file) => file.replace("output/series/", ""));
      fileList = fileList.filter((file) => file.includes(".json") && file !== ".json");
      fileList = fileList.map((file) => file.replace(".json", ""));
      setFileList(fileList);
    });
  }, []);

  const dataItems = fileList.map((file, index) => ({ id: index, value: file }));
  return (
    <View>
      <Title order={3}>Series Title</Title>
      <Stack className={classes.mainContent}>
        <StringTable data={dataItems} />
      </Stack>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}
