package userdomain

import (
	"dalkak/internal/core"
	userobject "dalkak/internal/domain/user/object"
)

type UserDomainService interface {
	CreateNotRegisteredUser(walletAddress string) (*userobject.UserEntity, error)
}

type UserDomainServiceImpl struct {
	Database     UserRepository
	Storage      core.StorageManager
	Keymanager   core.KeyManager
	EventManager core.EventManager
}

func NewUserDomainService(database UserRepository, storage core.StorageManager, keymanager core.KeyManager, eventManager core.EventManager) UserDomainService {
	return &UserDomainServiceImpl{
		Database:     database,
		Storage:      storage,
		Keymanager:   keymanager,
		EventManager: eventManager,
	}
}

func (service *UserDomainServiceImpl) CreateNotRegisteredUser(walletAddress string) (*userobject.UserEntity, error) {
	user, err := service.Database.FindUserByWalletAddress(walletAddress)
	if err != nil || user != nil {
		return nil, err
	}
	newUser := userobject.NewUserEntity(walletAddress)
	return newUser, nil
}

