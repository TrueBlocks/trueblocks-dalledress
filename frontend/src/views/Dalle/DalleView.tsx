import React, { useState, useEffect } from "react";
import { Select, Paper, Grid, Button, Group } from "@mantine/core";
import classes from "@/App.module.css";
import { View } from "@/components/view/View";
import EditableSelect from "@/components/EditableSelect";
import ResultDialog from "@/components/ResultDialog";
import {
  GetJson,
  GetData,
  GetSeries,
  GetTerse,
  GetPrompt,
  GetEnhanced,
  GetFilename,
  GenerateImage,
  GetLastTab,
  GetLastAddress,
  GetLastSeries,
  SetLastTab,
  SetLastAddress,
  SetLastSeries,
  Save,
} from "@gocode/app/App";
import { ImageDisplay } from "@/components/ImageDisplay";
import Tabber from "./Tabber";

function DalleView() {
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
    GetLastTab().then((lastTab: string) => {
      setActiveTab(lastTab);
    });
    GetLastAddress().then((lastAddress: string) => {
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
      SetLastAddress(address);
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
      SetLastTab(value);
    }
  };

  const tabItems = [
    { value: "series", label: "Series", content: series },
    { value: "json", label: "JSON", content: json },
    { value: "data", label: "Data", content: data },
    { value: "terse", label: "Terse", content: terse },
    { value: "prompt", label: "Prompt", content: prompt },
    { value: "enhanced", label: "Enhanced", content: enhanced },
  ];

  return (
    <View title="Dalle View">
      <div className={classes.content}>
        <Grid>
          <Grid.Col span={8} className={classes.gridColumn}>
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
            <Paper shadow="xs" p="md" className={classes.imageDisplayContainer}>
              <ImageDisplay address={image} loading={imageLoading} />
            </Paper>
          </Grid.Col>
          <Grid.Col span={4} className={classes.gridColumn}>
            <Tabber items={tabItems} activeTab={activeTab} onTabChange={handleTabChange} />
          </Grid.Col>
        </Grid>
        <ResultDialog opened={dialogOpened} onClose={() => setDialogOpened(false)} success={success} />
      </div>
    </View>
  );
}

export default DalleView;
