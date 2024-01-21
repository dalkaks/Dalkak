// src/locales/i18n.ts

import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import backend from 'i18next-http-backend'
import translationEN from './en/translation.json'
import translationKO from './ko/translation.json'

const resources = {
  en: {
    translation: translationEN,
  },
  ko: {
    translation: translationKO,
  },
}

i18n
  .use(LanguageDetector)
  .use(backend)
  .use(initReactI18next)
  .init({
    resources,
    // lng: 'ko', // 기본 설정 언어, 'cimode'로 설정할 경우 키 값으로 출력된다.
    fallbackLng: 'ko', // 번역 파일에서 찾을 수 없는 경우 기본 언어
    keySeparator: false, // 키 구분자 설정
    interpolation: {
      escapeValue: false,
    },
  })

export default i18n
