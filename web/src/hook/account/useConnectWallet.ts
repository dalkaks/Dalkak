import { useSetRecoilState } from 'recoil'
import { accountState } from '../../state/accountState'
import { useSDK } from '@metamask/sdk-react'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import { authenticateUserWithSignature } from '@/api/auth'

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

      await authenticateUserWithSignature(address, signResult)
      setAccount(address)
    } catch (err) {
      toast.error(t('error-metamask-connect'), { duration: 750 })
    }
  }

  return connect
}
