import { accountState } from "@/state/accountState"
import { useRecoilState } from "recoil"

export const useParseAccessToken = () => {
  const [account, setAccount] = useRecoilState(accountState)
  const accessToken = localStorage.getItem('accessToken')

  const parse = () => {
    if (accessToken) {
      const parsedAddress = parseAccessToken(accessToken)
      if (account && account != parsedAddress) {
        localStorage.removeItem('accessToken')
        setAccount(null)
      } else {
        setAccount(parsedAddress)
      }
    }
  }

  return parse
}

function parseAccessToken(token: string) {
  try {
    const base64Payload = token.split('.')[1]
    const decodedPayload = atob(base64Payload)
    const payload = JSON.parse(decodedPayload)

    // if exp 지나면 엑세스 토큰 재발급

    return payload.sub
  } catch (e) {
    console.error('Failed to parse access token', e)
    return null
  }
}
