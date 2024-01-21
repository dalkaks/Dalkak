import UserInfo from './userInfo'
import HomeLogo from './homeLogo'

export default function Header() {
  return (
    <div className="flex justify-between items-center h-16 px-4">
      <HomeLogo />
      <UserInfo />
    </div>
  )
}
