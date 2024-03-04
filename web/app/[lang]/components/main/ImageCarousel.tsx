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
import { useMediaQuery } from 'react-responsive';

type Props = {
  items: Item[];
};
const CAROUSEL_OPTIONS = {
  loop: true,
  active: true,
  duration: 50
};
const ImageCarousel = ({ items }: Props) => {
  const isDesktopAndTablet = useMediaQuery({
    query: '(min-width: 640px)'
  });
  return (
    <Carousel className="w-full p-4" opts={CAROUSEL_OPTIONS}>
      <CarouselContent>
        {items.map((item, index) => (
          <CustomCarouselItem key={index} item={item} />
        ))}
      </CarouselContent>
      {isDesktopAndTablet && <CarouselPrevious />}
      {isDesktopAndTablet && <CarouselNext />}
    </Carousel>
  );
};

export default ImageCarousel;
