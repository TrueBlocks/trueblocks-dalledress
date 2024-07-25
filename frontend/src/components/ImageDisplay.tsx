import React, { useState, useEffect } from "react";
import { Text, Image } from "@mantine/core";
import { GetTitle } from "@gocode/app/App";

export const ImageDisplay = ({ address, loading }: { address: string; loading: boolean }) => {
  var [title, setTitle] = useState<string>("");

  useEffect(() => {
    GetTitle(address).then((value: string) => {
      setTitle(value);
    });
  }, [address, loading]);

  const imgSrc = `http://localhost:8889/files/${address}&timestamp=${new Date().getTime()}`;
  if (loading || !address || address === "") {
    return (
      <div>
        <Text>{title}</Text>
        <Text>Loading...</Text>
      </div>
    );
  }

  return (
    <div>
      <Text c="white">{imgSrc}</Text>
      <Text c="white">{title}</Text>
      <Image fit="contain" src={imgSrc} alt={address} />
    </div>
  );
};
