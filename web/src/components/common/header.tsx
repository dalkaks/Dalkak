import Image from 'next/image'
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from '../ui/navigation-menu'
import Link from 'next/link'
import { MetaMaskButton } from '@metamask/sdk-react-ui'

export default function Header() {
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

      <MetaMaskButton theme={'light'} color="white"></MetaMaskButton>
    </div>
  )
}
