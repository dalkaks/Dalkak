'use client';
import {
  Carousel,
  CarouselContent,
  CarouselNext,
  CarouselPrevious
} from '@/components/ui/carousel';
import React from 'react';
import { Item } from '../../types/main/item';
import CustomCarouselItem from './CustomCarouselItem';

type Props = {
  items: Item[];
};
const CAROUSEL_OPTIONS = {
  loop: true,
  active: true,
  duration: 50
};
const ImageCarousel = ({ items }: Props) => {
  return (
    <Carousel className="w-9/12" opts={CAROUSEL_OPTIONS}>
      <CarouselContent>
        {items.map((item, index) => (
          <CustomCarouselItem key={index} item={item} />
        ))}
      </CarouselContent>
      <CarouselPrevious />
      <CarouselNext />
    </Carousel>
  );
};

export default ImageCarousel;
