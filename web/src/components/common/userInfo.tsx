import { useSDK } from '@metamask/sdk-react'
import { useState } from 'react'
import { DotFilledIcon, ReloadIcon } from '@radix-ui/react-icons'
import { useTranslation } from 'react-i18next'
import { Button } from '../ui/button'
import { toast } from 'sonner'

export default function UserInfo() {
  const { t } = useTranslation()
  const [account, setAccount] = useState<string>()
  const { sdk, connected, connecting, provider, chainId } = useSDK()

  const connect = async () => {
    if (window.navigator.userAgent.includes('SamsungBrowser')) {
      toast.error(t('error-metamask-browser'), { duration: 750 })
      return
    }

    try {
      const accounts: any = await sdk?.connect()
      setAccount(accounts?.[0])
    } catch (err) {
      toast.error(t('error-metamask-connect'), { duration: 750 })
    }
  }

  return (
    <div>
      {!connected && !connecting && (
        <Button className="pl-1" onClick={connect}>
          <DotFilledIcon className="mx-1 h-4 w-4 text-red-500" />
          Connect
        </Button>
      )}
      {connecting && (
        <Button className="pl-1" disabled>
          <ReloadIcon className="mx-1 h-4 w-4 animate-spin" />
          Connect
        </Button>
      )}
      {connected && (
        <Button className="pl-1">
          <DotFilledIcon className="mx-1 h-4 w-4 text-green-500" />
          {account?.slice(0, 4)}
          ...{account?.slice(-4)}
        </Button>
      )}
    </div>
  )
}
