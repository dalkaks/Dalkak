'use server';

import serviceModule from '../../serviceModule';

async function deleteMedia() {
  serviceModule('DELETE', 'media', {
    mediaType: 'image',
    prefix: 'board'
  });
}

export default deleteMedia;
