'use server'

import errorHandler, { errorGuard } from '../../common/util/errorHandler';
import { ResponseSuccess } from '../../common/type/response';
import serviceModule from '../../serviceModule';
import { NftFileExt } from '@/type/nft/fileExtension';
import { cookies } from 'next/headers';

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
  const token = cookies().get('access_token');
  if (!token) throw new Error('로그인이 필요합니다');
  const res = await serviceModule<ResponseSuccess<ResponsePresign>>(
    'POST',
    'user/media/presigned',
    param,
    {
      Authorization: `Bearer ${token.value}`
    }
  );
  if (errorGuard(res)) throw errorHandler(res, ERROR_CASE);
  console.log(res)
  return res;
};

export default presign;
