package user

import (
	"dalkak/pkg/interfaces"
)

type UserServiceImpl struct {
	db interfaces.Database
}

func NewUserService(db interfaces.Database) interfaces.UserService {
	return &UserServiceImpl{db: db}
}
