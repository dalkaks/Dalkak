'use client';

import React, { useState } from 'react';
import ImagePreview from './ImagePreview';
import MintingForm from './MintingForm';

const MintingLayout = () => {
  const [file, setFile] = useState<File>(new File([], ''));
  return (
    <div className="container flex h-1/2 w-full flex-col justify-around space-x-2 self-center py-44 sm:flex-row">
      <ImagePreview file={file} />
      <MintingForm setFile={setFile} />
    </div>
  );
};

export default MintingLayout;
