interface CustomErrors {
	[key: string]: string;
}

export const SERVER_ERR: CustomErrors = {
	NOT_FOUND: '잠시 후 다시 시도해 주세요',
	'Bad Request Exception': '잘못된 요청입니다',
};

export const CLIENT_ERR: CustomErrors = {
	UNKNOWN_ERR: '잠시 후 다시 시도해 주세요',
};
