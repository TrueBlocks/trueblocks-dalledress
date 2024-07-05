import React, { useState, useEffect } from 'react';
import { Image } from '@mantine/core';
import { GetImage } from '@gocode/app/App';

// Define the type for the props
interface ImageDisplayProps {
  address: string;
}

const ImageDisplay: React.FC<ImageDisplayProps> = ({ address }) => {
  const [imageSrc, setImageSrc] = useState<string>('');

  useEffect(() => {
    if (address) {
      GetImage(address)
        .then((base64Image: string) => {
          if (base64Image !== 'No image file found. Press Generate.') {
            const dataUrl = `data:image/png;base64,${base64Image}`;
            setImageSrc(dataUrl);
          }
        })
        .catch((error: any) => {
          console.error('Error fetching image path:', error);
        });
    }
  }, [address]);

  if (!imageSrc || imageSrc === 'No image file found. Press Generate.') {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <div>{address}</div>
      <Image src={imageSrc} alt={address} />
      <div>{address}</div>
    </div>
  );
};

export default ImageDisplay;
