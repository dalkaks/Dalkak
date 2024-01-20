import Image from 'next/image'
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from '../ui/navigation-menu'
import Link from 'next/link'

export default function Header() {
  return (
    <NavigationMenu className="w-full">
      <NavigationMenuList>
        <NavigationMenuItem>
          <Link href="/" legacyBehavior passHref>
            <NavigationMenuLink className={navigationMenuTriggerStyle()}>
              <Image
                src="/favicon.ico"
                alt="Image"
                width={20}
                height={20}
                className="rounded-md object-cover"
              />
              <h1 className="ml-1 text-2xl font-extrabold tracking-tight lg:text-4xl">
                dalkak
              </h1>
            </NavigationMenuLink>
          </Link>
        </NavigationMenuItem>
      </NavigationMenuList>
    </NavigationMenu>
  )
}
