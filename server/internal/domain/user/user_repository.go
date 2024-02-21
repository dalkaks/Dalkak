package userdomain

import (
	userobject "dalkak/internal/domain/user/object"
)

type UserRepository interface {
	FindUserByWalletAddress(walletAddress string) (*userobject.UserEntity, error)
}