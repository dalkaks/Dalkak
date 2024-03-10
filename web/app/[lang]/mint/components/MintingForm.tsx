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
import { NftFileExt } from '@/type/nft/fileExtension';
import uploadMedia, {
  RequestUploadMedia
} from '@/app/network/media/uploadMedia';
import deleteMedia from '@/app/network/media/delete/deleteMedia';

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
    alert('Form submitted');
  }
  async function onImageUpload(e: React.FormEvent<HTMLInputElement>) {
    if (e.currentTarget.files) {
      const file = e.currentTarget.files[0];
      form.setValue('file', file);
      setFile(file);

      const reqDto: RequestUploadMedia = {
        file: file,
        presign: {
          mediaType: 'image',
          ext: file.type.split('/')[1] as NftFileExt, //임시로 타입단언
          prefix: 'board'
        }
      };

      uploadMedia(reqDto)
        .then((res) => {
          console.log(res);
        })
        .catch((e) => {
          deleteMedia().then((res) => {
            console.log(res);
          });
          console.log(e);
        });
    }
  }
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
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
                      onImageUpload(e);
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
        <Button type="submit" className="mt-4 w-full">
          Submit
        </Button>
      </form>
    </Form>
  );
};

export default MintingForm;
