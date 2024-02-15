import errorHandler, { errorGuard } from '../../common/util/errorHandler';
import { ResponseSuccess } from '../../common/type/response';
import serviceModule from '../../serviceModule';
import { NftFileExt } from '@/type/nft/fileExtension';

interface RequestPresign {
  mediaType: "image",
  ext: NftFileExt,
  prefix: "board"
}

interface ResponsePresign {
  id: string;
  url: `https://${string}.${RequestPresign['ext']}` // https + 파일경로 + 확장자
  presignedUrl: string;
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

const presign = async (param: RequestPresign) => {
  const res = await serviceModule<ResponseSuccess<ResponsePresign>>(
    'POST',
    'user/media/presigned',
    param
  );
  if (errorGuard(res)) throw errorHandler(res, ERROR_CASE);
  return res;
};

export default presign;
