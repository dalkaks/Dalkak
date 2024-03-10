import presign, { RequestPresign } from './post/presign';
import upload from './put/upload';

export interface RequestUploadMedia {
  file: File;
  presign: RequestPresign;
}

async function uploadMedia(params: RequestUploadMedia) {
  const presignResult = await presign(params.presign);
  const { uploadUrl, id } = presignResult.data;
  console.log('UPLOAD URL:', uploadUrl);
  console.log(params.file.arrayBuffer());
  const uploadResult = await upload({
    uploadUrl,
    file: params.file
  });
}

export default uploadMedia;
