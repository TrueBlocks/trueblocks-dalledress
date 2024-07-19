import React, { useState, useEffect } from "react";
import { Stack, Title, Paper, Grid, Button, Group } from "@mantine/core";
import classes from "@/App.module.css";
import { View, ViewStatus, EditableSelect, ResultDialog, ImageDisplay } from "@components";
import {
  GetJson,
  GetData,
  GetSeries,
  GetTerse,
  GetPrompt,
  GetEnhanced,
  GetFilename,
  GenerateImage,
  GetLast,
  SetLast,
  Save,
} from "@gocode/app/App";
import Tabber from "./Tabber";

export function DalleView() {
  const [address, setAddress] = useState<string>("");
  const [json, setJson] = useState<string>("");
  const [data, setData] = useState<string>("");
  const [prompt, setPrompt] = useState<string>("");
  const [enhanced, setEnhanced] = useState<string>("");
  const [series, setSeries] = useState<string>("");
  const [terse, setTerse] = useState<string>("");
  const [image, setImage] = useState<string>("");
  const [imageLoading, setImageLoading] = useState<boolean>(false);
  const [activeTab, setActiveTab] = useState<string>("json");
  const [dialogOpened, setDialogOpened] = useState(false);
  const [success, setSuccess] = useState(false);
  const gridWidth = 8;

  const handleOpenDialog = (result: boolean) => {
    setSuccess(result);
    setDialogOpened(true);
  };

  const handleGenerate = () => {
    setEnhanced("Loading...");
    setImageLoading(true);
    GenerateImage(address).then((value: string) => {
      setImageLoading(false);
      setEnhanced(value);
    });
  };

  const handleSave = () => {
    Save(address).then((value: boolean) => {
      handleOpenDialog(value);
    });
  };

  // On first load of the view, set the most recently viewed tab and address
  useEffect(() => {
    GetLast("tab").then((lastTab: string) => {
      setActiveTab(lastTab);
    });
    GetLast("address").then((lastAddress: string) => {
      if (lastAddress && lastAddress.length > 0) {
        setAddress(lastAddress);
      } else {
        setAddress("0xf503017d7baf7fbc0fff7492b751025c6a78179b");
      }
    });
  }, []);

  // When address changes, update all the data
  useEffect(() => {
    if (address && !imageLoading) {
      SetLast("address", address);
      GetJson(address).then((value: string) => {
        setJson(value);
      });
      GetData(address).then((value: string) => {
        setData(value);
      });
      GetSeries(address).then((value: string) => {
        setSeries(value);
      });
      GetTerse(address).then((value: string) => {
        setTerse(value);
      });
      GetPrompt(address).then((value: string) => {
        setPrompt(value);
      });
      GetEnhanced(address).then((value: string) => {
        setEnhanced(value);
      });
      GetFilename(address).then((value: string) => {
        setImage(value);
      });
    }
  }, [address, imageLoading]);

  const handleTabChange = (value: string | null) => {
    if (value) {
      setActiveTab(value);
      SetLast("tab", value);
    }
  };

  const tabItems = [
    { pre: true, item: { value: "series", label: "Series", content: series } },
    { pre: true, item: { value: "json", label: "JSON", content: json } },
    { pre: true, item: { value: "data", label: "Data", content: data } },
    { pre: false, item: { value: "terse", label: "Terse", content: terse } },
    { pre: false, item: { value: "prompt", label: "Prompt", content: prompt } },
    { pre: false, item: { value: "enhanced", label: "Enhanced", content: enhanced } },
  ];

  return (
    <View>
      <Title order={3}>Dalle View</Title>
      <Stack h="100%" className={classes.mainContent}>
        <Grid>
          <Grid.Col span={gridWidth}>
            <Stack h="100%" className={classes.mainContent}>
              <Group mt="md" style={{ justifyContent: "flex-start" }}>
                <EditableSelect
                  value={address}
                  onChange={(value) => setAddress(value)}
                  label="Select or enter an address or ENS name"
                  placeholder="Enter or select an address"
                />
                <Button onClick={handleGenerate} style={{ marginTop: "22px" }}>
                  Generate
                </Button>
                <Button onClick={handleSave} style={{ marginTop: "22px" }}>
                  Save
                </Button>
              </Group>
              <Paper shadow="xs" p="md">
                <ImageDisplay address={image} loading={imageLoading} />
              </Paper>
            </Stack>
          </Grid.Col>
          <Grid.Col span={12 - gridWidth}>
            <Tabber items={tabItems} activeTab={activeTab} onTabChange={handleTabChange} />
          </Grid.Col>
        </Grid>
        <ResultDialog opened={dialogOpened} onClose={() => setDialogOpened(false)} success={success} />
      </Stack>
      <ViewStatus>Status / Progress</ViewStatus>
    </View>
  );
}
