import ENV from "@/app/resources/env-constants";

interface RequestLogin {
  walletAddress: string;
  signature: string;
}

const loginService = async (req: RequestLogin) => {
  const res = await fetch(`${ENV.SERVER_PATH}/user/login`, {
    method: "POST",
    body: JSON.stringify(req),
  });

  return res;
};

export default loginService;
