import React from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { useQueryErrorResetBoundary } from '@tanstack/react-query';
import useToast from '../hooks/useToast';
import { CLIENT_ERR, SERVER_ERR } from '../constant/errMsg';
import { isAxiosError } from 'axios';

const CriticalErrorBoundary = ({
	children,
}: React.PropsWithChildren<unknown>) => {
	const { reset } = useQueryErrorResetBoundary();
	const { errorToast } = useToast();

	return (
		<ErrorBoundary
			onReset={reset}
			onError={(error) => {
				if (
					// 이 ErrorBoundary에서 처리하면 안되는 오류의 경우 상위 ErrorBoundary로 위임
					!isAxiosError(error) ||
					!error.response
				) {
					return;
				}
				throw error;
			}}
			fallbackRender={({ error }) => {
				if (!SERVER_ERR[error.message]) {
					errorToast(SERVER_ERR[error.message]);
				} else {
					errorToast(CLIENT_ERR['UNKNOWN_ERR']);
				}

				return (
					<div>
						<h1> 알 수 없는 에러가 발생했습니다. </h1>
						<footer> {error.message} </footer>
					</div>
				);
			}}
		>
			{children}
		</ErrorBoundary>
	);
};

export default CriticalErrorBoundary;
