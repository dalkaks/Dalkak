import MintingForm from '@/app/[lang]/mint/components/MintingForm';
import type { Meta, StoryObj } from '@storybook/react';

const meta = {
  title: 'Component/MintingForm',
  component: MintingForm,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
  },
} satisfies Meta<typeof MintingForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    primary: true,
    label: '민팅 폼',
  },
};
