'use client';

import React from 'react';
import presign from '@/app/network/media/post/presign';

const PsButton = ({ title }: { title: string }) => {
  return (
    <button
      onClick={() =>
        presign({
          mediaType: 'image',
          ext: 'jpeg',
          prefix: 'board'
        })
      }
    >
      {title}
    </button>
  );
};

export default PsButton;
