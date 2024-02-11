package user

import (
	"dalkak/pkg/dtos"
	"net/http"
	"strings"
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
}

func GenerateUserBoardImageDataSk(prefix string, contentType string) (string, error) {
	parts := strings.Split(contentType, "/")
	if len(parts) < 2 {
		return "", &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to split content type",
		}
	}
	mediaType := parts[0]
	return `Media#` + prefix + `#` + mediaType, nil
}

func (b *UserMediaData) ToBoardImageDto() *dtos.BoardImageDto {
	return &dtos.BoardImageDto{
		Id:          b.Id,
		Extension:   b.Extension,
		ContentType: b.ContentType,
		Url:         b.Url,
	}
}
