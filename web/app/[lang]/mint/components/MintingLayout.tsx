'use client';

import React, { useState } from 'react';
import ImagePreview from './ImagePreview';
import MintingForm from './MintingForm';

const MintingLayout = () => {
  const [file, setFile] = useState<File>(new File([], ''));
  return (
    <div className="container flex h-1/2 w-full justify-around space-x-2 self-center">
      <ImagePreview file={file} />
      <MintingForm setFile={setFile} />
    </div>
  );
};

export default MintingLayout;
