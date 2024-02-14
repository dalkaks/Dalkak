import errorHandler, { errorGuard } from '../../common/util/errorHandler';
import { ResponseSuccess } from '../../common/type/response';
import serviceModule from '../../serviceModule';

interface ResponseRefresh {
  accessToken: string;
}

const ERROR_CASE: { [key in number]: any } = {
  '400': {
    DEFAULT: 'BAD_REQUEST'
  },
  '401': {
    DEFAULT: 'UNAUTHORIZED',
    REFRESH_TOKEN_NOT_FOUND: 'REFRESH_TOKEN_NOT_FOUND',
  },
  '500': {
    DEFAULT: 'INTERNAL_SERVER_ERROR'
  }
};

const refresh = async () => {
  const res = await serviceModule<ResponseSuccess<ResponseRefresh>>(
    'POST',
    'user/refresh'
  );
  if (errorGuard(res)) throw errorHandler(res, ERROR_CASE);
  return res;
};

export default refresh;
