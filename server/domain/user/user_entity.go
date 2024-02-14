package user

import (
	"dalkak/pkg/dtos"
)

type UserData struct {
	Pk         string
	Sk         string
	EntityType string
	Timestamp  int64

	WalletAddress string
}

const UserDataType = "User"

func GenerateUserDataPk(walletAddress string) string {
	return UserDataType + `#` + walletAddress
}

func (u *UserData) ToUserDto() *dtos.UserDto {
	return &dtos.UserDto{
		WalletAddress: u.WalletAddress,
		Timestamp:     u.Timestamp,
	}
}

type UserMediaData struct {
	Pk         string
	Sk         string
	EntityType string
	Timestamp  int64

	Id          string
	Prefix      string
	Extension   string
	ContentType string
	Url         string
	IsConfirm   bool
}

func GenerateUserBoardImageDataSk(prefix string, mediaType string) string {
	return `Media#` + prefix + `#` + mediaType
}

func (b *UserMediaData) ToMediaMeta() *dtos.MediaMeta {
	return &dtos.MediaMeta{
		ID:          b.Id,
		Prefix:      b.Prefix,
		Extension:   b.Extension,
		ContentType: b.ContentType,
		URL:         b.Url,
	}
}
