import MintingLayout from '@/app/[lang]/mint/components/MintingLayout';
import type { Meta, StoryObj } from '@storybook/react';

const meta = {
  title: 'Component/MintingLayout',
  component: MintingLayout,
  parameters: {
    layout: 'centered'
  },
  tags: ['autodocs'],
  argTypes: {}
} satisfies Meta<typeof MintingLayout>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {};
