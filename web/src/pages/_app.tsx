import '@/styles/globals.css';
import { Inter as FontSans } from 'next/font/google';
import type { AppProps } from 'next/app';

export const fontSans = FontSans({
	subsets: ['latin'],
	variable: '--font-sans',
});

export default function App({ Component, pageProps }: AppProps) {
	return (
		<div style={{ fontFamily: fontSans.variable }}>
			<Component {...pageProps} />;
		</div>
	);
}
