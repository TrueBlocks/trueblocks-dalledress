import React from "react";
import { Center, Stack, Title } from "@mantine/core";
import GlobalFooter from "./GlobalFooter";
import GlobalMenu from "./GlobalMenu";

function GlobalNavbar() {
  return (
    <Stack h={"100%"} justify="space-between">
      <div>
        <Center>
          <Title order={1}>Browse</Title>
        </Center>
        <GlobalMenu />
      </div>
      <GlobalFooter />
    </Stack>
  );
}

export default GlobalNavbar;
