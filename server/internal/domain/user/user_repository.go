package userdomain

import (
	"dalkak/internal/infrastructure/database/dao"
)

type UserRepository interface {
	FindUserByWalletAddress(walletAddress string) (*dao.UserDao, error)
}
