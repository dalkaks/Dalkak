'use client';

import React from 'react';
import presign from '@/app/network/media/post/presign';

const PsButton = ({ title }: { title: string }) => {
  return (
    <button
      onClick={async () => {
        try {
          await presign({
            mediaType: 'image',
            ext: 'jpeg',
            prefix: 'board'
          });
        } catch (error) {
          alert(error);
        }
      }}
    >
      {title}
    </button>
  );
};

export default PsButton;
