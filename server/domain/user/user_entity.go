package user

import "dalkak/pkg/dtos"

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

type UserBoardImageData struct {
	Pk         string
	Sk         string
	EntityType string
	Timestamp  int64

	Id          string
	Extension   string
	ContentType string
	Url         string
}

func GenerateUserBoardImageDataSk(boardImageId string) string {
	return `BoardImage#` + boardImageId
}

func (b *UserBoardImageData) ToBoardImageDto() *dtos.BoardImageDto {
	return &dtos.BoardImageDto{
		Id:          b.Id,
		Extension:   b.Extension,
		ContentType: b.ContentType,
		Url:         b.Url,
	}
}
