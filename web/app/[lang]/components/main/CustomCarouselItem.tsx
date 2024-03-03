import React from 'react';
import { Item } from '../../types/main/item';
import { CarouselItem } from '@/components/ui/carousel';
import Image from 'next/image';
import { Card } from '@/components/ui/card';

type Props = { item: Item };

const CustomCarouselItem = ({ item }: Props) => {
  return (
    <CarouselItem>
      <Card className="flex max-h-[400px] max-w-full justify-center overflow-hidden align-middle">
        <Image
          className="w-full"
          width={1200}
          height={300}
          src={item.src}
          alt={item.alt}
        />
      </Card>
    </CarouselItem>
  );
};

export default CustomCarouselItem;
