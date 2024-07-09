import React from "react";
import { Center, Stack, Title } from "@mantine/core";
import Footer from "./Footer";
import Menu from "./Menu";

function Navbar() {
  return (
    <Stack h={"100%"} justify="space-between">
      <div>
        <Center>
          <Title order={1}>Browse</Title>
        </Center>
        <Menu />
      </div>
      <Footer />
    </Stack>
  );
}

export default Navbar;
