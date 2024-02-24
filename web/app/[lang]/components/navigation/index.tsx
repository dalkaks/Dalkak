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

  useEffect(() => {
    account.walletAddress && account.signature && login(account);
    if (account.signature && account.walletAddress) alert('로그인 되었습니다');
  }, [account]);

  return (
    <div className="container sticky top-0 flex h-20 bg-white shadow-md">
      <MetaMaskProvider
        debug={true}
        sdkOptions={{
          dappMetadata: {
            name: 'My dapp',
            url: href
          }
        }}
      >
        <SideTab className="gap-5">
          <img
            className="h-[80%]"
            src="/images/dalkak_logo.png"
            alt="dalkak_logo"
          />
          <span className="text-2xl font-bold">Dalkak</span>
        </SideTab>

        <NavigationMenu className="flex w-8/12 items-center justify-center gap-10">
          <NavigationMenuList>
            <NavigationMenuItem>
              <Link href="/">
                <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                  Home
                </NavigationMenuLink>
              </Link>
              <Link href="/mint">
                <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                  Mint
                </NavigationMenuLink>
              </Link>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>

        <SideTab>
          <MetaButton setAccount={setAccount} />
        </SideTab>
      </MetaMaskProvider>
    </div>
  );
};

export default Navigation;
