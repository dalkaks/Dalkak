'use client';

import React from 'react';
import login, { RequestLogin } from '@/app/network/user/post/login';

const Button = ({ title, data }: { title: string; data: RequestLogin }) => {
  return (
    <button
      onClick={async () => {
        await login(data);
      }}
    >
      {title}
    </button>
  );
};

export default Button;
