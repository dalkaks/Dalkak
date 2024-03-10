import errorHandler, { errorGuard } from '../../common/util/errorHandler';
import { ResponseSuccess } from '../../common/type/response';
import serviceModule from '../../serviceModule';
import { UploadUrl } from '../type/media';

export interface RequestUpload {
  uploadUrl: UploadUrl;
  file: File;
}

interface ResponseUpload {}

const ERROR_CASE: { [key in number]: any } = {
  '400': {
    DEFAULT: 'BAD_REQUEST'
  },
  '401': {
    DEFAULT: 'UNAUTHORIZED',
    REFRESH_TOKEN_NOT_FOUND: 'REFRESH_TOKEN_NOT_FOUND'
  },
  '500': {
    DEFAULT: 'INTERNAL_SERVER_ERROR',
    REQUEST: {
      INVALID_REQUEST: 'PRESIGN_INVALID_REQUEST'
    }
  }
};

const upload = async ({ uploadUrl, file }: RequestUpload) => {
  const res = await serviceModule<ResponseSuccess<ResponseUpload>>(
    'PUT',
    uploadUrl,
    file
  );
  console.log('UPLOAD : ', res);
  if (errorGuard(res)) throw errorHandler(res, ERROR_CASE);
  return res;
};

export default upload;
