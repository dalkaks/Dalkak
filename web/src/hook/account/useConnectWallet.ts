import { useSetRecoilState } from 'recoil'
import { accountState } from '../../state/accountState'
import { useSDK } from '@metamask/sdk-react'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'

export const useConnectWallet = () => {
  const { sdk } = useSDK()
  const setAccount = useSetRecoilState(accountState)
  const { t } = useTranslation()

  const connect = async () => {
    if (window.navigator.userAgent.includes('SamsungBrowser')) {
      toast.error(t('error-metamask-browser'), { duration: 750 })
      return
    }

    try {
      const signResult: any = await sdk?.connectAndSign({
        msg: t('connect-wallet'),
      })
      const address = window.ethereum?.selectedAddress
      if (!address) {
        toast.error(t('error-metamask-connect'), { duration: 750 })
        return
      }

      setAccount(address)

      const apiUrl = process.env.NEXT_PUBLIC_API_URL
      const response = await fetch(`${apiUrl}/user/auth`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          WalletAddress: address,
          Signature: signResult,
        }),
      })
      if (response.status !== 200) {
        toast.error(t('error-metamask-connect'), { duration: 750 })
        return
      }
    } catch (err) {
      toast.error(t('error-metamask-connect'), { duration: 750 })
    }
  }

  return connect
}
