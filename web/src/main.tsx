import React from 'react';
import ReactDOM from 'react-dom/client';
import { MetaMaskUIProvider } from '@metamask/sdk-react-ui';
import App from './App.tsx';
import CriticalErrorBoundary from './pages/CriticalErrorBoundary.tsx';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { HelmetProvider } from 'react-helmet-async';
import { ThemeProvider } from '@emotion/react';
import { SkeletonTheme } from 'react-loading-skeleton';
import { RouterProvider } from 'react-router-dom';
import theme from './styles/theme.tsx';
import GlobalStyle from './styles/GlobalStyle.tsx';
import router from './router/router.tsx';

const root = ReactDOM.createRoot(
	document.getElementById('root')! as HTMLElement
);

const globalQueryClient = new QueryClient({
	defaultOptions: {
		queries: {
			retry: 0,
		},
		mutations: {
			retry: 0,
		},
	},
});

root.render(
	<React.StrictMode>
		<MetaMaskUIProvider
			debug={false}
			sdkOptions={{
				dappMetadata: {
					name: 'Example React Dapp',
					url: window.location.host,
				},
				// Other options
			}}
		>
			<App />
		</MetaMaskUIProvider>
		<CriticalErrorBoundary>
			<QueryClientProvider client={globalQueryClient}>
				<HelmetProvider>
					<ThemeProvider theme={theme}>
						<SkeletonTheme
							enableAnimation={false}
							baseColor={theme.color.gray50}
						>
							<RouterProvider router={router} />
							<GlobalStyle />
						</SkeletonTheme>
					</ThemeProvider>
				</HelmetProvider>
				<ReactQueryDevtools initialIsOpen={false} />
			</QueryClientProvider>
		</CriticalErrorBoundary>
	</React.StrictMode>
);
