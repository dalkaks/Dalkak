import { useQueryErrorResetBoundary } from '@tanstack/react-query';
import { isAxiosError } from 'axios';
import { ErrorBoundary } from 'react-error-boundary';
import useToast from '../hooks/useToast';
import { useNavigate } from 'react-router-dom';
import { useStore } from 'zustand';
import usePaymentStackStore from './store/historyStack';

interface Props {
	children: React.PropsWithChildren<React.ReactNode>;
}

const RetryErrorBoundary = ({ children }: Props) => {
	const { reset } = useQueryErrorResetBoundary();
	const { errorToast } = useToast();
	const navigate = useNavigate();
	const { stack, pop } = useStore(usePaymentStackStore);

	return (
		<ErrorBoundary
			onReset={reset}
			onError={(error) => {
				if (
					isAxiosError(error) &&
					(error?.response?.status === 500 ||
						error?.response?.data?.message === 'Internal Server Error')
				) {
					throw error;
				}
			}}
			fallbackRender={({ error, resetErrorBoundary }) => {
				errorToast(error.message);
				resetErrorBoundary();
				if (stack.length > 1) {
					navigate(stack[stack.length - 1], {
						replace: true,
					});
					pop();
				} else {
					navigate(-1);
				}
				return null;
			}}
		>
			{children}
		</ErrorBoundary>
	);
};

export default RetryErrorBoundary;
