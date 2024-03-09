import ENV from '@/resources/env-constants';
import { ResponseError } from './common/type/response';

type RequestType = 'POST' | 'GET' | 'PUT' | 'DELETE';

const deleteService = async <S>(path: string, param?: any, header?: object) => {
  const res = await fetch(`${ENV.SERVER_PATH}/${path}`, {
    method: 'DELETE',
    body: param && JSON.stringify(param),
    headers: {
      'Content-Type': 'application/json',
      ...header
    },
    credentials: 'include',
    ...header
  });
  if (res.ok) {
    return res.json() as Promise<S>;
  } else {
    return createErrorContext(res);
  }
};

const postService = async <S>(path: string, param?: any, header?: object) => {
  const res = await fetch(`${ENV.SERVER_PATH}/${path}`, {
    method: 'POST',
    body: param && JSON.stringify(param),
    headers: {
      'Content-Type': 'application/json',
      ...header
    },
    credentials: 'include',
    ...header
  });
  if (res.ok) {
    return res.json() as Promise<S>;
  } else {
    return createErrorContext(res);
  }
};

const putService = async <S>(path: string, param?: any, header?: object) => {
  // https 경로로 들어오는 경우 그대로 사용하고, 아닌 경우 서버 주소를 붙여준다.
  // ex) 업로드의 경우 https://로 들어오는 경우가 있음(아마존에 S3 업로드)
  path = path.startsWith('https') ? path : `${ENV.SERVER_PATH}/${path}`;

  const res = await fetch(path, {
    method: 'PUT',
    body: param,
    headers: {
      'Content-Type': 'application/octet-stream',
      ...header
    },
    credentials: 'include'
  });
  if (res.ok) {
    return res.json() as Promise<S>;
  } else {
    return createErrorContext(res);
  }
};

const getService = async <S>(path: string, param?: any, header?: object) => {
  const res = await fetch(
    `${ENV.SERVER_PATH}/${path}?${new URLSearchParams(param).toString()}`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        ...header
      },
      credentials: 'include'
    }
  );
  if (res.ok) {
    return res.json() as Promise<S>;
  } else {
    return createErrorContext(res);
  }
};

const serviceModule = <S>(
  type: RequestType,
  path: string,
  param?: any,
  header?: object
) => {
  switch (type) {
    case 'POST':
      return postService<S>(path, param, header);
    case 'GET':
      return getService<S>(path, param, header);
    case 'PUT':
      return putService<S>(path, param, header);
    case 'DELETE':
      return deleteService<S>(path, param, header);
    default:
      throw new Error('Invalid Request Type');
  }
};

export default serviceModule;

const createErrorContext = (res: Response): Promise<ResponseError> => {
  return res.json().then((error: { error: { message: string } }) => ({
    status: res.status,
    error: error.error.message
  }));
};
