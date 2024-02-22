package userdomain

import userentity "dalkak/internal/domain/user/object/entity"

type UserRepository interface {
	FindUserByWalletAddress(walletAddress string) (*userentity.UserEntity, error)
}