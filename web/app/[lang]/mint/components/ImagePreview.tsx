import { AspectRatio } from '@/components/ui/aspect-ratio';
import { Card } from '@/components/ui/card';
import Image from 'next/image';
import React from 'react';

type Props = {
  file: File;
};

const ImagePreview = ({ file }: Props) => {
  const defaultImage = '/images/no-image.webp';
  const image = file.size ? URL.createObjectURL(file) : defaultImage;
  return (
    <Card className="flex h-[450px] w-[450px]">
      <AspectRatio ratio={1 / 1}>
        <Image
          src={image}
          className="rounded-md object-cover"
          alt="image"
          fill
        />
      </AspectRatio>
    </Card>
  );
};

export default ImagePreview;
