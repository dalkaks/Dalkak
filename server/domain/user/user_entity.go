package user

import "dalkak/pkg/dtos"

type UserData struct {
	Pk            string
	Sk            string
	EntityType    string
	Timestamp     int64

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
