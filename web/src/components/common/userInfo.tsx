import { useSDK } from '@metamask/sdk-react'
import { useState } from 'react'
import { ReloadIcon } from '@radix-ui/react-icons'
import { Button } from '../ui/button'

// import { useTranslation } from 'react-i18next'

export default function UserInfo() {
  const [account, setAccount] = useState<string>()
  const { sdk, connected, connecting, provider, chainId } = useSDK()

  const connect = async () => {
    try {
      const accounts: any = await sdk?.connect()
      setAccount(accounts?.[0])
    } catch (err) {
      console.warn(`failed to connect..`, err)
    }
  }

  return (
    <div>
      {!connected && !connecting && <Button onClick={connect}>Connect</Button>}
      {connecting && (
        <Button disabled>
          <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
          Connect
        </Button>
      )}
      {connected && (
        <div>
          <>
            {chainId}
            <p />
            {/* {account} */}
          </>
        </div>
      )}
    </div>
  )
}
