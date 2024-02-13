import ENV from "../resources/env-constants";

type RequestType = 'POST' | 'GET';

const postService = async (path: string, param?: any) => {
  const res = await fetch(`${ENV.SERVER_PATH}/${path}`, {
    method: 'POST',
    body: param && JSON.stringify(param),
    headers: {
      "Content-Type": "application/json",
    },
  });

  return res;
}

const getService = async (path: string, param?: any) => {
  const res = await fetch(`${ENV.SERVER_PATH}/${path}?${new URLSearchParams(param).toString()}`, {
    method: 'GET',
    headers: {
      "Content-Type": "application/json",
    },
  });

  return res;
}

const serviceModule = (type: RequestType, path: string, param?: any) => {
  switch (type) {
    case 'POST':
      return postService(path, param);
    case 'GET':
      return getService(path, param);
  }
}

export default serviceModule;