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
