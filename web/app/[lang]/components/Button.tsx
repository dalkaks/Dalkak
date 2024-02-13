'use client';

import React from 'react';
import loginService, { RequestLogin } from '@/app/network/user/post/login';

const Button = ({ title, data }: { title: string; data: RequestLogin }) => {
  return <button onClick={() => loginService(data)}>{title}</button>;
};

export default Button;
