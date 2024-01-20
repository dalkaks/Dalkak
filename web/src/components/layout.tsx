import Header from './common/header'
import { Separator } from './ui/separator'

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  return (
    <>
      <Header />
      <Separator />

      {/* <Navbar /> */}
      <main>{children}</main>
    </>
  )
}
