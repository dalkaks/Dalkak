import { atom } from 'recoil'

export const accountState = atom<string | null>({
  key: 'accountState',
  default: null,
})
