import ENV from '@/app/resources/env-constants';
import loginService from './login';
import mockWallet from '@/mock/walletData.json';

jest.spyOn(global, 'fetch');
describe('POST /user/login', () => {
  beforeEach(() => {
    // process.env.SERVER_PATH = "http://test.domain.com";
    ENV.SERVER_PATH = 'http://localhost:80';
  });
  it('POST 메소드 작동 횟수 및 데이터 점검', async () => {
    await loginService(mockWallet);

    expect(fetch).toHaveBeenCalledTimes(1);
    expect(fetch).toHaveBeenCalledWith(`${ENV.SERVER_PATH}/user/auth`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(mockWallet)
    });
  });

  it('로그인 정상작동 - accessToken 존재 유무 및 타입 체크', async () => {
    const res = await loginService(mockWallet);

    expect(res).toHaveProperty(['data', 'accessToken']);
    expect(typeof res.data.accessToken === 'string').toBe(true);
  });
});
