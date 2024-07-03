import React, { useState } from "react";
import { Container, TextInput, Button, Group, Text } from "@mantine/core";
import classes from "../View.module.css";
import View from "@/components/view/View";

function DalleView() {
  const [address, setAddress] = useState("");
  const [firstHash, setFirstHash] = useState("");
  const [appendedString, setAppendedString] = useState("");

  const handleGenerate = () => {
    crypto.subtle
      .digest("SHA-256", new TextEncoder().encode(address))
      .then((hashBuffer) => {
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray
          .map((b) => b.toString(16).padStart(2, "0"))
          .join("");
        setFirstHash(hashHex);
        const appended = `${address.slice(0, 16)} | ${address.slice(16, 32)} | ${hashHex.slice(0, 16)} | ${hashHex.slice(16, 32)} | ${hashHex.slice(32, 48)}`;
        setAppendedString(appended);
      });
  };

  return (
    <View title="Dalle View">
      <TextInput
        value={address}
        onChange={(event) => setAddress(event.currentTarget.value)}
        placeholder="Enter address"
        label="Address"
      />
      <Group mt="md" style={{ justifyContent: "center" }}>
        <Button onClick={handleGenerate}>Generate</Button>
      </Group>
      <Text mt="md">First Hash: {firstHash}</Text>
      <Text mt="md">Appended String: {appendedString}</Text>
    </View>
  );
}

export default DalleView;

/*
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
  GetEnhanced,
  GetTerse,
  GetData,
  GetJson,
  GetImage
  // GetImageData
  // GetImprovedPrompt,
  // GetModeration
} from "../wailsjs/go/app/App";
import "./App.css";
import "@mantine/core/styles.css";

export default function App() {
  const [email, setEmail] = useState<string>(
    "0xf503017d7baf7fbc0fff7492b751025c6a78179b"
  );
  const [prompt, setPrompt] = useState<string>("");
  const [enhanced, setEnhanced] = useState<string>("");
  const [terse, setTerse] = useState<string>("");
  const [data, setData] = useState<string>("");
  const [json, setJson] = useState<string>("");
  // const [moderation, setModeration] = useState<string>("");
  // const [improvedPrompt, setImprovedPrompt] = useState<string>("");

  const form = useForm({
    initialValues: {
      email: "0xf503017d7baf7fbc0fff7492b751025c6a78179b",
      termsOfService: false
    }
  });

  useEffect(() => {
    GetPrompt(form.values["email"]).then((value: string) => {
      setPrompt(value);
    });
    GetEnhanced(form.values["email"]).then((value: string) => {
      setEnhanced(value);
    });
    GetTerse(form.values["email"]).then((value: string) => {
      setTerse(value);
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

export const CopyText = ({ prompt }: { prompt?: string }) => {
  const promptText = prompt ? prompt : "Book Now";

  return (
    <div className="shit-container">
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

*/
