import React, { useState, useEffect } from "react";
import { ActionIcon } from "@mantine/core";
import {
  MantineProvider,
  Paper,
  Button,
  Group,
  Box,
  TextInput,
  Text
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useClipboard } from "@mantine/hooks";
import { GetPrompt, GetData, GetJson } from "../wailsjs/go/main/App";

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

  return (
    <div id="App" style={{ width: "98vw", margin: "auto" }}>
      <MantineProvider>
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
        <Paper shadow="xs" p="md" style={{ maxWidth: "100vw", marginLeft: 0 }}>
          <Text>{email ? email : "Working..."}</Text>
          <CopyText prompt={prompt ? prompt : "Working..."} />
          <CopyText prompt={data ? data : "Working..."} />
          <CopyText prompt={json ? json : "Working..."} />
        </Paper>
      </MantineProvider>
    </div>
  );
}

export const CopyText = ({ prompt }: { prompt?: string }) => {
  const { copy } = useClipboard();
  const promptText = prompt ? prompt : "Book Now";
  return (
    <Group align="flex-start" justify="flex-start">
      <Text>{promptText}</Text>
      <ActionIcon onClick={() => copy(promptText)}>
        <Text>Copy</Text>
      </ActionIcon>
    </Group>
  );
};
