import ENV from '@/resources/env-constants';
import logout from './logout';

jest.spyOn(global, 'fetch');
describe('POST /user/logout', () => {
  beforeEach(() => {
    // process.env.SERVER_PATH = "http://test.domain.com";
    ENV.SERVER_PATH = 'http://localhost:80';
  });

  it('로그아웃 정상작동', async () => {
    const res = await logout();

    expect(res).toHaveProperty(['data']);
  });
});
