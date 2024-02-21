package mediadto

import appdto "dalkak/pkg/dto/app"

type CreateTempMediaDto struct {
	UserInfo  *appdto.UserInfo
	MediaType string
	Ext       string
	Prefix    string
}

func NewCreateTempMediaDto(userInfo *appdto.UserInfo, mediaType, ext, prefix string) *CreateTempMediaDto {
	return &CreateTempMediaDto{
		UserInfo:  userInfo,
		MediaType: mediaType,
		Ext:       ext,
		Prefix:    prefix,
	}
}
