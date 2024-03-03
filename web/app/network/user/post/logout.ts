import errorHandler, { errorGuard } from '../../common/util/errorHandler';
import { ResponseSuccess } from '../../common/type/response';
import serviceModule from '../../serviceModule';

interface ResponseLogout {}

const logout = async () => {
  const res = await serviceModule<ResponseSuccess<ResponseLogout>>(
    'POST',
    'user/logout'
  );
  if (errorGuard(res)) throw errorHandler(res);
  return res;
};

export default logout;
