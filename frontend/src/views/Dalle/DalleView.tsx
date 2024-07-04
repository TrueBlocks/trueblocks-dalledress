import React, { useState, useEffect } from 'react';
import { Paper, Grid, Tabs, TextInput, Button, Group, Text, ScrollArea } from '@mantine/core';
import classes from '../View.module.css';
import View from '@/components/view/View';
import { GetJson, GetData, GetTitle, GetTerse, GetPrompt, GetEnhanced, GetImage } from '@gocode/app/App';

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

  const handleGenerate = () => {
    setIsGenerateDisabled(true);
    crypto.subtle.digest('SHA-256', new TextEncoder().encode(address)).then((hashBuffer) => {
      setImagePath('x');
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
              border: '1px solid green',
            }}
          >
            <Text mt="md">Image: {imagePath}</Text>
            <div style={{ height: '500px' }}></div> {/* Blank space below */}
          </Paper>
        </Grid.Col>
        <Grid.Col span={4} className={classes.gridColumn}>
          <Tabs defaultValue="json">
            <ScrollArea style={{ whiteSpace: 'nowrap' }}>
              <Tabs.List>
                <Tabs.Tab value="json">JSON</Tabs.Tab>
                <Tabs.Tab value="data">Data</Tabs.Tab>
                <Tabs.Tab value="title">Title</Tabs.Tab>
                <Tabs.Tab value="terse">Terse</Tabs.Tab>
                <Tabs.Tab value="prompt">Prompt</Tabs.Tab>
                <Tabs.Tab value="enhanced">Enhanced</Tabs.Tab>
              </Tabs.List>
            </ScrollArea>

            <Tabs.Panel value="json">
              <Text mt="md">
                <pre>{json}</pre>
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="data">
              <Text mt="md">
                <pre>{data}</pre>
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="title">
              <Text mt="md">{title}</Text>
            </Tabs.Panel>
            <Tabs.Panel value="terse">
              <Text mt="md" style={{ textAlign: 'justify' }}>
                {terse}
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="prompt">
              <Text mt="md" style={{ textAlign: 'justify' }}>
                {prompt}
              </Text>
            </Tabs.Panel>
            <Tabs.Panel value="enhanced">
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

/*
  const form = useForm({
    initialValues: {
      address: "0xf503017d7baf7fbc0fff7492b751025c6a78179b",
    }
  });

  const openImage = (address: string) => {
    GetImage(address);
  };

  return (
    <div id="App" style={{ width: "98vw", margin: "auto" }}>
      <Box mx="auto">
        <form
          onSubmit={form.onSubmit((values) => setEmail(values["address"]))}
          onBlur={form.onSubmit((values) => setEmail(values["address"]))}
        >
          <TextInput
            label="Email"
            placeholder="your@address.com"
            {...form.getInputProps("address")}
          />
          <Group mt="md">
            <Button type="submit">Submit</Button>
          </Group>
        </form>
      </Box>
      <Button onClick={() => openImage(address)}>Generate</Button>{" "}
        <Text>{address ? address : "Working..."}</Text>
        <div>Prompt:</div>
        <CopyText prompt={prompt ? prompt : "Working..."} />
        <div>Enhanced:</div>
        <CopyText prompt={enhanced ? enhanced : "Working..."} />
        <div>Terse:</div>
        <CopyText prompt={terse ? terse : "Working..."} />
        <div>Data:</div>
        <CopyText prompt={data ? data : "Working..."} />
        <div>Json:</div>
        <CopyText prompt={json ? json : "Working..."} />
      </Paper>
    </div>
  );
}

function ImageDisplay() {
  const [imageSrc, setImageSrc] = useState("");

  const fetchImage = async () => {
    const imageData = await GetImageData(); // Adjust based on your actual IPC call
    setImageSrc(imageData);
  };

  return (
    <div>
      <button onClick={fetchImage}>Load Image</button>
      {imageSrc && <img src={imageSrc} alt="Dynamically loaded" />}
    </div>
  );
}

*/
