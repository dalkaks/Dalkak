import Image from 'next/image'
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from '../ui/navigation-menu'
import Link from 'next/link'
import { useSDK } from '@metamask/sdk-react'
import { useState } from 'react'

export default function Header() {
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
    <div className="flex justify-between items-center h-16">
      <NavigationMenu>
        <NavigationMenuList>
          <NavigationMenuItem>
            <Link href="/" legacyBehavior passHref>
              <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                <Image
                  src="/favicon.ico"
                  alt="Image"
                  width={25}
                  height={25}
                  className="rounded-md object-cover"
                />
                <h1 className="ml-2 text-2xl font-extrabold tracking-tight">
                  dalkak
                </h1>
              </NavigationMenuLink>
            </Link>
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>

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
      </div>
    </div>
  )
}
