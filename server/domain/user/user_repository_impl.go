package user

import (
	"dalkak/pkg/interfaces"
)

type UserRepositoryImpl struct {
	db interfaces.Database
}

func NewUserRepository(db interfaces.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (db *UserRepositoryImpl) FindOrCreateUser(walletAddress string) (string, error) {
	return "", nil
}
