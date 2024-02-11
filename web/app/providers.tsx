"use client";

import { MetaMaskProvider } from "@metamask/sdk-react";
import React from "react";

type Props = {
  children: any;
  href: string;
};

const Providers = ({ children, href }: Props) => {
  return (
    <MetaMaskProvider
      debug={false}
      sdkOptions={{
        dappMetadata: {
          name: "My dapp",
          url: href,
        },
      }}
    >
      {children}
    </MetaMaskProvider>
  );
};

export default Providers;
