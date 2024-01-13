import React from 'react';
import ReactDOM from 'react-dom/client';
import { MetaMaskUIProvider } from '@metamask/sdk-react-ui';
import App from './App.tsx';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
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
	</React.StrictMode>
);
