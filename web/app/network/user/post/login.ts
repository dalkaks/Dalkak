import serviceModule from "../../serviceModule";

export interface RequestLogin {
  walletAddress: string;
  signature: string;
}

const LOGIN_ERRORS = {
  '400': 'Invalid wallet address or signature',

}

const loginErrorHandle = (res: Response) => {
  if (!res.ok) {
    throw new Error("Network response was not ok");
  }
  return res;
}


const loginService = async (req: RequestLogin) => {
  const res = await serviceModule('POST', 'user/auth', req)

  console.log(res)

  return res;
};

export default loginService;
