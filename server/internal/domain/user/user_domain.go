package userdomain

import (
	"dalkak/internal/core"
	userobject "dalkak/internal/domain/user/object"
	userdto "dalkak/pkg/dto/user"
)

type UserDomainService interface {
	CheckAndCreateUser(*userdto.CheckAndCreateUserDto) (*userobject.UserEntity, error)
}

type UserDomainServiceImpl struct {
	Database     UserRepository
	Keymanager   core.KeyManager
	EventManager core.EventManager
}

func NewUserDomainService(database UserRepository, keymanager core.KeyManager, eventManager core.EventManager) UserDomainService {
	return &UserDomainServiceImpl{
		Database:     database,
		Keymanager:   keymanager,
		EventManager: eventManager,
	}
}

func (service *UserDomainServiceImpl) CheckAndCreateUser(dto *userdto.CheckAndCreateUserDto) (*userobject.UserEntity, error) {
	user, err := service.Database.FindUserByWalletAddress(dto.WalletAddress)
	if err != nil || user != nil {
		return nil, err
	}
	newUser := userobject.NewUserEntity(dto.WalletAddress)
	return newUser, nil
}
