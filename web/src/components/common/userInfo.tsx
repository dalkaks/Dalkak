import { useSDK } from '@metamask/sdk-react'
import { useEffect } from 'react'
import { useRecoilState } from 'recoil'
import { accountState } from '@/state/accountState'
import { DotFilledIcon, ReloadIcon } from '@radix-ui/react-icons'
import { Button } from '../ui/button'
import { useGetProvider } from '@/hook/account/useGetProvider'
import { useConnectWallet } from '@/hook/account/useConnectWallet'

export default function UserInfo() {
  const [account, setAccount] = useRecoilState(accountState)
  const { connected, connecting, provider, chainId } = useSDK()

  const getProvider = useGetProvider()
  const connect = useConnectWallet()

  useEffect(() => {
    if (!account && connected) {
      getProvider()
    }
  }, [connected])

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
