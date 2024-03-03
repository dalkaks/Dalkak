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

type ConfirmMediaTempDto struct {
	UserInfo  *appdto.UserInfo
	Id        string
	MediaType string
	Prefix    string
}

func NewConfirmMediaTempDto(userInfo *appdto.UserInfo, id string, mediaType, prefix string) *ConfirmMediaTempDto {
	return &ConfirmMediaTempDto{
		UserInfo:  userInfo,
		Id:        id,
		MediaType: mediaType,
		Prefix:    prefix,
	}
}

type CreateMediaNftDto struct {
	UserInfo  *appdto.UserInfo
	Timestamp int64
	Prefix    string
	PrefixId  string
	ImageId   *string
	VideoId   *string
}

func NewCreateMediaNftDto(userInfo *appdto.UserInfo, timestamp int64, prefix, prefixId string, imageId, videoId *string) *CreateMediaNftDto {
	return &CreateMediaNftDto{
		UserInfo:  userInfo,
		Timestamp: timestamp,
		Prefix:    prefix,
		PrefixId:  prefixId,
		ImageId:   imageId,
		VideoId:   videoId,
	}
}
