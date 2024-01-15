import { useState } from 'react';
import reactLogo from './assets/react.svg';
import viteLogo from '/vite.svg';
import {
	MetaMaskButton,
	useAccount,
	useSDK,
	useSignMessage,
} from '@metamask/sdk-react-ui';

// function App() {
// 	const [count, setCount] = useState(0);

// 	function generateCuid() {
// 		// 간단한 CUID 생성 함수
// 		return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(
// 			/[xy]/g,
// 			function (c) {
// 				const r = (Math.random() * 16) | 0,
// 					v = c === 'x' ? r : (r & 0x3) | 0x8;
// 				return v.toString(16);
// 			}
// 		);
// 	}

// 	function getCuid() {
// 		// 로컬 스토리지에서 CUID를 가져오거나 새로 생성
// 		let cuid = localStorage.getItem('cuid');
// 		if (!cuid) {
// 			cuid = generateCuid();
// 			localStorage.setItem('cuid', cuid);
// 		}
// 		return cuid;
// 	}

// 	function sendRequestWithCuid() {
// 		const cuid = getCuid();
// 		const url = import.meta.env.VITE_SERVER_URL as string;

// 		fetch(url, {
// 			method: 'GET',
// 			headers: {
// 				'X-Client-Id': cuid,
// 			},
// 		})
// 			.then((data) => console.log(data))
// 			.catch((error) => console.error('Error:', error));
// 	}

// 	sendRequestWithCuid();

// 	return (
// 		<>
// 			<div>
// 				<a href="https://vitejs.dev" target="_blank">
// 					<img src={viteLogo} className="logo" alt="Vite logo" />
// 				</a>
// 				<a href="https://react.dev" target="_blank">
// 					<img src={reactLogo} className="logo react" alt="React logo" />
// 				</a>
// 			</div>
// 			<h1>Vite + React</h1>
// 			<div className="card">
// 				<button onClick={() => setCount((count) => count + 1)}>
// 					count is {count}
// 				</button>
// 				<p>
// 					Edit <code>src/App.tsx</code> and save to test HMR
// 				</p>
// 			</div>
// 			<p className="read-the-docs">
// 				Click on the Vite and React logos to learn more
// 			</p>
// 		</>
// 	);
// }

function AppReady() {
	const {
		data: signData,
		isError: isSignError,
		isLoading: isSignLoading,
		isSuccess: isSignSuccess,
		signMessage,
	} = useSignMessage({
		message: 'gm wagmi frens',
	});

	const { isConnected } = useAccount();

	return (
		<div className="App">
			<header className="App-header">
				<MetaMaskButton theme={'light'} color="white"></MetaMaskButton>
				{isConnected && (
					<>
						<div style={{ marginTop: 20 }}>
							<button disabled={isSignLoading} onClick={() => signMessage()}>
								Sign message
							</button>
							{isSignSuccess && <div>Signature: {signData}</div>}
							{isSignError && <div>Error signing message</div>}
						</div>
					</>
				)}
			</header>
		</div>
	);
}

function App() {
	const { ready } = useSDK();

	if (!ready) {
		return <div>Loading...</div>;
	}

	return <AppReady />;
}

export default App;
