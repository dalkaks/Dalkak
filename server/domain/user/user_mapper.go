package user

import "dalkak/pkg/dtos"

func ConvertUserTableToUserDto(userTable *UserTable) *dtos.UserDto {
	if userTable == nil {
		return nil
	}

	return &dtos.UserDto{
		WalletAddress: userTable.WalletAddress,
		Timestamp:     userTable.Timestamp,
	}
}
