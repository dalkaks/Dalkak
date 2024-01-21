'use client'

import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import Head from 'next/head'
import Layout from '@/components/layout'
import { MetaMaskProvider } from '@metamask/sdk-react'
import { useEffect, useState } from 'react'
import '../locales/i18n'

export default function App({ Component, pageProps }: AppProps) {
  const [dappUrl, setDappUrl] = useState('')

  useEffect(() => {
    setDappUrl(window.location.href)
  }, [])

  return (
    <>
      <Head>
        <meta
          name="viewport"
          content="width=device-width, initial-scale=1.0, viewport-fit=cover"
        />
      </Head>
      <div>
        <MetaMaskProvider
          debug={false}
          sdkOptions={{
            dappMetadata: {
              name: 'Example React Dapp',
              url: dappUrl,
            },
            // Other options
          }}
        >
          <Layout>
            <Component {...pageProps} />
          </Layout>
        </MetaMaskProvider>
      </div>
    </>
  )
}
