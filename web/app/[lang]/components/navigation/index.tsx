'use client';
import login from '@/app/network/user/post/login';
import { MetaMaskProvider } from '@metamask/sdk-react';
import React, { useEffect, useState } from 'react';
import MetaButton from './MetaButton';
import Link from 'next/link';
import { Button } from '@/components/ui/button';

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
    <div className="container flex-row">
      <MetaMaskProvider
        debug={true}
        sdkOptions={{
          dappMetadata: {
            name: 'My dapp',
            url: href
          }
        }}
      >
        <div className="container fixed flex justify-center ">
          <nav>
            <Link href="/">
              <Button>Home</Button>
            </Link>
            <Link href="/mint">
              <Button>Mint</Button>
            </Link>
            <MetaButton setAccount={setAccount} />
          </nav>
        </div>
      </MetaMaskProvider>
    </div>
  );
};

export default Navigation;
