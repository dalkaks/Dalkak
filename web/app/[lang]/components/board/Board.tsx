import { AspectRatio } from '@/components/ui/aspect-ratio';
import { Card, CardContent } from '@/components/ui/card';
import Image from 'next/image';
import React from 'react';

type Props = {
  imageSrcs: string[];
};
const Board = ({ imageSrcs }: Props) => {
  return (
    <div className="container grid grid-cols-1 gap-4 sm:grid-cols-4">
      {imageSrcs.map((image, index) => (
        <Card>
          <CardContent className="overflow-hidden rounded-lg p-0 shadow-md hover:shadow-lg">
            <AspectRatio ratio={1 / 1} key={index}>
              <Image
                key={index}
                src={image}
                alt={`board-image-${index + 1}`}
                layout="fill"
              />
            </AspectRatio>
          </CardContent>
        </Card>
      ))}
    </div>
  );
};

export default Board;
