import React from 'react';

type Props = {
  children: React.ReactNode;
  className?: string;
};

const SideTab = ({ children, className }: Props) => {
  return (
    <div className={`flex w-2/12 items-center justify-center ${className}`}>
      {children}
    </div>
  );
};

export default SideTab;
