import styles from '../page.module.css';
import { Locale, getDictionary } from './dictionaries';
import { Toaster } from 'sonner';
import ImageCarousel from './components/main/ImageCarousel';
import { Item } from './types/main/item';
import Board from './components/board/Board';

export const metadata = {
  icons: {
    icon: '/icons/favicon.ico'
  }
};

interface HomeProps {
  params: {
    lang: Locale;
  };
}

const carouselItems: Item[] = [
  {
    src: '/images/carouselImageTest.webp',
    alt: 'Image 1'
  },
  {
    src: '/images/carouselImageTest.webp',
    alt: 'Image 2'
  },
  {
    src: '/images/carouselImageTest.webp',
    alt: 'Image 3'
  }
];

// TODO: Replace with actual board images
const BOARD_IMAGES = [
  '/images/mock/board/board-mock-1.jpg',
  '/images/mock/board/board-mock-2.jpg',
  '/images/mock/board/board-mock-3.jpg',
  '/images/mock/board/board-mock-4.jpg',
  '/images/mock/board/board-mock-5.jpg',
  '/images/mock/board/board-mock-6.jpg'
];

export default async function Home({ params: { lang } }: HomeProps) {
  const dict = await getDictionary(lang);
  return (
    <main className={`${styles.main} w-full overflow-hidden bg-slate-200`}>
      <ImageCarousel items={carouselItems} />
      <Board imageSrcs={BOARD_IMAGES} />
      <Toaster />
    </main>
  );
}
