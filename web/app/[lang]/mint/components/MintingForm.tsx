'use client';

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form';
import React from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

import { z } from 'zod';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

const MAX_FILE_SIZE = 5000000;
const ACCEPTED_IMAGE_TYPES = [
  'image/jpeg',
  'image/jpg',
  'image/png',
  'image/webp'
];

const formSchema = z.object({
  title: z.string().min(2, { message: 'Name is too short' }),
  description: z.string(),
  file: z
    .any()
    .refine((file) => file?.size <= MAX_FILE_SIZE, `Max image size is 5MB.`)
    .refine(
      (file) => ACCEPTED_IMAGE_TYPES.includes(file?.type),
      'Only .jpg, .jpeg, .png and .webp formats are supported.'
    )
});

type Props = {
  setFile: React.Dispatch<React.SetStateAction<File>>;
};

const MintingForm = ({ setFile }: Props) => {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: '',
      description: '',
      file: new File([], 'none')
    }
  });
  function onSubmit(data: z.infer<typeof formSchema>) {
    console.log(data);
    alert('Form submitted');
  }
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          name="title"
          control={form.control}
          render={({ field }) => (
            <FormItem>
              <FormLabel>NFT 이름 입력</FormLabel>
              <FormControl>
                <Input
                  onChangeCapture={(e) => {
                    console.log(e.currentTarget.value);
                  }}
                  placeholder="이름"
                  {...field}
                />
              </FormControl>
              <FormDescription>NFT의 이름을 입력해 주세요.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          name="description"
          control={form.control}
          render={({ field }) => {
            return (
              <FormItem>
                <FormLabel>설명</FormLabel>
                <FormControl>
                  <Input
                    onChangeCapture={(e) => {
                      console.log(e.currentTarget.value);
                    }}
                    placeholder="설명"
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  민팅하시려는 NFT의 설명을 작성해 주세요.
                </FormDescription>
                <FormMessage />
              </FormItem>
            );
          }}
        />
        <FormField
          name="file"
          control={form.control}
          render={({ field }) => {
            return (
              <FormItem>
                <FormLabel>설명</FormLabel>
                <FormControl>
                  <Input
                    onChangeCapture={(e) => {
                      if (e.currentTarget.files) {
                        form.setValue('file', e.currentTarget.files[0]);
                        setFile(e.currentTarget.files[0]);
                      }
                      console.log(e.currentTarget.files);
                    }}
                    type="file"
                  />
                </FormControl>
                <FormDescription>
                  민팅하시려는 NFT의 설명을 작성해 주세요.
                </FormDescription>
                <FormMessage />
              </FormItem>
            );
          }}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
};

export default MintingForm;
