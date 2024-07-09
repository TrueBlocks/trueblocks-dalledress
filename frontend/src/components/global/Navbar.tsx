import React from "react";
import { Center, Stack, Title } from "@mantine/core";
import Menu from "./Menu";
import ServerStatus from "./ServerStatus";

function Navbar() {
  return (
    <Stack h={"100%"} justify="space-between">
      <div>
        <Center>
          <Title order={1}>Browse</Title>
        </Center>
        <Menu />
      </div>
      <ServerStatus />
    </Stack>
  );
}

export default Navbar;
