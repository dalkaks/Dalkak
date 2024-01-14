import Toast from '../components/common/Toast';

const useToast = () => {
	const successToast = (title: string) => {
		Toast.fire({ icon: 'success', title });
	};

	const errorToast = (title: string) => {
		Toast.fire({ icon: 'error', title });
	};

	const warningToast = (title: string) => {
		Toast.fire({ icon: 'warning', title });
	};

	const infoToast = (title: string) => {
		Toast.fire({ icon: 'info', title });
	};

	const questionToast = (title: string) => {
		Toast.fire({ icon: 'question', title });
	};

	return { successToast, errorToast, warningToast, infoToast, questionToast };
};

export default useToast;
