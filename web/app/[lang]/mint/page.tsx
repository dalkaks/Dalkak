'use client';

import React, { useState } from 'react';
import MintingForm from './components/MintingForm';
import ImagePreview from './components/ImagePreview';

const page = () => {
  const [file, setFile] = useState<File>(new File([], ''));
  return (
    <div className="flex">
      <ImagePreview file={file} />
      <MintingForm setFile={setFile} />
    </div>
  );
};

export default page;
