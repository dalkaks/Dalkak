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
    <Card className="flex sm:min-h-96 sm:min-w-96">
      <div className="h-full w-full overflow-hidden rounded-md">
        <AspectRatio ratio={1 / 1}>
          <Image
            src={image}
            className="rounded-md object-cover"
            alt="image"
            fill
          />
        </AspectRatio>
      </div>
    </Card>
  );
};

export default ImagePreview;
