import styles from '../page.module.css';
import { Locale, getDictionary } from './dictionaries';
import Button from './components/Button';
import mockWallet from '@/mock/walletData.json';
import RfButton from './components/RfButton';
import { CardWithForm } from './components/CardWithForm';
import { CardDemo } from './components/CardDemo';
import PsButton from './components/PsButton';

interface HomeProps {
  params: {
    lang: Locale;
  };
}

export default async function Home({ params: { lang } }: HomeProps) {
  const dict = await getDictionary(lang);
  return (
    <main className={styles.main}>
      <Button data={mockWallet} title="login"></Button>
      <RfButton title="refresh"></RfButton>
      <PsButton title="presign"></PsButton>
      <CardWithForm />
      <CardDemo />
    </main>
  );
}
