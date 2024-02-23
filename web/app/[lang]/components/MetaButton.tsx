'use client';

import { Button } from '@/components/ui/button';
import { useSDK } from '@metamask/sdk-react';
import React, { useEffect, useState } from 'react';
import detectEthereumProvider from '@metamask/detect-provider';

type MetaButtonProps = {
  title: string;
  setAccount: React.Dispatch<
    React.SetStateAction<{
      walletAddress: string;
      signature: string;
    }>
  >;
};

const MetaButton = ({ title, setAccount }: MetaButtonProps) => {
  const [hasProvider, setHasProvider] = useState(false);
  const { sdk, connected } = useSDK();

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
      const accounts = (await sdk?.connect({
        msg: '안전하게 지갑 연결'
      })) as string[];
      const signature = await sdk?.sign({ data: '안전하게 지갑 연결' });
      setAccount({ walletAddress: accounts?.[0], signature: signature });
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
  return (
    <Button disabled={!hasProvider} onClick={connected ? disconnect : connect}>
      {connected ? 'Disconnect' : 'Connect'}
    </Button>
  );
};

export default MetaButton;
