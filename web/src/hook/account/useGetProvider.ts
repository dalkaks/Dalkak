import { useSetRecoilState } from 'recoil'
import { accountState } from '../../state/accountState'
import detectEthereumProvider from '@metamask/detect-provider'

export const useGetProvider = () => {
  const setAccount = useSetRecoilState(accountState)

  const isMobile = () => {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      navigator.userAgent,
    )
  }
  const isUnsupportedBrowser = () => {
    return window.navigator.userAgent.includes('SamsungBrowser')
  }

  const getProvider = async () => {
    if (isMobile() || isUnsupportedBrowser()) {
      return
    }

    const provider = (await detectEthereumProvider({
      mustBeMetaMask: true,
      silent: true,
    })) as any

    if (provider) {
      const accounts = await provider.request({ method: 'eth_requestAccounts' })
      if (accounts.length > 0) {
        setAccount(accounts[0])
      }
    }
  }

  return getProvider
}
