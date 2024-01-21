import Header from './common/header'
import { Separator } from './ui/separator'

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="flex-container">
      <Header />
      <Separator />

      {/* <Navbar /> */}
      <main className="bg-secondary flex-main">{children}</main>
    </div>
  )
}
