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
				if (isAxiosError(error)) {
					throw error;
				}
			}}
			fallbackRender={({ error }) => {
				if (SERVER_ERR[error.message]) {
					errorToast(SERVER_ERR[error.message]);
				} else {
					errorToast(CLIENT_ERR['UNKNOWN_ERR']);
				}
				return null;
			}}
		>
			{children}
		</ErrorBoundary>
	);
};

export default CriticalErrorBoundary;
