package boarddto

import (
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
	appdto "dalkak/pkg/dto/app"
)

type CreateBoardDto struct {
	UserInfo        *appdto.UserInfo
	Name            string
	Description     string
	ExternalLink    string
	BackgroundColor string
	Attributes      []boardvalueobject.NftAttribute
}

func NewCreateBoardDto(userInfo *appdto.UserInfo, name, description, externalLink, backgroundColor string, attributes []boardvalueobject.NftAttribute) *CreateBoardDto {
	return &CreateBoardDto{
		UserInfo:        userInfo,
		Name:            name,
		Description:     description,
		ExternalLink:    externalLink,
		BackgroundColor: backgroundColor,
		Attributes:      attributes,
	}
}
