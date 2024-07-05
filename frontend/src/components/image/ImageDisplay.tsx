import React from 'react';
import { Image } from '@mantine/core';

export const ImageDisplay = ({ address }: { address: string }) => {
  const imgSrc = `http://localhost:8082/files/${address}`;

  return (
    <div>
      <div>{address}</div>
      <Image src={imgSrc} alt={address} />
      <div>{address}</div>
    </div>
  );
};
