'use client';

import React from 'react';
import refresh from '@/app/network/user/post/refresh';

const RfButton = ({ title }: { title: string }) => {
  return <button onClick={async () => await refresh()}>{title}</button>;
};

export default RfButton;
