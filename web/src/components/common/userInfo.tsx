import { useSDK } from '@metamask/sdk-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

function Welcome() {
  const { t } = useTranslation()

  return <div>{t('dashboard')}</div>
}

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
      {!connected && !connecting && (
        <button style={{ padding: 10, margin: 10 }} onClick={connect}>
          Connect
        </button>
      )}
      {connected && (
        <div>
          <>
            {chainId && `Connected chain: ${chainId}`}
            <p></p>
            {account && `Connected account: ${account}`}
          </>
        </div>
      )}
      <Welcome />
    </div>
  )
}
