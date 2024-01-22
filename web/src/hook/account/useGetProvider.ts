import { useSetRecoilState } from 'recoil'
import { accountState } from '../../state/accountState'
import detectEthereumProvider from '@metamask/detect-provider'

export const useGetProvider = () => {
  const setAccount = useSetRecoilState(accountState)

  const getProvider = async () => {
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
