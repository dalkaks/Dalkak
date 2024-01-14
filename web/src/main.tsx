import React from 'react';
import ReactDOM from 'react-dom/client';
import { MetaMaskUIProvider } from '@metamask/sdk-react-ui';
import App from './App.tsx';
import './index.css';
import CriticalErrorBoundary from './pages/CriticalErrorBoundary.tsx';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

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
			<p>hello</p>
			<QueryClientProvider client={globalQueryClient}>
				{/* <HelmetProvider>
					<ThemeProvider theme={theme}>
						<SkeletonTheme
							enableAnimation={false}
							baseColor={theme.color.gray50}
						>
							<RouterProvider router={router} />
							<GlobalStyle />
							<Toaster toastOptions={toasterOptions} />
						</SkeletonTheme>
					</ThemeProvider>
				</HelmetProvider>
				<ReactQueryDevtools initialIsOpen={false} /> */}
			</QueryClientProvider>
		</CriticalErrorBoundary>
	</React.StrictMode>
);
