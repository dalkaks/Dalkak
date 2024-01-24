package user

import (
	"dalkak/internal/interfaces"
)

type UserService interface {
}

type UserServiceImpl struct {
	db interfaces.Database
}

func NewUserService(db interfaces.Database) UserService {
	return &UserServiceImpl{db: db}
}
