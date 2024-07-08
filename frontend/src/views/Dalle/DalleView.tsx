import React, { useState, useEffect } from 'react';
import { Select, Paper, Grid, Tabs, TextInput, Button, Group, Text } from '@mantine/core';
import classes from '../View.module.css';
import View from '@/components/view/View';
import EditableSelect from '@/components/EditableSelect';
import ResultDialog from '@/components/ResultDialog';
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
} from '@gocode/app/App';
import { ImageDisplay } from '@/components/ImageDisplay';

function DalleView() {
  const [address, setAddress] = useState<string>('');
  const [json, setJson] = useState<string>('');
  const [data, setData] = useState<string>('');
  const [prompt, setPrompt] = useState<string>('');
  const [enhanced, setEnhanced] = useState<string>('');
  const [series, setSeries] = useState<string>('');
  const [terse, setTerse] = useState<string>('');
  const [image, setImage] = useState<string>('');
  const [imageLoading, setImageLoading] = useState<boolean>(false);
  const [activeTab, setActiveTab] = useState<string>('json');
  const [dialogOpened, setDialogOpened] = useState(false);
  const [success, setSuccess] = useState(false);

  const handleOpenDialog = (result: boolean) => {
    setSuccess(result);
    setDialogOpened(true);
  };

  const handleGenerate = () => {
    setEnhanced('Loading...');
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
        setAddress('0xf503017d7baf7fbc0fff7492b751025c6a78179b');
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

  return (
    <View title="Dalle View">
      <Grid>
        <Grid.Col span={8} className={classes.gridColumn}>
          <Group mt="md" style={{ justifyContent: 'flex-start' }}>
            <EditableSelect
              value={address}
              onChange={(value) => setAddress(value)}
              label="Select or enter an address or ENS name"
              placeholder="Enter or select an address"
            />
            <Button onClick={handleGenerate} style={{ marginTop: '22px' }}>
              Generate
            </Button>
            <Button onClick={handleSave} style={{ marginTop: '22px' }}>
              Save
            </Button>
          </Group>
          <Paper
            shadow="xs"
            p="md"
            style={{
              maxWidth: '100vw',
              marginLeft: 0,
            }}
          >
            <ImageDisplay address={image} loading={imageLoading} />
          </Paper>
        </Grid.Col>
        <Grid.Col span={4} className={classes.gridColumn}>
          <Tabs value={activeTab} onChange={handleTabChange}>
            <Tabs.List>
              <Tabs.Tab value="series">Series</Tabs.Tab>
              <Tabs.Tab value="json">JSON</Tabs.Tab>
              <Tabs.Tab value="data">Data</Tabs.Tab>
              <Tabs.Tab value="terse">Terse</Tabs.Tab>
              <Tabs.Tab value="prompt">Prompt</Tabs.Tab>
              <Tabs.Tab value="enhanced">Enhanced</Tabs.Tab>
            </Tabs.List>

            <Tabs.Panel value="series" className={classes.tabPanel}>
              <Text mt="md">
                <pre>{series}</pre>
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="json" className={classes.tabPanel}>
              <Text mt="md">
                <pre>{json}</pre>
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="data" className={classes.tabPanel}>
              <Text mt="md">
                <pre>{data}</pre>
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="terse" className={classes.tabPanel}>
              <Text mt="md" style={{ textAlign: 'justify' }}>
                {terse}
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="prompt" className={classes.tabPanel}>
              <Text mt="md" style={{ textAlign: 'justify' }}>
                {prompt}
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="enhanced" className={classes.tabPanel}>
              <Text mt="md" style={{ textAlign: 'justify' }}>
                {enhanced}
              </Text>
            </Tabs.Panel>
          </Tabs>
        </Grid.Col>
      </Grid>
      <ResultDialog opened={dialogOpened} onClose={() => setDialogOpened(false)} success={success} />{' '}
    </View>
  );
}

export default DalleView;
