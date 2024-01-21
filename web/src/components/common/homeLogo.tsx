import Image from 'next/image'

export default function HomeLogo() {
  const handleLogoClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault() // 기본 링크 동작 방지
    window.location.href = '/' // 홈 화면으로 이동 및 새로고침
  }

  return (
    <a href="/" onClick={handleLogoClick} className="flex items-center">
      <Image
        src="/favicon.ico"
        alt="Image"
        width={25}
        height={25}
        className="rounded-md object-cover"
      />
      <h1 className="ml-2 text-2xl font-extrabold tracking-tight">dalkak</h1>
    </a>
  )
}
