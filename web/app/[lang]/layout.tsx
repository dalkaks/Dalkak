import React from 'react';
import Navigation from './components/navigation';

type Props = {
  children: React.ReactNode;
};

function layout({ children }: Props) {
  return (
    <div>
      <Navigation />
      {children}
    </div>
  );
}

export default layout;
