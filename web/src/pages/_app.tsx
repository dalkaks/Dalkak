'use client'

import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import Head from 'next/head'
import Layout from '@/components/layout'
import { MetaMaskUIProvider } from '@metamask/sdk-react-ui'
import { useEffect, useState } from 'react'

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
        <MetaMaskUIProvider
          sdkOptions={{
            dappMetadata: {
              name: 'dalkak dapp',
              url: dappUrl,
            },
            // Other options
          }}
        >
          <Layout>
            <Component {...pageProps} />
          </Layout>
        </MetaMaskUIProvider>
      </div>
    </>
  )
}
