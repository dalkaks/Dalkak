'use client';
import login from '@/app/network/user/post/login';
import { MetaMaskProvider } from '@metamask/sdk-react';
import React, { useEffect, useState } from 'react';
import MetaButton from './MetaButton';
import Link from 'next/link';
import SideTab from './containers/SideTab';
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList
} from '@radix-ui/react-navigation-menu';
import { navigationMenuTriggerStyle } from '@/components/ui/navigation-menu';
import LogoTab from './LogoTab';
import WalletInfo from './wallet/WalletInfo';
import { useMediaQuery } from 'react-responsive';

const Navigation = () => {
  const [account, setAccount] = useState<{
    walletAddress: string;
    signature: string;
  }>({
    walletAddress: '',
    signature: ''
  });
  const [href, setHref] = useState('');
  useEffect(() => {
    setHref(window.location.href);
  }, []);
  const isDesktopAndTablet = useMediaQuery({
    query: '(min-width: 640px)'
  });

  useEffect(() => {
    account.walletAddress && account.signature && login(account);
    if (account.signature && account.walletAddress) alert('로그인 되었습니다');
  }, [account]);

  return (
    <div className="container sticky top-0 z-10 flex h-20 w-full justify-evenly bg-white shadow-md">
      <MetaMaskProvider
        debug={true}
        sdkOptions={{
          dappMetadata: {
            name: 'My dapp',
            url: href
          }
        }}
      >
        <LogoTab />
        <NavigationMenu className="flex w-7/12 items-center justify-center gap-10">
          <NavigationMenuList>
            <NavigationMenuItem>
              <Link href="/" legacyBehavior passHref>
                <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                  Home
                </NavigationMenuLink>
              </Link>
              <Link href="/mint" legacyBehavior passHref>
                <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                  Mint
                </NavigationMenuLink>
              </Link>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>

        <SideTab className="w-2/12 gap-2 sm:w-3/12 sm:gap-5">
          {isDesktopAndTablet && <WalletInfo />}
          <MetaButton setAccount={setAccount} />
        </SideTab>
      </MetaMaskProvider>
    </div>
  );
};

export default Navigation;
