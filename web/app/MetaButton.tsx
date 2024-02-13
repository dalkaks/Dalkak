'use client';

import { useSDK } from '@metamask/sdk-react';
import React, { useState } from 'react';

const MetaButton = ({ title }: { title: string }) => {
  const [account, setAccount] = useState<string>();
  const { sdk } = useSDK();

  const connect = async () => {
    try {
      const accounts = (await sdk?.connect()) as string[];
      console.log(accounts);
      setAccount(accounts?.[0]);
    } catch (err) {
      console.warn(`failed to connect..`, err);
    }
  };
  return <button onClick={connect}>{title}</button>;
};

export default MetaButton;
