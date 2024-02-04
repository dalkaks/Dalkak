import { reissueAccessToken } from '@/api/auth'
import { accountState } from '@/state/accountState'
import { useRecoilState } from 'recoil'

export const useParseAccessToken = () => {
  const [account, setAccount] = useRecoilState(accountState);

  const parse = async () => {
    if (typeof window === 'undefined') return;

    const accessToken = localStorage.getItem('accessToken');
    if (!accessToken) return;

    try {
      const parsedAddress = await parseAccessToken(accessToken);
      if (account && account !== parsedAddress) {
        localStorage.removeItem('accessToken');
        setAccount(null);
      } else {
        setAccount(parsedAddress);
      }
    } catch (error) {
      console.error('Error parsing access token:', error);
    }
  };

  return parse;
};

async function parseAccessToken(token: string, retryCount = 0) {
  if (!token) throw new Error('No token provided');

  const payload = decodeTokenPayload(token);
  if (!payload || !payload.sub || !payload.exp) {
    throw new Error('Invalid token structure');
  }

  if (payload.exp < Date.now() / 1000) {
    if (retryCount >= 1) throw new Error('Token reissue limit reached');

    await reissueAccessToken();
    const newToken = localStorage.getItem('accessToken');
    if (!newToken) throw new Error('Failed to reissue token');
    
    return parseAccessToken(newToken, retryCount + 1);
  }

  return payload.sub;
}

function decodeTokenPayload(token: string) {
  try {
    const base64Payload = token.split('.')[1];
    const decodedPayload = atob(base64Payload);
    return JSON.parse(decodedPayload);
  } catch (e) {
    throw new Error('Failed to decode token payload');
  }
}