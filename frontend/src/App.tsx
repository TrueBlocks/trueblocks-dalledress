import React, { useState, useEffect } from "react";
import { ActionIcon } from "@mantine/core";
import {
  Paper,
  Button,
  Group,
  Box,
  TextInput,
  Text,
  CopyButton,
  Tooltip
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { IconCopy, IconCheck } from "@tabler/icons-react";
import {
  GetPrompt,
  GetData,
  GetJson,
  GetImage,
  GetImageData
  // GetImprovedPrompt,
  // GetModeration
} from "../wailsjs/go/main/App";
import "./App.css";
import "@mantine/core/styles.css";

export default function App() {
  const [email, setEmail] = useState<string>("trueblocks.eth");
  const [prompt, setPrompt] = useState<string>("");
  const [data, setData] = useState<string>("");
  const [json, setJson] = useState<string>("");
  // const [moderation, setModeration] = useState<string>("");
  // const [improvedPrompt, setImprovedPrompt] = useState<string>("");

  const form = useForm({
    initialValues: {
      email: "trueblocks.eth",
      termsOfService: false
    }
  });

  useEffect(() => {
    GetPrompt(form.values["email"]).then((value: string) => {
      setPrompt(value);
    });
    GetData(form.values["email"]).then((value: string) => {
      setData(value);
    });
    GetJson(form.values["email"]).then((value: string) => {
      setJson(value);
    });
    // GetModeration(form.values["email"]).then((value: string) => {
    //   setModeration(value);
    // });
    // GetImprovedPrompt(form.values["email"]).then((value: string) => {
    //   setImprovedPrompt(value);
    // });
  }, [email]);

  const openImage = (email: string) => {
    GetImage(email);
  };

  return (
    <div id="App" style={{ width: "98vw", margin: "auto" }}>
      <Box mx="auto">
        <form
          onSubmit={form.onSubmit((values) => setEmail(values["email"]))}
          onBlur={form.onSubmit((values) => setEmail(values["email"]))}
        >
          <TextInput
            label="Email"
            placeholder="your@email.com"
            {...form.getInputProps("email")}
          />
          <Group mt="md">
            <Button type="submit">Submit</Button>
          </Group>
        </form>
      </Box>
      <Button onClick={() => openImage(email)}>Generate</Button>{" "}
      {/* <div>{improvedPrompt}</div> */}
      {/* <div>{moderation}</div> */}
      <Paper
        shadow="xs"
        p="md"
        style={{
          maxWidth: "100vw",
          marginLeft: 0,
          border: "1px solid green"
        }}
      >
        <Text>{email ? email : "Working..."}</Text>
        <CopyText prompt={prompt ? prompt : "Working..."} />
        <CopyText prompt={data ? data : "Working..."} />
        <CopyText prompt={json ? json : "Working..."} />
      </Paper>
    </div>
  );
}

export const CopyText = ({ prompt }: { prompt?: string }) => {
  const promptText = prompt ? prompt : "Book Now";

  return (
    <div className="shit-container">
      <ImageDisplay />
      <div className="shit">{promptText}</div>
      <CopyButton value={promptText} timeout={2000}>
        {({ copied, copy }) => (
          <Tooltip
            label={copied ? "Copied" : "Copy"}
            withArrow
            position="right"
          >
            <ActionIcon
              color={copied ? "teal" : "gray"}
              variant="subtle"
              onClick={copy}
            >
              {copied ? (
                <IconCheck style={{ width: "1rem" }} />
              ) : (
                <IconCopy style={{ width: "1rem" }} />
              )}
            </ActionIcon>
          </Tooltip>
        )}
      </CopyButton>
    </div>
  );
};

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
