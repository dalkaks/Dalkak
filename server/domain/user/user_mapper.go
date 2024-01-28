package user

import "dalkak/pkg/interfaces"

func ConvertUserTableToUserDto(userTable *UserTable) *interfaces.UserDto {
	if userTable == nil {
		return nil
	}

	return &interfaces.UserDto{
		WalletAddress: userTable.WalletAddress,
		Timestamp:     userTable.Timestamp,
	}
}
