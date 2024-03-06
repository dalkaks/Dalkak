'use server';

import errorHandler, { errorGuard } from '../../common/util/errorHandler';
import { ResponseSuccess } from '../../common/type/response';
import serviceModule from '../../serviceModule';
import { NftFileExt } from '@/type/nft/fileExtension';
import { cookies } from 'next/headers';

export interface RequestPresign {
  mediaType: 'image';
  ext: NftFileExt;
  prefix: 'board';
}

interface ResponsePresign {
  id: string;
  accessUrl: `https://${string}.${RequestPresign['ext']}`; // https + 파일경로 + 확장자
  uploadUrl: string;
}

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

const presign = async (param: RequestPresign) => {
  const token = cookies().get('access_token')?.value;
  if (!token) throw new Error('로그인이 필요합니다');
  const res = await serviceModule<ResponseSuccess<ResponsePresign>>(
    'POST',
    'media/presigned',
    param,
    {
      Authorization: `Bearer ${token}`
    }
  );
  if (errorGuard(res)) throw errorHandler(res, ERROR_CASE);
  return res;
};

export default presign;
