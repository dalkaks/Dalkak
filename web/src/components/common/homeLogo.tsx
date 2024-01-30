import Image from 'next/image'
import Link from 'next/link'

export default function HomeLogo() {
  return (
    <Link href="/" passHref className="flex items-center">
      <Image
        src="/favicon.ico"
        alt="Image"
        width={25}
        height={25}
        className="rounded-md object-cover"
      />
      <h1 className="ml-2 text-2xl font-extrabold tracking-tight">dalkak</h1>
    </Link>
  )
}
