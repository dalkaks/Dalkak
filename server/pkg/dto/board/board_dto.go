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

func NewCreateBoardDto(userInfo *appdto.UserInfo, title, content, imageId, videoId, externalLink, backgroundColor string, attributes blockchaintype.NftAttribute) *CreateBoardDto {
	return &CreateBoardDto{
		UserInfo:        userInfo,
		Title:           title,
		Content:         content,
		ImageId:         imageId,
		VideoId:         videoId,
		ExternalLink:    externalLink,
		BackgroundColor: backgroundColor,
		Attributes:      attributes,
	}
}
