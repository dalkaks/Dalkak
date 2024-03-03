import React from 'react';
import SideTab from './containers/SideTab';
import Link from 'next/link';
import { useMediaQuery } from 'react-responsive';

const LogoTab = () => {
  const Logo = () => (
    <img className="h-[50%]" src="/images/dalkak_logo.png" alt="dalkak_logo" />
  );

  const isDesktopAndTablet = useMediaQuery({
    query: '(min-width: 640px)'
  });

  return (
    <SideTab className="gap-1 pl-2 sm:gap-5">
      <Logo />
      {isDesktopAndTablet && (
        <Link href="/">
          <button className="text-2xl font-bold">Dalkak</button>
        </Link>
      )}
    </SideTab>
  );
};

export default LogoTab;
