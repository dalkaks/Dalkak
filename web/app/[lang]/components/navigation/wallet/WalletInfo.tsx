import { useSDK } from '@metamask/sdk-react';
import React from 'react';
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger
} from '@/components/ui/hover-card';
import Image from 'next/image';

type Props = {};

const WalletInfo = (props: Props) => {
  const { connected } = useSDK();

  if (!connected) return <></>;

  const WalletIcon = () => (
    <Image className="h-full" src="/icons/wallet.svg" alt="wallet" />
  );

  return (
    <HoverCard>
      <HoverCardTrigger className="h-[50%]">
        <WalletIcon />
      </HoverCardTrigger>
      <HoverCardContent className="relative top-5 w-48 border-solid">
        <h4 className="w-full overflow-hidden text-ellipsis text-sm font-semibold">
          {window.ethereum?.selectedAddress}
        </h4>
      </HoverCardContent>
    </HoverCard>
  );
};

export default WalletInfo;
