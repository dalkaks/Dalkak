import '@/styles/globals.css';
import { Inter as FontSans } from 'next/font/google';
import type { AppProps } from 'next/app';
import Head from 'next/head';

export const fontSans = FontSans({
	subsets: ['latin'],
	variable: '--font-sans',
});

export default function App({ Component, pageProps }: AppProps) {
	return (
		<>
			<Head>
				<meta
					name="viewport"
					content="width=device-width, initial-scale=1.0, viewport-fit=cover"
				/>
			</Head>
			<div style={{ fontFamily: fontSans.variable }}>
				<Component {...pageProps} />;
			</div>
		</>
	);
}
