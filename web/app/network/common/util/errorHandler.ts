import { ResponseError } from '../type/response';

export const ERROR_CASE: { [key in number]: any } = {
  '400': {
    DEFAULT: 'BAD_REQUEST'
  },
  '401': {
    DEFAULT: 'UNAUTHORIZED'
  },
  '500': {
    DEFAULT: 'INTERNAL_SERVER_ERROR'
  }
};

export const errorGuard = (arg: any): arg is ResponseError => {
  return arg.error;
};

const splitErrorMessage = (message: string) => {
  const specificError = message.split(':');
  return {
    origin: specificError[0],
    detail: specificError[1]
  };
};

const errorHandler = (
  res: ResponseError,
  customErrorCase?: typeof ERROR_CASE
) => {
  const errorCase = customErrorCase ?? ERROR_CASE;
  if (!errorCase[res.status]) throw new Error('Unknown Error');

  const errorQuery = Object.keys(errorCase[res.status]);
  const { origin, detail } = splitErrorMessage(res.error);

  if (errorQuery.includes(origin)) {
    const error =
      errorCase[res.status][origin][detail] ?? errorCase[res.status]['DEFAULT'];
    return new Error(error);
  }
  const error = errorCase[res.status]['DEFAULT'];
  return new Error(error);
};

export default errorHandler;
