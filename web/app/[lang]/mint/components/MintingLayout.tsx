'use client';

import React, { useState } from 'react';
import ImagePreview from './ImagePreview';
import MintingForm from './MintingForm';

const MintingLayout = () => {
  const [file, setFile] = useState<File>(new File([], '')); //여기 에러있음 바꿔야함 - File is not defined
  return (
    <div className="container flex h-full w-full flex-col self-center py-4 sm:h-1/2 sm:flex-row sm:justify-around sm:py-0">
      <ImagePreview file={file} />
      <MintingForm setFile={setFile} />
    </div>
  );
};

export default MintingLayout;
