'use client';

import { Button } from '@/components/ui/button';
import { useSDK } from '@metamask/sdk-react';
import React, { useEffect, useState } from 'react';
import detectEthereumProvider from '@metamask/detect-provider';
import StatusDot from './StatusDot';
import { useMediaQuery } from 'react-responsive';
import WalletInfo from './wallet/WalletInfo';

type MetaButtonProps = {
  setAccount: React.Dispatch<
    React.SetStateAction<{
      walletAddress: string;
      signature: string;
    }>
  >;
};

const MetaButton = ({ setAccount }: MetaButtonProps) => {
  const [hasProvider, setHasProvider] = useState(false);
  const { sdk, connected } = useSDK();
  const isDesktopAndTablet = useMediaQuery({
    query: '(min-width: 640px)'
  });

  useEffect(() => {
    const getProvider = async () => {
      const provider = await detectEthereumProvider();

      if (provider) {
        setHasProvider(true);
      }
    };
    getProvider();
  }, []);

  const connect = async () => {
    try {
      if (!hasProvider) throw new Error('No provider found');

      const signature = (await sdk?.connectAndSign({
        msg: '안전하게 지갑 연결'
      })) as string;

      const walletAddress = window.ethereum?.selectedAddress as string;

      setAccount({ walletAddress, signature });
    } catch (err) {
      console.warn(`failed to connect..`, err);
    }
  };

  const disconnect = async () => {
    try {
      await sdk?.disconnect();
    } catch (err) {
      console.warn(`failed to disconnect..`, err);
    }
  };

  const handleConnection = async () => {
    if (connected) {
      await disconnect();
    } else {
      await connect();
    }
  };
  return (
    <Button disabled={!hasProvider} onClick={handleConnection}>
      <StatusDot connected={connected} />
      {isDesktopAndTablet && (
        <span>{connected ? 'Disconnect' : 'Connect'}</span>
      )}
      {!isDesktopAndTablet && <WalletInfo />}
    </Button>
  );
};

export default MetaButton;
