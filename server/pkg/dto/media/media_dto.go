package mediadto

import appdto "dalkak/pkg/dto/app"

type CreateMediaTempDto struct {
	UserInfo  *appdto.UserInfo
	MediaType string
	Ext       string
	Prefix    string
}

func NewCreateMediaTempDto(userInfo *appdto.UserInfo, mediaType, ext, prefix string) *CreateMediaTempDto {
	return &CreateMediaTempDto{
		UserInfo:  userInfo,
		MediaType: mediaType,
		Ext:       ext,
		Prefix:    prefix,
	}
}

type GetMediaTempDto struct {
	UserInfo  *appdto.UserInfo
	MediaType string
	Prefix    string
}

func NewGetMediaTempDto(userInfo *appdto.UserInfo, mediaType, prefix string) *GetMediaTempDto {
	return &GetMediaTempDto{
		UserInfo:  userInfo,
		MediaType: mediaType,
		Prefix:    prefix,
	}
}
