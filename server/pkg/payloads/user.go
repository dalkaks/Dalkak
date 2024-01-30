package payloads

type UserAuthAndSignUpRequest struct {
	WalletAddress string
	Signature     string
}

type UserAuthAndSignUpResponse struct {
	AccessToken string `json:"accessToken"`
}
