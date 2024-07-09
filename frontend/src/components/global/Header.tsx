import React from "react";
import { Title } from "@mantine/core";

const Header = ({ title }: { title: string }) => {
  return <Title order={1}>{title}</Title>;
};

export default Header;
