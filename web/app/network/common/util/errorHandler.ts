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

const errorHandler = (res: ResponseError, errorCase?: typeof ERROR_CASE) => {
  errorCase ?? ERROR_CASE;
  if (!ERROR_CASE[res.status]) throw new Error('Unknown Error');

  const error = ERROR_CASE[res.status][res.error.message || 'DEFAULT'];

  return new Error(error);
};

export default errorHandler;
