import { ScrollRestoration, createBrowserRouter } from 'react-router-dom';
import RetryErrorBoundary from '../pages/RetryErrorBoundary';
import Root from './Root';

const router = createBrowserRouter([
	{
		path: '/',
		element: (
			<RetryErrorBoundary>
				<Root />
				<ScrollRestoration />
			</RetryErrorBoundary>
		),
		children: [
			// {
			//   element: <Layout headerType="textWhite" isPaddingTop={false} />,
			//   children: [
			//     {
			//       index: true,
			//       element: <MainPage />,
			//     },
			//     {
			//       // 존재하지 않는 경로 진입시 메인페이지로 리다이렉트
			//       path: '*',
			//       element: <Navigate replace to="/" />,
			//     },
			//   ],
			// },
		],
	},
]);

export default router;
