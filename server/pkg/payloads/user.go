package payloads

import "dalkak/config"

type UserAuthAndSignUpRequest struct {
	WalletAddress string
	Signature     string
}

type UserAccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type UserBoardImagePresignedRequest struct {
	MediaType string `json:"mediaType"`
	Ext       string `json:"ext"`
}

func (req *UserBoardImagePresignedRequest) IsValid() bool {
	switch req.MediaType {
	case "image":
		if _, ok := config.AllowedImageExtensions[req.Ext]; ok {
			return true
		}
		// Todo: video
	}
	return false
}

type UserBoardImagePresignedResponse struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
