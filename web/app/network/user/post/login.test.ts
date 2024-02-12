import ENV from "@/app/resources/env-constants";
import loginService from "./login";

const mockReq = {
  walletAddress: "myWalletAddress",
  signature: "mySignature",
};

const mockRes = {
  data: {
    accessToken: "successGetAccessToken",
  },
};

jest.spyOn(global, "fetch");
describe("POST /user/login", () => {
  beforeEach(() => {
    // process.env.SERVER_PATH = "http://test.domain.com";
    ENV.SERVER_PATH = "http://test.domain.com";
  });
  it("should send a POST request to the correct endpoint", async () => {
    await loginService(mockReq);

    expect(fetch).toHaveBeenCalledTimes(1);
    expect(fetch).toHaveBeenCalledWith(`http://test.domain.com/user/login`, {
      method: "POST",
      body: JSON.stringify(mockReq),
    });
  });

  it("should return the response", async () => {
    (fetch as jest.Mock).mockResolvedValueOnce(mockRes);

    const res = await loginService(mockReq);

    expect(res).toEqual(mockRes);
  });
});
