import React, { useState, useEffect } from "react";
import { GetSeries } from "@gocode/app/App";
import classes from "@/App.module.css";
import { View, ViewStatus, ViewTitle } from "@components";
import { Stack, Text, Title } from "@mantine/core";

export function SeriesView() {
  const [fileList, setFileList] = useState<string[]>([]);

  useEffect(() => {
    GetSeries("./output/series").then((fileList: string[]) => {
      fileList = fileList.map((file) => file.replace("output/series/", ""));
      fileList = fileList.filter((file) => file.includes(".json") && file !== ".json");
      fileList = fileList.map((file) => file.replace(".json", ""));
      setFileList(fileList);
    });
  }, []);

  const dataItems = fileList.map((file, index) => ({ id: index, value: file }));
  return (
    <View>
      <Stack className={classes.mainContent}>
        <ViewTitle />
        <div>NEEDS</div>
      </Stack>
      <ViewStatus />
    </View>
  );
}
