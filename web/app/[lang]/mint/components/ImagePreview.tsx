import { AspectRatio } from '@/components/ui/aspect-ratio';
import { Card } from '@/components/ui/card';
import Image from 'next/image';
import React from 'react';

type Props = {
  file: File;
};

const ImagePreview = ({ file }: Props) => {
  const image = URL.createObjectURL(file);
  return (
    <Card className="w-[450px] flex">
      <AspectRatio ratio={16 / 9}>
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
