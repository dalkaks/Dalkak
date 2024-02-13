import errorHandler, { errorGuard } from '../../common/util/errorHandler';
import { ResponseSuccess } from '../../common/type/response';
import serviceModule from '../../serviceModule';

export interface RequestLogin {
  walletAddress: string;
  signature: string;
}

interface ResponseLogin {
  accessToken: string;
}

const ERROR_CASE: { [key in number]: any } = {
  '400': {
    DEFAULT: 'BAD_REQUEST',
    MISSING_SIGNATURE: 'MISSING_SIGNATURE',
    MISSING_WALLET_ADDRESS: 'MISSING_WALLET_ADDRESS'
  },
  '401': {
    DEFAULT: 'UNAUTHORIZED',
    INVALID_SIGNATURE: 'INVALID_SIGNATURE',
    INVALID_WALLET_ADDRESS: 'INVALID_WALLET_ADDRESS'
  },
  '500': {
    DEFAULT: 'INTERNAL_SERVER_ERROR'
  }
};

const loginService = async (req: RequestLogin) => {
  const res = await serviceModule<ResponseSuccess<ResponseLogin>>(
    'POST',
    'user/auth',
    req
  );
  if (errorGuard(res)) throw errorHandler(res, ERROR_CASE);
  return res;
};

export default loginService;
