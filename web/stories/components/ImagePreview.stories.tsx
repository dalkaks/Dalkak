import ImagePreview from '@/app/[lang]/mint/components/ImagePreview';
import type { Meta, StoryObj } from '@storybook/react';

const meta = {
  title: 'Component/ImagePreview',
  component: ImagePreview,
  parameters: {
    layout: 'centered'
  },
  tags: ['autodocs'],
  argTypes: {
    file: { control: { type: 'file' } }
  },

  render: (_, { loaded }) => {
    return <ImagePreview file={loaded.file} />;
  }
} satisfies Meta<typeof ImagePreview>;

export default meta;
type Story = StoryObj<typeof meta>;

const url = 'https://nickelodeonuniverse.com/wp-content/uploads/Spongebob.png';

export const Primary: Story = {
  args: { file: new File([], '') },
  loaders: [
    async () => ({
      file: await fetch(url)
        .then(async (res) => await res.blob())
        .then((blob) => {
          return new File([blob], 'spongebob.png', { type: 'image/png' });
        })
    })
  ]
};
