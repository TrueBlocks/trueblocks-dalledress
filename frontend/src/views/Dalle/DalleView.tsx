import React, { useState, useEffect } from 'react';
import { Paper, Grid, Tabs, TextInput, Button, Group, Text } from '@mantine/core';
import classes from '../View.module.css';
import View from '@/components/view/View';
import { GetJson, GetData, GetTitle, GetTerse, GetPrompt, GetEnhanced, GetImage } from '@gocode/app/App';
import { GetLastTab, SetLastTab } from '@gocode/app/App';
import ImageDisplay from '@/components/image/ImageDisplay';

function DalleView() {
  const [address, setAddress] = useState<string>('0xf503017d7baf7fbc0fff7492b751025c6a78179b');
  const [json, setJson] = useState<string>('');
  const [data, setData] = useState<string>('');
  const [prompt, setPrompt] = useState<string>('');
  const [enhanced, setEnhanced] = useState<string>('');
  const [title, setTitle] = useState<string>('');
  const [terse, setTerse] = useState<string>('');
  const [imagePath, setImagePath] = useState<string>('');
  const [isGenerateDisabled, setIsGenerateDisabled] = useState<boolean>(false);
  const [activeTab, setActiveTab] = useState<string>('json');

  useEffect(() => {
    GetLastTab().then((lastTab: string) => {
      setActiveTab(lastTab);
    });
  }, []);

  const handleGenerate = () => {
    setIsGenerateDisabled(true);
    crypto.subtle.digest('SHA-256', new TextEncoder().encode(address)).then((hashBuffer) => {
      GetImage(address).then((value: string) => {
        setImagePath(value);
      });
      GetEnhanced(address).then((value: string) => {
        setEnhanced(value);
        setIsGenerateDisabled(false);
      });
    });
  };

  useEffect(() => {
    GetJson(address).then((value: string) => {
      setJson(value);
    });
    GetData(address).then((value: string) => {
      setData(value);
    });
    GetTitle(address).then((value: string) => {
      setTitle(value);
    });
    GetTerse(address).then((value: string) => {
      setTerse(value);
    });
    GetPrompt(address).then((value: string) => {
      setPrompt(value);
    });
  }, [address]);

  useEffect(() => {
    setIsGenerateDisabled(false);
  }, [address]);

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
            <TextInput
              value={address}
              onChange={(event) => setAddress(event.currentTarget.value)}
              placeholder="Enter address"
              label="Address"
              style={{ width: '600px' }}
            />
            <Button onClick={handleGenerate} style={{ marginTop: '22px' }} disabled={isGenerateDisabled}>
              Generate
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
            <Text mt="md">Image: {imagePath}</Text>
            <ImageDisplay address={address} />
          </Paper>
        </Grid.Col>
        <Grid.Col span={4} className={classes.gridColumn}>
          <Tabs value={activeTab} onChange={handleTabChange}>
            <Tabs.List>
              <Tabs.Tab value="json">JSON</Tabs.Tab>
              <Tabs.Tab value="data">Data</Tabs.Tab>
              <Tabs.Tab value="title">Title</Tabs.Tab>
              <Tabs.Tab value="terse">Terse</Tabs.Tab>
              <Tabs.Tab value="prompt">Prompt</Tabs.Tab>
              <Tabs.Tab value="enhanced">Enhanced</Tabs.Tab>
            </Tabs.List>

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
            <Tabs.Panel value="title" className={classes.tabPanel}>
              <Text mt="md">{title}</Text>
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
    </View>
  );
}

export default DalleView;
