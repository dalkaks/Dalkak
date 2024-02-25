package boarddto

import (
	blockchaintype "dalkak/internal/infrastructure/blockchain/type"
	appdto "dalkak/pkg/dto/app"
)

type CreateBoardDto struct {
	UserInfo        *appdto.UserInfo
	Title           string
	Content         string
	ImageId         string
	VideoId         string
	ExternalLink    string
	BackgroundColor string
	Attributes      blockchaintype.NftAttribute
}
