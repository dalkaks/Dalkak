import React from 'react';
import SideTab from './containers/SideTab';
import Link from 'next/link';

const LogoTab = () => {
  const Logo = () => (
    <img className="h-[50%]" src="/images/dalkak_logo.png" alt="dalkak_logo" />
  );

  return (
    <SideTab className="gap-5">
      <Logo />
      <Link href="/">
        <button className="text-2xl font-bold">Dalkak</button>
      </Link>
    </SideTab>
  );
};

export default LogoTab;
