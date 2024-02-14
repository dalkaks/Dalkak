import ENV from '../resources/env-constants';
import { ResponseError } from './common/type/response';

type RequestType = 'POST' | 'GET';

const postService = async <S, E extends ResponseError>(
  path: string,
  param?: any
) => {
  const res = await fetch(`${ENV.SERVER_PATH}/${path}`, {
    method: 'POST',
    body: param && JSON.stringify(param),
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include'
  });
  console.log(res);
  if (res.ok) {
    return res.json() as Promise<S>;
  } else {
    return createErrorContext(res) as Promise<E>;
  }
};

const getService = async <S, E extends ResponseError>(
  path: string,
  param?: any
) => {
  const res = await fetch(
    `${ENV.SERVER_PATH}/${path}?${new URLSearchParams(param).toString()}`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include'
    }
  );
  console.log(res);
  if (res.ok) {
    return res.json() as Promise<S>;
  } else {
    return createErrorContext(res) as Promise<E>;
  }
};

const serviceModule = <S>(type: RequestType, path: string, param?: any) => {
  switch (type) {
    case 'POST':
      return postService<S, ResponseError>(path, param);
    case 'GET':
      return getService<S, ResponseError>(path, param);
  }
};

export default serviceModule;

const createErrorContext = (res: Response) => {
  return res.json().then((error: { error: { message: string } }) => ({
    status: res.status,
    ...error
  }));
};
