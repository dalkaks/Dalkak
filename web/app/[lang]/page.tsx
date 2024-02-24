import styles from '../page.module.css';
import { Locale, getDictionary } from './dictionaries';
import RfButton from './components/RfButton';
import { CardWithForm } from './components/CardWithForm';
import { CardDemo } from './components/CardDemo';
import PsButton from './components/PsButton';
import { Toaster } from 'sonner';

interface HomeProps {
  params: {
    lang: Locale;
  };
}

export default async function Home({ params: { lang } }: HomeProps) {
  const dict = await getDictionary(lang);
  return (
    <main className={`${styles.main} bg-slate-200 `}>
      <RfButton title="refresh"></RfButton>
      <PsButton title="presign"></PsButton>
      <CardWithForm />
      <CardDemo />
      <Toaster />
    </main>
  );
}
