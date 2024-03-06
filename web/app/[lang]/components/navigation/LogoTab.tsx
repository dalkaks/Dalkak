import React from 'react';
import SideTab from './containers/SideTab';
import Link from 'next/link';
import { useMediaQuery } from 'react-responsive';
import Image from 'next/image';

const LogoTab = () => {
  const Logo = () => (
    <Image
      className="h-[50%]"
      src="/images/dalkak_logo.png"
      alt="dalkak_logo"
      width={50}
      height={50}
    />
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
