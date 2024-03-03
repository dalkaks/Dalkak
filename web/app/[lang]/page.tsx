import styles from '../page.module.css';
import { Locale, getDictionary } from './dictionaries';
import RfButton from './components/RfButton';
import { CardWithForm } from './components/CardWithForm';
import { CardDemo } from './components/CardDemo';
import PsButton from './components/PsButton';
import { Toaster } from 'sonner';
import ImageCarousel from './components/main/ImageCarousel';
import { Item } from './types/main/item';

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

export default async function Home({ params: { lang } }: HomeProps) {
  const dict = await getDictionary(lang);
  return (
    <main className={`${styles.main} w-full overflow-hidden bg-slate-200`}>
      <ImageCarousel items={carouselItems} />
      <Toaster />
    </main>
  );
}
