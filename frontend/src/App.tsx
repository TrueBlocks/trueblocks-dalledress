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
import { GetPrompt, GetData, GetJson, GetImage } from "../wailsjs/go/main/App";
import "./App.css";
import "@mantine/core/styles.css";

export default function App() {
  const [email, setEmail] = useState<string>("");
  const [prompt, setPrompt] = useState<string>("");
  const [data, setData] = useState<string>("");
  const [json, setJson] = useState<string>("");

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
  }, [form.values]);

  const openImage = (email: string) => {
    GetImage(email);
  };

  return (
    <div id="App" style={{ width: "98vw", margin: "auto" }}>
      <Box mx="auto">
        <form onSubmit={form.onSubmit((values) => setEmail(values["email"]))}>
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
      <Button onClick={() => openImage(email)}>Generate</Button>
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
  //     <Group align="flex-start" justify="flex-start">
  //   <Text>{promptText}</Text>
  //   <ActionIcon onClick={() => copy(promptText)}>
  //     <Text>Copy</Text>
  //   </ActionIcon>
  // </Group>
};
