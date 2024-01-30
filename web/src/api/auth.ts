const apiUrl = process.env.NEXT_PUBLIC_API_URL

export async function authenticateUserWithSignature(
  walletAddress: string,
  signature: string,
) {
  try {
    const response = await fetch(`${apiUrl}/user/auth`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        WalletAddress: walletAddress,
        Signature: signature,
      }),
    }).then((res) => res.json())

    localStorage.setItem('accessToken', response.data.accessToken)
  } catch (err) {
    throw err
  }
}
